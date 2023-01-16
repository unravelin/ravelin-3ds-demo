package handler

import (
	"log"
	"net/http"

	"github.com/unravelin/ravelin-3ds-demo/domain"
)

// MethodNotification is called by the customer's browser when the ACS has completed collecting
// info about the customer's browser.
//
// For more detail see: https://developer.ravelin.com/guides/3d-secure/browser-flow/#method-request
func (h Handler) MethodNotification(w http.ResponseWriter, r *http.Request) {
	addCommonHeaders(w, "text/html;charset=UTF-8")

	if r.Method == http.MethodOptions {
		return
	}

	log.Printf("Handling %s request", MethodNotificationEndpoint)

	methodNotificationResponse := &domain.MethodNotificationResponse{}
	err := decodeFormData(r, "threeDSMethodData", methodNotificationResponse)
	if err != nil {
		log.Printf("failed to decode Method Notification Response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.ThreeDSTransactionStore.SetMethodStatus(methodNotificationResponse.ThreeDSServerTransID, MethodStatusCompleted)
	if err != nil {
		log.Printf("failed to set method status for threeDSServerID %s: %v", methodNotificationResponse.ThreeDSServerTransID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.MethodNotificationResponseTemplate.Execute(w, methodNotificationResponse.ThreeDSServerTransID)
	if err != nil {
		log.Printf("failed to write method notification response - %s", err)
	}
}
