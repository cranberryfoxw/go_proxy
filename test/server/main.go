package main

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"strings"
	"time"
)

// 生成自签名证书
func generateSelfSignedCert() (tls.Certificate, error) {
	// 生成私钥
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return tls.Certificate{}, err
	}

	// 设置证书信息
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: "localhost",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
	}

	// 创建证书
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return tls.Certificate{}, err
	}

	// 打包证书和私钥
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})

	return tls.X509KeyPair(certPEM, keyPEM)
}

func main() {
	// 生成自签名证书
	cert, err := generateSelfSignedCert()
	if err != nil {
		log.Fatalf("生成证书失败: %v", err)
	}

	// 配置TLS
	tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}}

	// 监听TLS连接
	listener, err := tls.Listen("tcp", ":8088", tlsConfig)
	if err != nil {
		log.Fatalf("监听端口失败: %v", err)
	}
	defer listener.Close()

	fmt.Println("TLS服务器已启动，监听端口 8088")

	// 接受客户端连接
	conn, err := listener.Accept()
	if err != nil {
		log.Fatalf("接受连接失败: %v", err)
	}
	defer conn.Close()

	fmt.Println("客户端已连接")

	// 启动协程处理客户端消息
	go handleClientMessages(conn)

	// 从标准输入读取消息并发送给客户端
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("请输入要发送给客户端的消息（输入q退出）: ")
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("读取输入失败: %v", err)
			break
		}

		message = strings.TrimSpace(message)
		if message == "q" {
			break
		}

		// 发送消息到客户端
		_, err = conn.Write([]byte(message + "\n"))
		if err != nil {
			log.Printf("发送消息失败: %v", err)
			break
		}
	}

	fmt.Println("服务器已关闭")
}

// 处理客户端发送的消息
func handleClientMessages(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("接收消息失败: %v", err)
			break
		}

		fmt.Printf("客户端: %s", message)
	}
}
   