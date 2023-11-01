package errcode

import (
	"errors"
	"fmt"
	"strconv"
)

func Parse(errorCodeString string) (ErrorCode, error) {
	eint, err := strconv.ParseInt(errorCodeString, 16, 32)
	if err != nil {
		return 0, errors.New("response error code is not a hexadecimal string")
	}
	e := ErrorCode(eint)
	return e, nil
}

func ParseBCD(bytes []byte) (ErrorCode, error) {
	if len(bytes) != 2 {
		return 0, fmt.Errorf("failed to parse error code as BCD: expecting two bytes but got %d", len(bytes))
	}
	eint := (bytes[0]&0xF)<<4 | (bytes[1] & 0xF)
	e := ErrorCode(eint)
	return e, nil
}
