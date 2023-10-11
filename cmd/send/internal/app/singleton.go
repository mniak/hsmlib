package app

import (
	"log"

	"github.com/mniak/hsmlib"
)

var __connection *_Connection

func Connect(target string, useTls bool, clientCertFile, clientKeyFile string, skipVerify bool) {
	var err error
	__connection, err = newConnection(target, useTls, clientCertFile, clientKeyFile, skipVerify)
	if err != nil {
		log.Fatalln(err)
	}
}

func conn() *_Connection {
	if __connection == nil || __connection.ReadWriteCloser == nil {
		log.Fatalln("There is no connection open")
		return nil
	}
	return __connection
}

func Finish() {
	if __connection != nil {
		if __connection.ReadWriteCloser != nil {
			__connection.Close()
		}
		__connection = nil
	}
}

func SendFrame(data []byte) hsmlib.ResponseWithHeader {
	return conn().SendFrame(data)
}

func SendPacket(packet hsmlib.Packet) hsmlib.ResponseWithHeader {
	return conn().SendPacket(packet)
}
