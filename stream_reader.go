package go_ernie

import (
	"bufio"
	"bytes"
	"fmt"
	erinie "github.com/LazyLinchen/go-ernie/internal"
	"io"
	"net/http"
)

var (
	headerData  = []byte("data: ")
	errorSuffix = []byte(`"error_msg"`)
)

type streamable interface {
	ChatCompletionStreamResponse
}

type streamReader[T streamable] struct {
	emptyMessagesLimit uint
	isFinished         bool

	reader         *bufio.Reader
	response       *http.Response
	errAccumulator erinie.ErrorAccumulator
	unmarshaler    erinie.Unmarshaler
}

func (stream *streamReader[T]) Recv() (response T, err error) {
	if stream.isFinished {
		err = io.EOF
		return
	}
	response, err = stream.processLines()
	return
}

func (stream *streamReader[T]) processLines() (T, error) {
	var (
		emptyMessagesCount uint
		hasErrorSuffix     bool
	)

	for {
		rawLine, readErr := stream.reader.ReadBytes('\n')
		if readErr != nil || hasErrorSuffix {
			respErr := stream.unmarshalError()
			if respErr != nil {
				return *new(T), fmt.Errorf("error, %w", respErr.Error)
			}
			return *new(T), readErr
		}

		noSpaceLine := bytes.TrimSpace(rawLine)
		if bytes.HasSuffix(noSpaceLine, errorSuffix) {
			hasErrorSuffix = true
		}
		if !bytes.HasPrefix(noSpaceLine, headerData) || hasErrorSuffix {
			if hasErrorSuffix {
				noSpaceLine = bytes.TrimPrefix(noSpaceLine, headerData)
			}
			writeErr := stream.errAccumulator.Write(noSpaceLine)
			if writeErr != nil {
				return *new(T), writeErr
			}
			emptyMessagesCount++
			if emptyMessagesCount > stream.emptyMessagesLimit {
				return *new(T), ErrTooManyEmptyStreamMessages
			}

			continue
		}
		noPrefixLine := bytes.TrimPrefix(noSpaceLine, headerData)
		if string(noPrefixLine) == "[DONE]" {
			stream.isFinished = true
			return *new(T), io.EOF
		}

		var response T
		unmarshalErr := stream.unmarshaler.Unmarshal(noPrefixLine, &response)
		if unmarshalErr != nil {
			return *new(T), unmarshalErr
		}
		return response, nil
	}
}

func (stream *streamReader[T]) unmarshalError() (errResp *ErrorResponse) {
	errBytes := stream.errAccumulator.Bytes()
	if len(errBytes) == 0 {
		return
	}
	err := stream.unmarshaler.Unmarshal(errBytes, &errResp)
	if err != nil {
		errResp = nil
	}
	return
}

func (stream *streamReader[T]) Close() {
	stream.response.Body.Close()
}
