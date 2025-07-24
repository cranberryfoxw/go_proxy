package basic

import (
	"crypto/tls"
	"errors"
	"log"
	"net"
	"sync"
)

type TCPServer struct {
	Addr      string
	TLS       bool
	TLSConfig *tls.Config
	Listener  net.Listener
	mu        sync.Mutex
	stopped   bool
}

func (s *TCPServer) Listen() (net.Listener, error) {
	ln, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return nil, err
	}

	// if use tls
	if s.TLS {
		// if not tls cert key
		if s.TLSConfig == nil {
			// generate self sign cert
			s.TLSConfig, err = GenerateSelfSignedTLSConfig([]string{"localhost"})
			if err != nil {
				log.Fatalln("Failed to generate TLS config:", err)
			}
		}
		ln = tls.NewListener(ln, s.TLSConfig)
	}
	s.Listener = ln
	log.Println("TCP Server started at", s.Addr)
	return ln, nil
}

func (s *TCPServer) Accept() (net.Conn, error) {
	s.mu.Lock()
	if s.stopped {
		s.mu.Unlock()
		return nil, errors.New("server already stopped")
	}
	s.mu.Unlock()

	if s.Listener == nil {
		return nil, errors.New("listener not initialized")
	}
	return s.Listener.Accept()
}

func (s *TCPServer) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.stopped {
		return nil
	}
	s.stopped = true
	if s.Listener != nil {
		err := s.Listener.Close()
		if err != nil {
			return err
		}
		s.Listener = nil
		log.Println("TCP Server stopped at", s.Addr)
	}
	return nil
}
