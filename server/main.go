package main

import (
	"flag"
	"go_proxy/basic"
	"log"
	"net"
	"os"
	"strings"
)

type StringArrayFlag []string

func (s *StringArrayFlag) String() string {
	return strings.Join(*s, ",")
}

func (s *StringArrayFlag) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("mode parameter is required: server")
	}
	model := os.Args[1]

	var ss basic.ServerProtocolHandler
	switch model {
	case "server":
		serverFlags := flag.NewFlagSet("server", flag.ExitOnError)
		listen := serverFlags.String("listen", "", "listen: address")
		//proto := serverFlags.String("proto", "", "proto: tcp/udp/http")
		//useTLS := serverFlags.Bool("tls", false, "enable tls")
		tlsCert := serverFlags.String("tls_cert", "", "tls cert")
		tlsKey := serverFlags.String("tls_key", "", "tls key")
		var clientIDs StringArrayFlag
		var passwords StringArrayFlag
		serverFlags.Var(&clientIDs, "client_id", "Multiple client IDs")
		serverFlags.Var(&passwords, "password", "Multiple passwords")
		err := serverFlags.Parse(os.Args[2:])
		if len(clientIDs) != len(passwords) {
			log.Fatalf("client_id and password count must match")
		}
		clientMap := make(map[string]string)
		for i := 0; i < len(clientIDs); i++ {
			clientMap[clientIDs[i]] = passwords[i]
		}

		if err != nil {
			log.Fatalf("failed to parse server flags: %v\n", err)
		}

		tlsConfig := basic.NewTLSConfig(*tlsCert, *tlsKey)
		ss = &basic.TCPServer{
			Addr:      *listen,
			TLS:       true,
			TLSConfig: tlsConfig,
		}
		_, err = ss.Listen()
		if err != nil {
			log.Fatalf("Listen failed: %v", err)
		}
		c := make(chan net.Conn, 1)
		for {
			conn, err := ss.Accept()
			if err != nil {
				log.Println("Accept error:", err)
				break
			}
			c <- conn
			go handleConnection(c, clientMap)

		}
	}
}
