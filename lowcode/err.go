package lowcode

import "errors"

var (
	// ErrUnkonow - unknow error
	ErrUnkonow = errors.New("unknow error")

	// ErrMustHaveMainPaytables - must have main paytables
	ErrMustHaveMainPaytables = errors.New("must have main paytables")
)
