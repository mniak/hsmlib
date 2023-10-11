package app

import (
	"crypto/tls"
	"io"
	"net"
	"os"
)

type _Connection struct {
	io.ReadWriteCloser
}

func newConnection(target string, useTls bool, clientCertFile, clientKeyFile string, skipVerify bool) (*_Connection, error) {
	if !useTls {
		conn, err := net.Dial("tcp", target)
		return &_Connection{conn}, err
	}

	tlsConfig := tls.Config{
		InsecureSkipVerify: skipVerify,
	}
	if clientCertFile != "" && clientKeyFile != "" {
		clientCertBytes, err := os.ReadFile(clientCertFile)
		if err != nil {
			return nil, err
		}
		clientKeyBytes, err := os.ReadFile(clientKeyFile)
		if err != nil {
			return nil, err
		}
		cert, err := tls.X509KeyPair(clientCertBytes, clientKeyBytes)
		if err != nil {
			return nil, err
		}
		tlsConfig.Certificates = append(tlsConfig.Certificates, cert)
	}

	conn, err := tls.Dial("tcp", target, &tlsConfig)
	return &_Connection{conn}, err
}
