package handler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/unravelin/ravelin-3ds-demo/domain"
)

// ChallengeNotification is called by the customer's browser when the ACS has completed
// the challenge. This then calls Ravelin's /3ds/result endpoint to discover
// the outcome of the challenge.
//
// For more detail see:
// https://developer.ravelin.com/guides/3d-secure/browser-flow/#challenge-request
// https://developer.ravelin.com/guides/3d-secure/browser-flow/#result-request
func (h Handler) ChallengeNotification(w http.ResponseWriter, r *http.Request) {
	addCommonHeaders(w, "text/html;charset=UTF-8")

	if r.Method == http.MethodOptions {
		return
	}

	log.Printf("Handling %s request", ChallengeNotificationEndpoint)

	challengeResponse := &domain.ChallengeResponse{}
	err := decodeFormData(r, "cres", challengeResponse)
	if err != nil {
		log.Printf("failed to decode Challenge Response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Printf("/challenge-notification transStatus = %s", challengeResponse.TransStatus)

	resultRequest := &domain.RavelinResultRequest{
		ThreeDSServerTransID: challengeResponse.ThreeDSServerTransID,
	}

	log.Printf("Making Ravelin /3ds/result request for threeDSServerTransID %s", challengeResponse.ThreeDSServerTransID)

	rsp, err := h.sendToRavelin3DSServer(http.MethodPost, resultRequest, domain.RavelinThreeDSResultEndpoint)
	if err != nil {
		log.Printf("failed to send Result Request to ravelin threeds server: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := readBody(rsp.Body)
	if err != nil {
		log.Printf("failed to read 3ds server response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resultResponse := &domain.RavelinResultResponse{}
	err = json.Unmarshal(body, resultResponse)
	if err != nil {
		log.Printf("failed to read 3ds server response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	log.Printf("Ravelin /3ds/result response received. transStatus = %s", resultResponse.Data.TransStatus)

	var result string
	if resultResponse != nil && resultResponse.Data != nil &&
		(resultResponse.Data.TransStatus == "Y" || resultResponse.Data.TransStatus == "A") &&
		resultResponse.Data.AuthenticationValue != "" {
		result = "SUCCESS"
	} else {
		result = "FAILED"
	}

	err = h.ChallengeNotificationResponseTemplate.Execute(w, result)
	if err != nil {
		log.Printf("failed to write web challenge notification response - %s", err)
	}
}

// decodeFormData extracts a parameter from form data, and decodes it into a struct
func decodeFormData(r *http.Request, paramName string, decodeToStruct interface{}) error {
	paramValue, err := getFormVar(r, paramName)
	if err != nil {
		return err
	}

	// URL decode string first
	paramB64, err := url.QueryUnescape(paramValue)
	if err != nil {
		return fmt.Errorf("failed to unescape url encoded %s: %v", paramName, err)
	}

	// remove any padding
	paramB64 = strings.ReplaceAll(paramB64, "=", "")
	paramJSON, err := base64.RawStdEncoding.DecodeString(paramB64)
	if err != nil {
		return fmt.Errorf("failed to decode base64 %s: %v", paramName, err)
	}

	err = json.Unmarshal(paramJSON, decodeToStruct)
	if err != nil {
		return fmt.Errorf("failed to unmarshal %s JSON: %v", paramName, err)
	}

	return nil
}

func getFormVar(r *http.Request, paramName string) (string, error) {
	err := r.ParseForm()
	if err != nil {
		return "", fmt.Errorf("failed to parse form data - %s", err)
	}

	var paramValues []string
	// find first instance of param
	for k, v := range r.Form {
		if k == paramName {
			paramValues = v
			break
		}
	}

	if paramValues == nil {
		return "", fmt.Errorf("request does not contain `%s`", paramName)
	}

	if len(paramValues) > 1 {
		return "", fmt.Errorf("too many `%s` parameters", paramName)
	}

	return paramValues[0], nil
}
