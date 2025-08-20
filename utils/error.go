package utils

import (
	"errors"
	"syscall"
)

var ErrorDeviceNotAllowedCall = errors.New("device not allowed to call this function")
var ErrorDeviceNotFound = errors.New("device not found")
var ErrorContextNil = errors.New("context is nil")
var ErrorDeviceNotOpen = errors.New("device not open")
var ErrorAsyncNeedCallBack = errors.New("async operation requires a callback handler")
var ErrorFormatNotMatched = errors.New("media type format does not match the expected format")
var ErrorParameterInvalid = errors.New("invalid parameter provided")
var ErrorInternalError = errors.New("internal error occurred")

func CheckError(err error) error {
	if err == nil {
		return nil
	}
	var errno syscall.Errno
	if errors.As(err, &errno) && errno != 0 {
		return err
	}
	return nil
}
