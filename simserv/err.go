package simserv

import "errors"

var (
	// ErrUnkonow - unknow error
	ErrUnkonow = errors.New("unknow error")

	// ErrNonStatusOK - non statusOK
	ErrNonStatusOK = errors.New("non statusOK")
)
