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
	ChatCompletionResponse
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

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("Content-Type", "application/json; charset=utf-8")
	headers.Set("Cache-Control", "no-cache")

	request.Stream = true
	req, err := c.newRequest(ctx, http.MethodPost, c.fullURL(urlSuffix), withBody(request), withHeader(headers))
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
