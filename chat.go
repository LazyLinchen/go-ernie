package go_ernie

import (
	"context"
	"github.com/pkg/errors"
	"net/http"
)

const (
	ChatMessageRoleUser      = "user"
	ChatMessageRoleAssistant = "assistant"
)

var (
	ErrChatCompletionInvalidModel       = errors.New("this model is not supported with this method, please use CreateCompletion client method instead") //nolint:lll
	ErrChatCompletionStreamNotSupported = errors.New("streaming is not supported with this method, please use CreateChatCompletionStream")              //nolint:lll
)

type ChatCompletionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ResponseFormat string

var (
	ResponseFormatDefault ResponseFormat = "text"
	ResponseFormatJSON    ResponseFormat = "json"
)

type ChatCompletionRequest struct {
	Model        string                  `json:"model"`
	Messages     []ChatCompletionMessage `json:"messages"`
	Stream       bool                    `json:"stream,omitempty"`
	Temperature  float32                 `json:"temperature,omitempty"`
	TopP         float32                 `json:"top_p,omitempty"`
	PenaltyScore float32                 `json:"penalty_score,omitempty"`
	UserId       string                  `json:"user_id,omitempty"`
	ResponseFormat ResponseFormat          `json:"response_format,omitempty"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type ChatCompletionResponse struct {
	ID               string `json:"id"`
	Object           string `json:"object"`
	Created          int64  `json:"created"`
	SentenceId       int64  `json:"sentence_id"`
	IsEnd            bool   `json:"is_end"`
	Result           string `json:"result"`
	NeedClearHistory bool   `json:"need_clear_history"`
	Usage            Usage  `json:"usage"`
	ErrorCode        int    `json:"error_code"`
	ErrorMsg         string `json:"error_msg"`
}

func (c *Client) CreateChatCompletion(
	ctx context.Context,
	request ChatCompletionRequest,
) (response ChatCompletionResponse, err error) {
	if request.Stream {
		err = ErrChatCompletionStreamNotSupported
		return
	}
	if !isSupportedChatCompletionModel(request.Model) {
		err = ErrChatCompletionInvalidModel
		return
	}

	urlSuffix := chatCompletionUri(request.Model)

	req, err := c.newRequest(ctx, http.MethodPost, c.fullURL(urlSuffix), withBody(request))
	if err != nil {
		return
	}
	err = c.sendRequest(req, &response)
	if response.ErrorCode != 0 {
		errResp := ErrorResponse{&APIError{ErrorCode: response.ErrorCode, ErrorMsg: response.ErrorMsg}}
		err = errResp.Error
	}
	return
}
