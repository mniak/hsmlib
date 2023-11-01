package hsmlib

import "github.com/mniak/hsmlib/errcode"

func ErrorCodeResponse(errorCode errcode.ErrorCode) (Response, error) {
	return NewResponse(errorCode, nil), nil
}

func CommandDisabledResponse() (Response, error) {
	return ErrorCodeResponse(errcode.CommandDisabled)
}

func InvalidInputDataResponse() (Response, error) {
	return ErrorCodeResponse(errcode.InvalidInput)
}
