package main

import (
	"chatroom/server/model"
	"chatroom/server/utils"
	"crypto/tls"
	"log"
	"net"
	"time"
)

func process(conn net.Conn) {
	defer conn.Close()
	defer func() {
		if r := recover(); r != nil {
			utils.Log("Recovered from panic: %v", r)
		}
	}()
	processor := &Processor{
		Conn:          conn,
		RemoteAddress: conn.RemoteAddr().String(),
	}
	processor.process2()
}

func initUserDao() {
	model.MyUserDao = model.NewUserDao(pl)
}

func main() {
	initPool("localhost:6379", 16, 0, 300*time.Second)
	initUserDao()

	cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		log.Fatal("加载服务器的证书和私钥失败：", err)
	}

	//创建TLS配置
	config := &tls.Config{Certificates: []tls.Certificate{cert}}
	tlsListener, err := tls.Listen("tcp", "0.0.0.0:8889", config)
	if err != nil {
		utils.Log("tls.Listen err = %v", err)
		return
	}
	defer tlsListener.Close()

	utils.Log("服务器在8889端口监听......")

	for {
		utils.Log("等待客户端来连接服务器")
		conn, err := tlsListener.Accept()
		if err != nil {
			utils.Log("listen.Accept err=%v", err)
			continue
		}
		go process(conn)
	}
}
