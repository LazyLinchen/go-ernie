package go_ernie

import (
	"context"
	"net/http"
)

type ChatCompletionStreamChoice struct {
	Content string `json:"content,omitempty"`
	Role    string `json:"role,omitempty"`
}

type ChatCompletionStreamResponse struct {
	ID               string `json:"id"`
	Object           string `json:"object"`
	Created          int64  `json:"created"`
	SentenceId       int64  `json:"sentence_id"`
	IsEnd            bool   `json:"is_end"`
	Result           string `json:"result"`
	NeedClearHistory bool   `json:"need_clear_history"`
	Usage            Usage  `json:"usage"`
	ErrorMsg         string `json:"error_msg"`
	ErrorCode        string `json:"error_code"`
}

type ChatCompletionStream struct {
	*streamReader[ChatCompletionStreamResponse]
}

func (c *Client) CreateChatCompletionStream(
	ctx context.Context,
	request ChatCompletionRequest,
) (stream *ChatCompletionStream, err error) {
	if !isSupportedChatCompletionModel(request.Model) {
		err = ErrChatCompletionInvalidModel
		return
	}
	urlSuffix := chatCompletionUri(request.Model)

	request.Stream = true
	req, err := c.newRequest(ctx, http.MethodPost, c.fullURL(urlSuffix), withBody(request))
	if err != nil {
		return nil, err
	}
	resp, err := sendRequestStream[ChatCompletionStreamResponse](c, req)
	if err != nil {
		return
	}
	stream = &ChatCompletionStream{
		streamReader: resp,
	}
	return
}
