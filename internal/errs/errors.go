package errs

import "errors"

var (
	ErrUnknown = errors.New("this error is no define")
)

var (
	ErrInvalidRecord = errors.New("invalid record")
	ErrCRCValidation = errors.New("crc validation failed")
)
