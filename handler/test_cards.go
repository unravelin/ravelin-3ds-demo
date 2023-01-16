package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/unravelin/ravelin-3ds-demo/domain"
)

func (h Handler) TestCards(rw http.ResponseWriter, r *http.Request) {
	addCommonHeaders(rw, jsonContentType)
	if r.Method == http.MethodOptions {
		return
	}

	rsp, err := h.sendToRavelin3DSServer(http.MethodGet, nil, domain.RavelinThreeDSTestCardsEndpoint)
	if err != nil {
		log.Printf("failed to send Ravelin 3DS Test Cards Request: %v", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := readBody(rsp.Body)
	if err != nil {
		log.Printf("failed to read /3ds/testcards response: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	testCardsResponse := domain.RavelinTestCardsResponse{}
	err = json.Unmarshal(body, &testCardsResponse)
	if err != nil {
		log.Printf("failed to JSON decode /3ds/testcards response: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, err = json.Marshal(testCardsResponse.Data)
	if err != nil {
		log.Printf("failed to JSON encode Test Cards response: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = rw.Write(body)
	if err != nil {
		log.Printf("failed to write Test Cards response: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}
