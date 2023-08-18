package go_ernie

import "github.com/pkg/errors"

var (
	ErrTooManyEmptyStreamMessages = errors.New("stream has sent too many empty message")
)
