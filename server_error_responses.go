package hsmlib

func ErrorCodeResponse(errorCode string) (Response, error) {
	return Response{
		ErrorCode: errorCode,
	}, nil
}

func CommandDisabledResponse() (Response, error) {
	return ErrorCodeResponse("68") // Command has been disabled
}

func InvalidInputDataResponse() (Response, error) {
	return ErrorCodeResponse("15") // Invalid input data (invalid format, invalid characters, or not enough data provided)
}
