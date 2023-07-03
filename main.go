package main

import (
	"context"
	"io"
	"log"
	"net"
	"sync/atomic"
	"time"

	"github.com/samber/lo"
	"github.com/spf13/viper"
)

func h(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	viper.SetDefault("listen_addr", "0.0.0.0:1500")
	viper.SetDefault("timeout_seconds", 5)

	listenAddress := viper.GetString("listen_addr")
	targetAddress := viper.GetString("target_addr")
	timeoutSeconds := viper.GetInt("timeout")

	if targetAddress == "" {
		return
	}

	outaddr := lo.Must(net.ResolveTCPAddr("tcp", targetAddress))
	inaddr := lo.Must(net.ResolveTCPAddr("tcp", listenAddress))

	outConn := lo.Must(net.DialTCP("tcp", nil, outaddr))
	defer outConn.Close()

	reactor := NewReactor(outConn)
	reactor.Start()
	log.Println("reactor started")

	listener := lo.Must(net.ListenTCP("tcp", inaddr))
	inHSMProto := NewHSMProtocol()

	var connectionIDs atomic.Uint64
	for {
		inConn, err := listener.AcceptTCP()
		if err != nil {
			log.Println("failed to accept incoming connection", err)
			return
		}

		connectionID := connectionIDs.Add(1)
		log.Printf("incoming connection %d accepted\n", connectionID)

		go func(conn io.ReadWriteCloser, connID uint64) {
			defer conn.Close()

			log.Printf("waiting data on connection %d\n", connID)
			request, err := inHSMProto.Receive(conn)
			if err != nil {
				log.Println("failed to read packet", err)
				return
			}

			timeoutCtx, cancelCtx := context.WithTimeout(context.Background(), time.Duration(timeoutSeconds)*time.Second)
			defer cancelCtx()

			replyData, err := reactor.Post(timeoutCtx, request.Data)
			if err != nil {
				log.Printf("failed to forward message (ID %d:%02X): %v\n", connID, request.ID, err)
				return
			}

			err = inHSMProto.Send(conn, PacketWithID{
				ID:   request.ID,
				Data: replyData,
			})
			if err != nil {
				log.Printf("failed to send reply back (ID %d:%02X): %v\n", connID, request.ID, err)
				return
			}
		}(inConn, connectionID)
	}
}
