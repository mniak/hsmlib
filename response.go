package hsmlib

import (
	"bytes"
	"errors"
	"io"
)

type Response struct {
	ErrorCode string
	Data      []byte
}

func (r Response) ForCommandCode(commandCode []byte) ResponseWithCode {
	responseCode := CalculateResponseCode(commandCode)
	return ResponseWithCode{
		ResponseCode: responseCode,
		Response:     r,
	}
}

func CalculateResponseCode(commandCode []byte) string {
	if len(commandCode) != 2 {
		return "ZZ"
	}
	return string([]byte{commandCode[0], commandCode[1] + 1})
}

type ResponseWithCode struct {
	Response
	ResponseCode string
}

func (r ResponseWithCode) Bytes() []byte {
	var buf bytes.Buffer
	buf.WriteString(r.ResponseCode)
	buf.WriteString(r.ErrorCode)
	buf.Write(r.Data)
	return buf.Bytes()
}

func (r ResponseWithCode) WithHeader(header []byte) ResponseWithHeader {
	return ResponseWithHeader{
		Header:           header,
		ResponseWithCode: r,
	}
}

func ParseResponse(data []byte) (ResponseWithCode, error) {
	const codeLength = 2
	const errorCodeLength = 2

	if len(data) < codeLength {
		return ResponseWithCode{}, errors.New("response data does not contain a response code")
	}
	code := string(data[:codeLength])
	data = data[codeLength:]

	if len(data) < errorCodeLength {
		return ResponseWithCode{}, errors.New("response data does not contain an error code")
	}
	errorCode := string(data[:errorCodeLength])
	data = data[errorCodeLength:]

	respC := ResponseWithCode{
		ResponseCode: code,
		Response: Response{
			ErrorCode: errorCode,
			Data:      data,
		},
	}
	return respC, nil
}

type ResponseWithHeader struct {
	Header []byte
	ResponseWithCode
}

func ReceiveResponse(r io.Reader) (ResponseWithHeader, error) {
	packet, err := ReceivePacket(r)
	if err != nil {
		return ResponseWithHeader{}, err
	}
	respC, err := ParseResponse(packet.Payload)
	respH := ResponseWithHeader{
		Header:           packet.Header,
		ResponseWithCode: respC,
	}
	return respH, err
}

func (r ResponseWithHeader) AsPacket() Packet {
	return Packet{
		Header:  r.Header,
		Payload: r.ResponseWithCode.Bytes(),
	}
}
