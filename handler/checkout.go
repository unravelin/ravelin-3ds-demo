package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"

	"github.com/unravelin/ravelin-3ds-demo/domain"
)

var (
	ErrCardRangeNotFound = errors.New("card range not found")
	ErrUnauthorised      = errors.New("authorization token not valid")
)

// Checkout is an example of the handler which is called when the customer click the "pay" button.
// This calls the Ravelin /3ds/version endpoint and handles the response.
//
// For more detail see: https://developer.ravelin.com/guides/3d-secure/browser-flow/#version-request
func (h Handler) Checkout(rw http.ResponseWriter, r *http.Request) {
	addCommonHeaders(rw, jsonContentType)

	if r.Method == http.MethodOptions {
		return
	}

	log.Printf("Handling %s request", CheckoutEndpoint)

	checkoutReqBytes, err := readBody(r.Body)
	if err != nil {
		log.Printf("failed to read checkout request body :%v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	checkoutRequest := domain.MerchantCheckoutRequest{}
	err = json.Unmarshal(checkoutReqBytes, &checkoutRequest)
	if err != nil {
		log.Printf("failed to unmarshal checkout request json: %v", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if checkoutRequest.AccountNumber == "" {
		log.Printf("invalid checkout request: %v", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	versionRequest := domain.RavelinVersionRequest{
		TransactionID: uuid.New().String(),
		PAN:           checkoutRequest.AccountNumber,
	}

	log.Printf("Making Ravelin /3ds/version request for card ending in %s", getLastFour(versionRequest.PAN))
	rsp, err := h.sendToRavelin3DSServer(http.MethodPost, versionRequest, domain.RavelinThreeDSVersionEndpoint)
	if err != nil {
		if err == ErrCardRangeNotFound {
			// card range not found
			log.Printf("Card range not found for current pan: %v", err)
			rw.WriteHeader(http.StatusNotFound)
			return

		}
		if err == ErrUnauthorised {
			// token not authorised
			log.Printf("API token not valid: %v", err)
			rw.WriteHeader(http.StatusUnauthorized)
			return

		}
		log.Printf("failed to send version request to threeds server: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rspBytes, err := readBody(rsp.Body)
	if err != nil {
		log.Printf("failed to read version response body: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	versionResponse := domain.RavelinVersionResponse{}
	err = json.Unmarshal(rspBytes, &versionResponse)
	if err != nil {
		log.Printf("failed to decode version response: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	if versionResponse.Data == nil {
		log.Printf("no version information in version response: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("Ravelin /3ds/version response received")

	methodStatus := MethodStatusNotCompleted // set to completed in method notification
	if versionResponse.Data.ThreeDSMethodURL == "" {
		methodStatus = MethodStatusUnavailable
	}

	tx := ThreeDSTransaction{
		MessageVersion: versionResponse.Data.VersionRecommendation,
		MethodStatus:   methodStatus,
	}
	h.ThreeDSTransactionStore.Add(versionResponse.Data.ThreeDSServerTransID, tx)

	checkoutResp := domain.MerchantCheckoutResponse{
		MessageVersion:        versionResponse.Data.VersionRecommendation,
		ThreeDSServerTransID:  versionResponse.Data.ThreeDSServerTransID,
		TransactionID:         versionResponse.Data.TransactionID,
		ThreeDSMethodURL:      versionResponse.Data.ThreeDSMethodURL,
		MethodNotificationURL: h.MerchantUrl + MethodNotificationEndpoint,
	}
	respond(checkoutResp, rw)
}
