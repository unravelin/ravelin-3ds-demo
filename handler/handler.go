package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/unravelin/ravelin-3ds-demo/domain"
)

const (
	jsonContentType = "application/json;charset=UTF-8"

	CheckoutEndpoint              = "/checkout"
	AuthenticateEndpoint          = "/authenticate"
	MethodNotificationEndpoint    = "/method-notification"
	ChallengeNotificationEndpoint = "/challenge-notification"
	TestCardsEndpoint             = "/test-cards"
)

type Handler struct {
	RavelinApiUrl                         string
	RavelinApiKey                         string
	MerchantUrl                           string
	ThreeDSTransactionStore               ThreeDSTransactionStore
	MethodNotificationResponseTemplate    *template.Template
	ChallengeNotificationResponseTemplate *template.Template
}

func (h Handler) sendToRavelin3DSServer(method string, body interface{}, endpoint string) (*http.Response, error) {
	var requestBody io.Reader

	if body == nil || method == http.MethodGet {
		requestBody = http.NoBody
	} else {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal threeds request: %v", err)
		}

		requestBody = bytes.NewBuffer(bodyBytes)
	}

	apiURL := h.RavelinApiUrl + endpoint

	request, err := http.NewRequest(method, apiURL, requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	request.Header.Set("Authorization", "token "+h.RavelinApiKey)
	request.Header.Set("Content-Type", jsonContentType)

	rsp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("http request fail: %v", err)
	}

	if rsp.StatusCode == http.StatusNotFound && endpoint == domain.RavelinThreeDSVersionEndpoint {
		return nil, ErrCardRangeNotFound
	}

	if rsp.StatusCode == http.StatusUnauthorized {
		return nil, ErrUnauthorised
	}

	if rsp.StatusCode != 200 {
		body, err := readBody(rsp.Body)
		if err == nil {
			log.Default().Printf("response body: %s\n", string(body))
		}
		return nil, fmt.Errorf("received bad status code %s ", rsp.Status)
	}

	return rsp, nil
}

func respond(data interface{}, rw http.ResponseWriter) {
	bb, err := json.Marshal(data)
	if err != nil {
		log.Printf("failed to encode data for response: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = rw.Write(bb)
	if err != nil {
		log.Printf("failed to write response: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func readBody(rc io.ReadCloser) ([]byte, error) {
	defer rc.Close()

	bb, err := io.ReadAll(rc)
	if err != nil {
		return nil, err
	}

	return bb, nil
}

func addCommonHeaders(rw http.ResponseWriter, contentType string) {
	rw.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	rw.Header().Set("X-Content-Type-Options", "nosniff")
	rw.Header().Set("Content-Type", contentType)
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Methods", "POST")
	rw.Header().Set("Access-Control-Allow-Headers", "*")
}

func getLastFour(pan string) string {
	if len(pan) > 4 {
		return pan[len(pan)-4:]
	} else {
		return pan
	}
}
