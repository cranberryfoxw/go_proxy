package basic

import "net"

type LoginMessage struct {
	ClientName string `json:"clientName"`
	Password   string `json:"password"`
}

type ServerInstanceConfig struct {
	Remote string `json:"remote"`
	Local  string `json:"local"`
	Proto  string `json:"proto"`
}

type ConfigMessage struct {
	InstanceType string `json:"instanceType"`
	Resources    string `json:"resources"`
}

type ServerProtocolHandler interface {
	Listen() (net.Listener, error)
	Accept() (net.Conn, error)
	Stop() error
}

type ClientProtocolHandler interface {
	Dial() (net.Conn, error)
	Verify() error
	Stop() error
}
