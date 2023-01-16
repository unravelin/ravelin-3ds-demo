package domain

import "encoding/json"

// This file contains example structs for Ravelin's 3DS API.
// For more detail see: https://developer.ravelin.com/apis/3d-secure/

const (
	RavelinThreeDSVersionEndpoint      = "/3ds/version"
	RavelinThreeDSAuthenticateEndpoint = "/3ds/authenticate"
	RavelinThreeDSResultEndpoint       = "/3ds/result"
	RavelinThreeDSTestCardsEndpoint    = "/3ds/testcards"
)

// RavelinVersionRequest
// For more detail see: https://developer.ravelin.com/apis/3d-secure/version/
type RavelinVersionRequest struct {
	TransactionID string `json:"transactionId,omitempty"`
	PAN           string `json:"pan,omitempty"`
}

type RavelinVersionResponse struct {
	Code      int                         `json:"status"`
	Message   string                      `json:"message,omitempty"`
	Timestamp int64                       `json:"timestamp,omitempty"`
	Data      *RavelinVersionResponseData `json:"data,omitempty"`
}

type RavelinVersionResponseData struct {
	TransactionID         string `json:"transactionId,omitempty"`
	ThreeDSServerTransID  string `json:"threeDSServerTransID,omitempty"`
	ThreeDSMethodURL      string `json:"threeDSMethodURL,omitempty"`
	VersionRecommendation string `json:"versionRecommendation,omitempty"`
}

type MethodNotificationResponse struct {
	ThreeDSServerTransID string `json:"threeDSServerTransID,omitempty"`
}

// RavelinAuthenticateRequest only contains the minimum number of fields required by Ravelin's 3DS API.
// Additional fields should be added and populated as required.
// For more detail see: https://developer.ravelin.com/apis/3d-secure/authenticate/
type RavelinAuthenticateRequest struct {
	Timestamp     int64    `json:"timestamp,omitempty"`
	CustomerID    string   `json:"customerId,omitempty"`
	TransactionID string   `json:"transactionId,omitempty"`
	AReqData      AReqData `json:"areqData,omitempty"`
}

type AReqData struct {
	MessageCategory                   string `json:"messageCategory,omitempty"`
	MessageVersion                    string `json:"messageVersion,omitempty"`
	DeviceChannel                     string `json:"deviceChannel,omitempty"`
	ThreeDSRequestorID                string `json:"threeDSRequestorID,omitempty"`
	ThreeDSRequestorName              string `json:"threeDSRequestorName,omitempty"`
	ThreeDSRequestorURL               string `json:"threeDSRequestorURL,omitempty"`
	ThreeDSServerTransID              string `json:"threeDSServerTransID,omitempty"`
	ThreeDSRequestorAuthenticationInd string `json:"threeDSRequestorAuthenticationInd,omitempty"`
	AcquirerMerchantID                string `json:"acquirerMerchantID,omitempty"`
	AcquirerBIN                       string `json:"acquirerBIN,omitempty"`
	PAN                               string `json:"pan,omitempty"`
	CardExpiryDate                    string `json:"cardExpiryDate,omitempty"`
	MerchantCountryCode               string `json:"merchantCountryCode,omitempty"`
	MerchantName                      string `json:"merchantName,omitempty"`
	MCC                               string `json:"mcc,omitempty"`
	PurchaseAmount                    string `json:"purchaseAmount,omitempty"`
	PurchaseCurrency                  string `json:"purchaseCurrency,omitempty"`
	PurchaseExponent                  string `json:"purchaseExponent,omitempty"`
	PurchaseDate                      string `json:"purchaseDate,omitempty"`
	ThreeDSCompInd                    string `json:"threeDSCompInd,omitempty"`
	BrowserAcceptHeader               string `json:"browserAcceptHeader,omitempty"`
	BrowserJavaEnabled                bool   `json:"browserJavaEnabled"`
	BrowserJavascriptEnabled          bool   `json:"browserJavascriptEnabled"`
	BrowserLanguage                   string `json:"browserLanguage,omitempty"`
	BrowserColorDepth                 string `json:"browserColorDepth,omitempty"`
	BrowserScreenHeight               string `json:"browserScreenHeight,omitempty"`
	BrowserScreenWidth                string `json:"browserScreenWidth,omitempty"`
	BrowserTZ                         string `json:"browserTZ,omitempty"`
	BrowserUserAgent                  string `json:"browserUserAgent,omitempty"`
	NotificationURL                   string `json:"notificationURL,omitempty"`
}

type RavelinAuthenticateResponse struct {
	Code      int                              `json:"status"`
	Message   string                           `json:"message,omitempty"`
	Timestamp int64                            `json:"timestamp,omitempty"`
	Data      *RavelinAuthenticateResponseData `json:"data,omitempty"`
}

type RavelinAuthenticateResponseData struct {
	MessageVersion        string            `json:"messageVersion,omitempty"`
	ThreeDSServerTransID  string            `json:"threeDSServerTransID,omitempty"`
	SDKTransID            string            `json:"sdkTransID,omitempty"`
	ACSChallengeMandated  string            `json:"acsChallengeMandated,omitempty"`
	ACSDecConInd          string            `json:"acsDecConInd,omitempty"`
	ACSOperatorID         string            `json:"acsOperatorID,omitempty"`
	ACSReferenceNumber    string            `json:"acsReferenceNumber,omitempty"`
	ACSTransID            string            `json:"acsTransID,omitempty"`
	ACSRenderingType      *ACSRenderingType `json:"acsRenderingType,omitempty"`
	ACSSignedContent      string            `json:"acsSignedContent,omitempty"`
	ACSURL                string            `json:"acsURL,omitempty"`
	AuthenticationType    string            `json:"authenticationType,omitempty"`
	AuthenticationValue   string            `json:"authenticationValue,omitempty"`
	BroadInfo             json.RawMessage   `json:"broadInfo,omitempty"`
	CardholderInfo        string            `json:"cardholderInfo,omitempty"`
	DSReferenceNumber     string            `json:"dsReferenceNumber,omitempty"`
	DSTransID             string            `json:"dsTransID,omitempty"`
	ECI                   string            `json:"eci,omitempty"`
	TransStatus           string            `json:"transStatus,omitempty"`
	TransStatusReason     string            `json:"transStatusReason,omitempty"`
	WhitelistStatus       string            `json:"whiteListStatus,omitempty"`
	WhitelistStatusSource string            `json:"whiteListStatusSource,omitempty"`

	ErrorCode        string `json:"errorCode,omitempty"`
	ErrorComponent   string `json:"errorComponent,omitempty"`
	ErrorDescription string `json:"errorDescription,omitempty"`
	ErrorDetail      string `json:"errorDetail,omitempty"`
	ErrorMessageType string `json:"errorMessageType,omitempty"`
}

type ACSRenderingType struct {
	ACSInterface  string `json:"acsInterface,omitempty"`
	ACSUITemplate string `json:"acsUiTemplate,omitempty"`
}

type ChallengeResponse struct {
	ThreeDSServerTransID   string                 `json:"threeDSServerTransID,omitempty"`
	ACSCounterAtoS         string                 `json:"acsCounterAtoS,omitempty"`
	ACSTransID             string                 `json:"acsTransID,omitempty"`
	ChallengeCompletionInd string                 `json:"challengeCompletionInd,omitempty"`
	MessageExtension       map[string]interface{} `json:"messageExtension,omitempty"`
	MessageType            string                 `json:"messageType,omitempty"`
	MessageVersion         string                 `json:"messageVersion,omitempty"`
	SDKTransID             string                 `json:"sdkTransID,omitempty"`
	TransStatus            string                 `json:"transStatus,omitempty"`
}

// RavelinResultRequest
// For more detail see: https://developer.ravelin.com/apis/3d-secure/result/
type RavelinResultRequest struct {
	ThreeDSServerTransID string `json:"threeDSServerTransID,omitempty"`
}

type RavelinResultResponse struct {
	Code      int                        `json:"status"`
	Message   string                     `json:"message,omitempty"`
	Timestamp int64                      `json:"timestamp,omitempty"`
	Data      *RavelinResultResponseData `json:"data,omitempty"`
}

type RavelinResultResponseData struct {
	ThreeDSServerTransID         string `json:"threeDSServerTransID,omitempty"`
	SDKTransID                   string `json:"sdkTransID,omitempty"`
	MessageVersion               string `json:"messageVersion,omitempty"`
	MessageCategory              string `json:"messageCategory,omitempty"`
	TransStatus                  string `json:"transStatus,omitempty"`
	TransStatusReason            string `json:"transStatusReason,omitempty"`
	ECI                          string `json:"eci,omitempty"`
	AuthenticationType           string `json:"authenticationType,omitempty"`
	AuthenticationValue          string `json:"authenticationValue,omitempty"`
	AuthenticationValueAlgorithm string `json:"authenticationValueAlgorithm,omitempty"`
	ChallengeCancel              string `json:"challengeCancel,omitempty"`
	InteractionCounter           string `json:"interactionCounter,omitempty"`
	WhiteListStatus              string `json:"whiteListStatus,omitempty"`
	WhiteListStatusSource        string `json:"whiteListStatusSource,omitempty"`
}

type RavelinTestCardsResponse struct {
	Code      int        `json:"status"`
	Message   string     `json:"message,omitempty"`
	Timestamp int64      `json:"timestamp,omitempty"`
	Data      []TestCard `json:"data,omitempty"`
}

type TestCard struct {
	TestPan     string `json:"testPan,omitempty"`
	Description string `json:"description,omitempty"`
}
