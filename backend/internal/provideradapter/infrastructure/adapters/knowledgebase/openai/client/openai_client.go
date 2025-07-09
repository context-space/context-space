package client

import (
	"net/http"
	"time"

	openai "github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

const (
	defaultOpenaiBaseURL = "https://api.openai.com/v1"
	defaultApiKey        = ""
	defaultTimeout       = 30 * time.Second
)

// Define the OpenaiClient struct which wraps the SDK client
type OpenaiClient struct {
	SdkClient *openai.Client
	// No need for baseURL here as sdkClient handles it
}

// NewOpenaiClient creates a new wrapper client.
func NewOpenaiClient(apiKey string, baseURL string) *OpenaiClient { // Return *OpenaiClient
	if baseURL == "" {
		baseURL = defaultOpenaiBaseURL
	}

	if apiKey == "" {
		apiKey = defaultApiKey
	}

	httpClient := &http.Client{
		Timeout: defaultTimeout,
	}

	opts := []option.RequestOption{
		option.WithAPIKey(apiKey),
		option.WithHTTPClient(httpClient),
		option.WithBaseURL(baseURL),
	}

	sdkClient := openai.NewClient(opts...)

	// Return an instance of our wrapper struct
	return &OpenaiClient{
		SdkClient: &sdkClient,
	}
}
