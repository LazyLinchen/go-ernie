package go_ernie

import "net/http"

const (
	ernieAPIURL   = "https://aip.baidubce.com"
	ernieAPIURIv1 = "/rpc/2.0/ai_custom/v1"

	defaultEmptyMessagesLimit uint = 300
)

type ClientConfig struct {
	AK                 string
	SK                 string
	BaseURL            string
	AiApiURL           string
	APIVersion         string
	EmptyMessagesLimit uint

	HTTPClient *http.Client
}

func DefaultConfig(ak, sk string) ClientConfig {
	return ClientConfig{
		AK:                 ak,
		SK:                 sk,
		BaseURL:            ernieAPIURL,
		AiApiURL:           ernieAPIURIv1,
		HTTPClient:         &http.Client{},
		EmptyMessagesLimit: defaultEmptyMessagesLimit,
	}
}

func (ClientConfig) String() string {
	return "<Ernie API ClientConfig>"
	
}
