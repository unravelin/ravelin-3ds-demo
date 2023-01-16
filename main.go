package main

import (
	"embed"
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/unravelin/ravelin-3ds-demo/handler"
)

var (
	//go:embed static templates
	embeddedFS embed.FS
)

const (
	defaultRavelinApiUrl = "https://pci.ravelin.com"
	defaultMerchantUrl   = "http://localhost:8085"
)

func main() {
	var ravelinApiKey string
	var ravelinApiUrl string
	var merchantUrl string

	flag.StringVar(&ravelinApiKey, "ravelin-api-key", ravelinApiKey, "Ravelin API Key - Can also be set as $RAVELIN_API_KEY")
	flag.StringVar(&ravelinApiUrl, "ravelin-api-url", defaultRavelinApiUrl, "Ravelin API URL")
	flag.StringVar(&merchantUrl, "merchant-url", defaultMerchantUrl, "Merchant URL - If url does not contain a port, server is run on $PORT")
	flag.Parse()

	if ravelinApiKey == "" {
		ravelinApiKey = os.Getenv("RAVELIN_API_KEY")
		if ravelinApiKey == "" {
			panic("Ravelin API Key not set")
		}
	}

	if ravelinApiUrl == "" {
		panic("Ravelin API URL not set")
	}

	if merchantUrl == "" {
		panic("Merchant URL not set")
	}

	mUrl, err := url.Parse(merchantUrl)
	if err != nil {
		panic("failed to parse Merchant URL")
	}

	h := handler.Handler{
		RavelinApiUrl:           ravelinApiUrl,
		RavelinApiKey:           ravelinApiKey,
		MerchantUrl:             merchantUrl,
		ThreeDSTransactionStore: handler.NewThreeDSTransactionStore(),
	}

	h.MethodNotificationResponseTemplate, err = loadTemplate(embeddedFS, "templates/method-notification-response.html")
	if err != nil {
		panic(err)
	}

	h.ChallengeNotificationResponseTemplate, err = loadTemplate(embeddedFS, "templates/challenge-notification-response.html")
	if err != nil {
		panic(err)
	}

	staticFS, err := fs.Sub(embeddedFS, "static")
	if err != nil {
		panic(err)
	}
	frontEnd := http.FileServer(http.FS(staticFS))

	mux := http.NewServeMux()
	mux.Handle("/", frontEnd)
	mux.HandleFunc(handler.CheckoutEndpoint, h.Checkout)
	mux.HandleFunc(handler.AuthenticateEndpoint, h.Authenticate)
	mux.HandleFunc(handler.MethodNotificationEndpoint, h.MethodNotification)
	mux.HandleFunc(handler.ChallengeNotificationEndpoint, h.ChallengeNotification)
	mux.HandleFunc(handler.TestCardsEndpoint, h.TestCards)

	port := mUrl.Port()
	if port == "" {
		port = os.Getenv("PORT")
	}

	server := http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	log.Printf("Using Ravelin API URL %s", ravelinApiUrl)
	log.Printf("Starting server on port %q using merchant URL %s", server.Addr, merchantUrl)

	panic(server.ListenAndServe())
}

func loadTemplate(fs fs.ReadFileFS, filename string) (*template.Template, error) {
	file, err := fs.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to load template from file %s : %v", filename, err)
	}

	t, err := template.New("template").Parse(string(file))
	if err != nil {
		return nil, fmt.Errorf("failed to load template from file %s : %v", filename, err)
	}

	return t, nil
}
