// errors 统一定义bitcask.internal的错误
package errors

import "errors"

var (
	ErrUnknown = errors.New("this error is no define")
)

var (
	ErrInvalidRecord = errors.New("invalid record")
	ErrCRCValidation = errors.New("crc validation failed")
)
