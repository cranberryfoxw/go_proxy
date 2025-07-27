package main

import (
	"flag"
	"go_proxy/basic"
	"log"
	"os"
)

func main() {

	if len(os.Args) < 2 {
		log.Fatal("mode parameter is required: server")
	}
	model := os.Args[1]

	var cc basic.ClientProtocolHandler
	switch model {

	case "proxy_client":
		pfs := flag.NewFlagSet("pf", flag.ExitOnError)
		sa := pfs.String("server", "", "server address")
		tlsCert := pfs.String("tls_cert", "", "tls cert")
		tlsKey := pfs.String("tls_key", "", "tls key")

		clientId := pfs.String("client_id", "", "login id")
		password := pfs.String("password", "", "login password")

		remote := pfs.String("remote", "", "remote: server startup instance addr:port")
		local := pfs.String("local", "", "local: local listen addr:port")
		proto := pfs.String("proto", "", "remote: server use proto: tcp/udp/http")
		err := pfs.Parse(os.Args[2:])
		if err != nil {
			log.Fatalf("failed to parse pfs flags: %v\n", err)
		}

		tlsConfig := basic.NewTLSConfig(*tlsCert, *tlsKey)
		serverConfig := &basic.ServerInstanceConfig{
			Remote: *remote,
			Local:  *local,
			Proto:  *proto,
		}
		cc = &basic.TCPClient{
			Addr:         *sa,
			TLS:          true,
			TLSConfig:    tlsConfig,
			ServerConfig: serverConfig,
			Login:        false,
			ClientId:     *clientId,
			Password:     *password,
		}
		_, err = cc.Dial()
		err = cc.Verify()
		if err != nil {
			log.Fatal("contact server error -> ", err)
		}

	}
}
