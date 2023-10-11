package app

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/mniak/hsmlib"
	"golang.org/x/exp/slog"
)

func (conn *_Connection) SendFrame(data []byte) hsmlib.ResponseWithHeader {
	err := hsmlib.SendFrame(conn, data)
	if err != nil {
		log.Fatalln(err)
	}

	reply, err := hsmlib.ReceiveResponse(conn)
	if err != nil {
		log.Fatalln(err)
	}

	slog.Info("Response received:",
		"header", hex.EncodeToString(reply.Header),
		"response_code", reply.ResponseCode,
		"error_code", reply.ErrorCode,
		"data", fmt.Sprintf("%2X", reply.Data),
	)

	return reply
}

func (conn *_Connection) SendPacket(packet hsmlib.Packet) hsmlib.ResponseWithHeader {
	response := conn.SendFrame(packet.Bytes())
	if !bytes.Equal(response.Header, packet.Header) {
		log.Fatalln("invalid response header")
	}
	if len(packet.Payload) >= 2 {
		expectedResponseCode := packet.Payload[:2]
		expectedResponseCode[1]++
		if !bytes.Equal([]byte(response.ResponseCode), expectedResponseCode) {
			log.Fatalf("invalid response code: %q\n", response.ResponseCode)
		}
	}
	return response
}
