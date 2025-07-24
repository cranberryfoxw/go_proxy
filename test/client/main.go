package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	// 配置TLS，跳过证书验证
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true, // 跳过证书验证（仅用于测试）
	}

	// 连接TLS服务器
	conn, err := tls.Dial("tcp", "127.0.0.1:8088", tlsConfig)
	if err != nil {
		log.Fatalf("连接服务器失败: %v", err)
	}
	defer conn.Close()

	fmt.Println("已通过TLS连接到服务器")

	// 启动协程处理服务器发送的消息
	go handleServerMessages(conn)

	// 从标准输入读取消息并发送给服务器
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("请输入要发送给服务器的消息（输入q退出）: ")
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("读取输入失败: %v", err)
			break
		}

		message = strings.TrimSpace(message)
		if message == "q" {
			break
		}

		// 发送消息到服务器
		_, err = conn.Write([]byte(message + "\n"))
		if err != nil {
			log.Printf("发送消息失败: %v", err)
			break
		}
	}

	fmt.Println("客户端已关闭")
}

// 处理服务器发送的消息
func handleServerMessages(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("接收消息失败: %v", err)
			break
		}

		fmt.Printf("服务器: %s", message)
	}
}
   