package go_ernie

import (
	"bufio"
	"bytes"
	"context"
	erinie "github.com/LazyLinchen/go-ernie/internal"
	"io"
	"net/http"
)

type streamable interface {
	ChatCompletionStreamResponse | CompletionStreamResponse
}

type streamReader[T streamable] struct {
	emptyMessagesLimit uint
	isFinished         bool

	scanner        *bufio.Scanner
	response       *http.Response
	ctx            context.Context
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
	var eventData []byte
	if stream.scanner == nil {
		stream.scanner = bufio.NewScanner(stream.response.Body)
	}
	for len(eventData) == 0 {
		for {
			if !stream.scanner.Scan() {
				stream.isFinished = true
				return *new(T), stream.scanner.Err()
			}
			line := stream.scanner.Bytes()
			if len(line) == 0 {
				break
			}
			var value []byte
			if i := bytes.IndexRune(line, ':'); i != -1 {
				value = line[i+1:]
				if len(value) != 0 && value[0] == ' ' {
					value = value[1:]
				}
			}
			eventData = append(eventData, value...)
		}
	}
	var response T
	err := stream.unmarshaler.Unmarshal(eventData, &response)
	if err != nil {
		return *new(T), err
	}
	return response, nil
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
	_ = stream.response.Body.Close()
}
