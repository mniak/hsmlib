package app

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/mniak/hsmlib"
)

func (conn *_Connection) SendFrame(data []byte) hsmlib.ResponseWithHeader {
	logger.Info("Sending frame",
		"data", data,
	)
	err := hsmlib.SendFrame(conn, data)
	if err != nil {
		log.Fatalln(err)
	}

	reply, err := hsmlib.ReceiveResponse(conn)
	if err != nil {
		log.Fatalln(err)
	}

	logger.Info("Response received:",
		"header", hex.EncodeToString(reply.Header),
		"response_code", reply.ResponseCode,
		"error_code", fmt.Sprintf("%q", reply.ErrorCode()),
		"raw", reply.Bytes(),
		"data", reply.Data(),
	)

	return reply
}

func (conn *_Connection) SendPacket(packet hsmlib.Packet) hsmlib.ResponseWithHeader {
	logger.Info("Sending packet",
		"header", fmt.Sprintf("%2X", packet.Header),
		"payload", packet.Payload,
	)
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

func (conn *_Connection) SendPacketPayload(payload []byte) hsmlib.ResponseWithHeader {
	packet := hsmlib.Packet{
		Header:  RandomHeader(),
		Payload: payload,
	}
	return conn.SendPacket(packet)
}

func (conn *_Connection) SendCommand(cmd hsmlib.Command) hsmlib.ResponseWithHeader {
	logger.Info("Sending command",
		"code", cmd.Code(),
		"data", cmd.Data(),
	)
	payload := hsmlib.CommandBytes(cmd)
	packet := hsmlib.Packet{
		Header:  RandomHeader(),
		Payload: payload,
	}
	reply := conn.SendPacket(packet)
	return reply
}
