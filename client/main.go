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

	var cc basic.ProtocolHandler
	switch model {

	case "p":
		pfs := flag.NewFlagSet("pf", flag.ExitOnError)
		sa := pfs.String("server", "", "server address")
		tlsCert := pfs.String("tls_cert", "", "tls cert")
		tlsKey := pfs.String("tls_key", "", "tls key")

		clientId := pfs.String("client_id", "", "login id")
		pwd := pfs.String("pwd", "", "login password")

		remote := pfs.String("remote", "", "remote: server startup instance addr:port")
		local := pfs.String("local", "", "local: local listen addr:port")
		proto := pfs.String("proto", "", "remote: server use proto: tcp/udp/http")
	}
}
