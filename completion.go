package go_ernie

import (
	"context"
	"net/http"
)

type CompletionRequest struct {
	Model          string  `json:"model"`
	Prompt         string  `json:"prompt"`
	EndPoint       string  `json:"end_point"`
	Temperature    float64 `json:"temperature"`
	TopK           int     `json:"top_k"`
	TopP           float64 `json:"top_p"`
	PenaltyScore   float64 `json:"penalty_score"`
	Stream         bool    `json:"stream"`
	RetryCount     int     `json:"retry_count"`
	RequestTimeout float64 `json:"request_timeout"`
	UserId         string  `json:"user_id"`
}

type CompletionResponse struct {
	ID               string `json:"id"`
	Object           string `json:"object"`
	Created          int64  `json:"created"`
	SentenceId       int64  `json:"sentence_id"`
	IsEnd            bool   `json:"is_end"`
	IsTruncated      bool   `json:"is_truncated"`
	Result           string `json:"result"`
	NeedClearHistory bool   `json:"need_clear_history"`
	BanRound         int    `json:"ban_round"`
	Usage            Usage  `json:"usage"`
	ErrorCode        int    `json:"error_code"`
	ErrorMsg         string `json:"error_msg"`
}

func (c *Client) CreateCompletion(
	ctx context.Context,
	request CompletionRequest,
) (response CompletionResponse, err error) {
	if request.Stream {
		err = ErrChatCompletionStreamNotSupported
		return
	}
	if !isSupportedCompletionModel(request.Model) {
		err = ErrChatCompletionInvalidModel
		return
	}
	urlSuffix := completionUri(request.Model)

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
