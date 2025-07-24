package basic

import (
	"crypto/tls"
	"net"
	"sync"
)

type TCPClient struct {
	Addr      string
	TLS       bool
	TLSConfig *tls.Config
	Conn      net.Conn
	mu        sync.Mutex
	stopped   bool
}

func (s *TCPClient) Dial(network, address string) (net.Conn, error) {
	return nil, nil
}
