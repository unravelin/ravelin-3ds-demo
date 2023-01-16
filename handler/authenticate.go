package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/unravelin/ravelin-3ds-demo/domain"
)

// Authenticate calls the Ravelin /3ds/authenticate endpoint and handles the response.
// In this handler the merchant should provide merchant, customer and transaction data
// required by the 3DS authentication.
//
// For more detail see: https://developer.ravelin.com/guides/3d-secure/browser-flow/#authenticate-request
func (h Handler) Authenticate(w http.ResponseWriter, r *http.Request) {
	addCommonHeaders(w, jsonContentType)

	if r.Method == http.MethodOptions {
		return
	}

	log.Printf("Handling %s request", AuthenticateEndpoint)

	requestBody, err := readBody(r.Body)
	if err != nil {
		log.Printf("failed to read authenticate request body :%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	authenticateRequest := domain.MerchantAuthenticateRequest{}
	err = json.Unmarshal(requestBody, &authenticateRequest)
	if err != nil {
		log.Printf("failed to unmarshal authenticate request json: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = validateMerchantAuthenticateRequest(authenticateRequest)
	if err != nil {
		log.Printf("invalid authenticate request: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	acceptHeader := r.Header.Get("Accept")
	if authenticateRequest.BrowserData != nil {
		authenticateRequest.BrowserData.BrowserAcceptHeader = acceptHeader
	}

	ravelinAuthenticateRequest, err := h.createRavelinAuthenticateRequest(authenticateRequest)
	if err != nil {
		log.Printf("failed to create Ravelin 3DS Authenticate Request: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("Making Ravelin /3ds/authenticate request for card ending in %s", getLastFour(ravelinAuthenticateRequest.AReqData.PAN))
	rsp, err := h.sendToRavelin3DSServer(http.MethodPost, ravelinAuthenticateRequest, domain.RavelinThreeDSAuthenticateEndpoint)
	if err != nil {
		log.Printf("failed to send Ravelin 3DS Authenticate Request: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	responseBytes, err := readBody(rsp.Body)
	if err != nil {
		log.Printf("failed to read Ravelin 3DS Authenticate Response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ravelinAuthenticateResponse := &domain.RavelinAuthenticateResponse{}
	err = json.Unmarshal(responseBytes, ravelinAuthenticateResponse)
	if err != nil {
		log.Printf("failed to decode Ravelin 3DS Authenticate Response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	log.Printf("Ravelin /3ds/authenticate response received. MessageVersion: %s", ravelinAuthenticateResponse.Data.MessageVersion)

	merchantAuthenticateResponse := domain.MerchantAuthenticateResponse{}

	switch ravelinAuthenticateResponse.Data.MessageVersion {
	case "2.1.0", "2.2.0":
		switch ravelinAuthenticateResponse.Data.TransStatus {
		case "Y", "A":
			// proceed to authorisation
			merchantAuthenticateResponse.Status = "SUCCESS"
		case "C":
			merchantAuthenticateResponse.Status = "CHALLENGE_REQUIRED"
			merchantAuthenticateResponse.MessageVersion = ravelinAuthenticateResponse.Data.MessageVersion
			merchantAuthenticateResponse.ThreeDSServerTransID = ravelinAuthenticateResponse.Data.ThreeDSServerTransID
			merchantAuthenticateResponse.ACSTransID = ravelinAuthenticateResponse.Data.ACSTransID
			merchantAuthenticateResponse.ACSURL = ravelinAuthenticateResponse.Data.ACSURL
		case "N", "U", "R":
			merchantAuthenticateResponse.Status = "FAILED"
		}
	}

	respond(merchantAuthenticateResponse, w)
}

// createRavelinAuthenticateRequest prepares a Ravelin Authenticate request.
// For more detail see: https://developer.ravelin.com/apis/3d-secure/authenticate/
//
// Many of the fields in this function have been populated with example values for demonstration purposes.
// In a live implementation these fields should be populated with real merchant and transaction values.
func (h Handler) createRavelinAuthenticateRequest(request domain.MerchantAuthenticateRequest) (domain.RavelinAuthenticateRequest, error) {
	validColorDepth, err := convertToValidColorDepth(request.BrowserData.BrowserColorDepth)
	if err != nil {
		return domain.RavelinAuthenticateRequest{}, err
	}
	areqData := domain.AReqData{
		MessageCategory:                   "01",
		MessageVersion:                    request.MessageVersion,
		DeviceChannel:                     "02",
		ThreeDSRequestorAuthenticationInd: "01",
		ThreeDSRequestorID:                "example-3ds-merchant",
		ThreeDSRequestorName:              "Example 3DS Merchant",
		ThreeDSRequestorURL:               "https://www.ravelin.com/example-merchant",
		ThreeDSServerTransID:              request.ThreeDSServerTransID,
		AcquirerBIN:                       "000000999",
		PAN:                               request.AccountNumber,
		CardExpiryDate:                    request.CardExpiryDate,
		AcquirerMerchantID:                "9876543210001",
		MerchantCountryCode:               "826",
		MerchantName:                      "Example 3DS Merchant",
		MCC:                               "7922",
		PurchaseAmount:                    "80000",
		PurchaseCurrency:                  "826",
		PurchaseExponent:                  "2",
		PurchaseDate:                      time.Now().UTC().Format("20060102150405"),
		BrowserAcceptHeader:               request.BrowserData.BrowserAcceptHeader,
		BrowserJavaEnabled:                request.BrowserData.BrowserJavaEnabled,
		BrowserJavascriptEnabled:          request.BrowserData.BrowserJavascriptEnabled,
		BrowserLanguage:                   request.BrowserData.BrowserLanguage,
		BrowserColorDepth:                 strconv.Itoa(validColorDepth),
		BrowserScreenHeight:               strconv.Itoa(request.BrowserData.BrowserScreenHeight),
		BrowserScreenWidth:                strconv.Itoa(request.BrowserData.BrowserScreenWidth),
		BrowserTZ:                         strconv.Itoa(request.BrowserData.BrowserTZ),
		BrowserUserAgent:                  request.BrowserData.BrowserUserAgent,
		NotificationURL:                   h.MerchantUrl + ChallengeNotificationEndpoint,
	}

	tx, ok := h.ThreeDSTransactionStore.Get(request.ThreeDSServerTransID)
	if ok {
		areqData.MessageVersion = tx.MessageVersion
		areqData.ThreeDSCompInd = string(tx.MethodStatus)
	} else {
		areqData.ThreeDSCompInd = MethodStatusNotCompleted
	}

	r := domain.RavelinAuthenticateRequest{
		Timestamp:     time.Now().Unix(),
		CustomerID:    uuid.New().String(),
		TransactionID: uuid.New().String(),
		AReqData:      areqData,
	}

	return r, nil
}

// According to the EMVCo 3DS V220 spec 1, 4, 8, 15, 16, 24, 32 and 48 are the
// only valid values for color depth. Some browsers support 10 bit colors where
// the value for color depth can be 30 which will cause validation errors. For
// such cases the recommendation is to use the closest lower valid value.
func convertToValidColorDepth(colorDepth int) (int, error) {
	if colorDepth <= 0 {
		return 0, errors.New("color depth must be greater than zero")
	}

	valid := []int{1, 4, 8, 15, 16, 24, 32, 48}
	for i := 1; i < len(valid); i++ {
		if colorDepth == valid[i] {
			return colorDepth, nil
		}
		if colorDepth < valid[i] {
			return valid[i-1], nil
		}
	}

	return valid[len(valid)-1], nil
}

func validateMerchantAuthenticateRequest(request domain.MerchantAuthenticateRequest) error {
	if request.ProductQuantity == 0 { // Has to have quantity
		return fmt.Errorf("product quantity is zero")
	}

	if request.ProductSKU == "" { // Has to have a product
		return fmt.Errorf("no product selected")
	}

	return nil
}
