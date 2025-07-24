package basic

import "net"

type LoginMessage struct {
	ClientName string `json:"clientName"`
	Password   string `json:"password"`
}

type ConfigMessage struct {
	InstanceType string `json:"instanceType"`
	Resources    string `json:"resources"`
}

type ProtocolHandler interface {
	Listen() (net.Listener, error)
	Accept() (net.Conn, error)
	Stop() error
}
