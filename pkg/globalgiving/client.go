package globalgiving

import (
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

type client struct {
	apiKey     string
	baseURL    *url.URL
	httpClient *http.Client
}

func newClient() *client {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &client{
		apiKey: os.Getenv("GLOBAL_GIVING_API_KEY"),
		baseURL: &url.URL{
			Scheme: "https",
			Host:   "api.globalgiving.org",
		},
		httpClient: &http.Client{},
	}
}
