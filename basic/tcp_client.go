package basic

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"log"
	"net"
	"sync"
)

type TCPClient struct {
	Addr         string
	TLS          bool
	TLSConfig    *tls.Config
	ClientId     string
	Password     string
	Conn         net.Conn
	ServerConfig *ServerInstanceConfig
	mu           sync.Mutex
	Login        bool
	stopped      bool
}

func (c *TCPClient) Dial() (net.Conn, error) {
	c.mu.Lock()
	if c.TLSConfig == nil {
		c.TLSConfig = &tls.Config{
			InsecureSkipVerify: true,
			MinVersion:         tls.VersionTLS12,
		}
	}
	conn, err := tls.Dial("tcp", c.Addr, c.TLSConfig)
	c.Conn = conn
	c.mu.Unlock()
	return conn, err
}

func (c *TCPClient) Stop() error {
	return nil
}

func (c *TCPClient) Verify() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.Login {
		return nil
	}
	if c.Conn == nil {
		return errors.New("no connection established")
	}
	loginMessage := LoginMessage{
		ClientName: c.ClientId,
		Password:   c.Password,
	}
	loginMessageBytes, _ := json.Marshal(loginMessage)

	_, err := c.Conn.Write(append(loginMessageBytes))
	if err != nil {
		return errors.New("[x] failed to send login message: " + err.Error())
	}

	buffer := make([]byte, 4096)
	n, err := c.Conn.Read(buffer)
	if err != nil {
		return errors.New("[x] failed to read login response: " + err.Error())
	}
	loginResp := string(buffer[:n])
	if loginResp != "OK" {
		return errors.New("[x] login failed")
	}
	c.Login = true
	log.Println("Login successful")
	return nil
}

func S() {

}
