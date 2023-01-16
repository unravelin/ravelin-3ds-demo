package domain

// This file contains example structs for the merchant's back-end API.
// These have only been provided for demonstration purposes.
// The format of the actual merchant back-end requests and responses can be defined by the merchant.

type MerchantCheckoutRequest struct {
	AccountNumber string `json:"accountNumber,omitempty"`
}

type MerchantCheckoutResponse struct {
	MessageVersion        string `json:"messageVersion,omitempty"`
	ThreeDSServerTransID  string `json:"threeDSServerTransID,omitempty"`
	TransactionID         string `json:"transactionId,omitempty"`
	ThreeDSMethodURL      string `json:"threeDSMethodURL,omitempty"`
	MethodNotificationURL string `json:"methodNotificationURL,omitempty"`
}

type MerchantAuthenticateRequest struct {
	MessageVersion       string       `json:"messageVersion,omitempty"`
	ThreeDSServerTransID string       `json:"threeDSServerTransID,omitempty"`
	ProductSKU           string       `json:"productSKU,omitempty"`
	ProductQuantity      int          `json:"productQuantity,omitempty"`
	AccountNumber        string       `json:"accountNumber,omitempty"`
	CardExpiryDate       string       `json:"cardExpiryDate,omitempty"`
	BrowserData          *BrowserData `json:"browserData,omitempty"`
}

type BrowserData struct {
	BrowserAcceptHeader      string `json:"browserAcceptHeader,omitempty"`
	BrowserJavaEnabled       bool   `json:"browserJavaEnabled"`
	BrowserJavascriptEnabled bool   `json:"browserJavascriptEnabled"`
	BrowserLanguage          string `json:"browserLanguage,omitempty"`
	BrowserColorDepth        int    `json:"browserColorDepth,omitempty"`
	BrowserScreenHeight      int    `json:"browserScreenHeight,omitempty"`
	BrowserScreenWidth       int    `json:"browserScreenWidth,omitempty"`
	BrowserTZ                int    `json:"browserTZ,omitempty"`
	BrowserUserAgent         string `json:"browserUserAgent,omitempty"`
	NotificationURL          string `json:"notificationURL,omitempty"`
}

type MerchantAuthenticateResponse struct {
	Status               string `json:"status"`
	ThreeDSServerTransID string `json:"threeDSServerTransID,omitempty"`
	ACSTransID           string `json:"acsTransID,omitempty"`
	ACSURL               string `json:"acsURL,omitempty"`
	MessageVersion       string `json:"messageVersion,omitempty"`
	Error                string `json:"error,omitempty"`
}
