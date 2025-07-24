package main

import (
	"encoding/json"
	"fmt"
	"go_proxy/basic"
	"log"
	"net"
	"sync"
)

func handleConnection(c chan net.Conn, clientMap map[string]string) {
	conn := <-c
	defer conn.Close()

	// login lock
	var loggedIn bool
	var loginMutex sync.Mutex

	// msg buff
	buffer := make([]byte, 4096)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			log.Printf("[x] read  client %s message error: %v", conn.RemoteAddr(), err)
			return
		}
		messageData := make([]byte, n)
		copy(messageData, buffer[:n])

		// check login
		loginMutex.Lock()
		isLoggedIn := loggedIn
		loginMutex.Unlock()

		// not login only rec login message
		if !isLoggedIn {
			var loginMsg basic.LoginMessage
			if err := json.Unmarshal(messageData, &loginMsg); err != nil {
				// pass not login message
				continue
			}

			// valid name pwd
			if AuthHandler(loginMsg, clientMap) {
				// success
				loginMutex.Lock()
				loggedIn = true
				loginMutex.Unlock()

				// send ok
				if _, err := conn.Write([]byte("OK")); err != nil {
					log.Printf("notify the client of login success or failure: %v", err)
				}

				// handle
				handleClientConfig(conn, loginMsg.ClientName)
			} else {
				// valid error
				if _, err := conn.Write([]byte("FAIL")); err != nil {
					log.Printf("Failed to send login failure response: %v", err)
				}
			}
		} else {
			continue
		}
	}
}

func handleClientConfig(conn net.Conn, clientName string) {
	buffer := make([]byte, 4096)

	// 读取配置消息
	n, err := conn.Read(buffer)
	if err != nil {
		log.Printf("读取客户端 %s 配置失败: %v", clientName, err)
		return
	}

	// 解析配置
	var config basic.ConfigMessage
	if err := json.Unmarshal(buffer[:n], &config); err != nil {
		log.Printf("解析客户端 %s 配置失败: %v", clientName, err)
		return
	}
	conn.Write([]byte("hahha"))
	go startInstance(clientName, config)
}

// 启动实例的函数（示例实现）
func startInstance(clientName string, config basic.ConfigMessage) {
	// 这里实现启动实例的逻辑
	log.Printf("为客户端 %s 启动实例: %+v", clientName, config)

	// 示例：模拟启动过程
	fmt.Printf("启动实例: 类型=%s, 资源=%s\n", config.InstanceType, config.Resources)
}

func AuthHandler(loginMessage basic.LoginMessage, clientMap map[string]string) bool {
	expectedPassword, registered := clientMap[loginMessage.ClientName]
	if !registered {
		return false
	}
	return loginMessage.Password == expectedPassword
}
