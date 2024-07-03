package go_ernie

import (
	"context"
	"net/http"
)

type CompletionStreamResponse struct {
	CompletionResponse
}

type CompletionStream struct {
	*streamReader[CompletionStreamResponse]
}

func (c *Client) CreateCompletionStream(
	ctx context.Context,
	request CompletionRequest,
) (stream *CompletionStream, err error) {
	if request.Stream {
		err = ErrChatCompletionStreamNotSupported
		return
	}
	if !isSupportedCompletionModel(request.Model) {
		err = ErrChatCompletionInvalidModel
		return
	}
	urlSuffix := completionUri(request.Model)

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("Content-Type", "application/json; charset=utf-8")
	headers.Set("Cache-Control", "no-cache")

	req, err := c.newRequest(ctx, http.MethodPost, c.fullURL(urlSuffix), withBody(request), withHeader(headers))
	if err != nil {
		return
	}
	resp, err := sendRequestStream[CompletionStreamResponse](c, req)
	if err != nil {
		return
	}
	stream = &CompletionStream{
		streamReader: resp,
	}
	return
}
