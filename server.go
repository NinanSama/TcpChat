package main

import (
	"net"

	log "github.com/sirupsen/logrus"
)

func ServerMain() {
	//创建监听套接字
	listener, err := net.Listen("tcp", "127.0.0.1:20000")
	if err != nil {
		log.WithError(err).Error("Listen is falied")
		return
	}

	defer listener.Close()

	//创建用户消息发送管理进程（优先启动）
	go Manager()

	//创建信号通道，接收来自系统的指令
	//signalChan := make(chan os.Signal, 2)
	//signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT)

	//循环监听客户端请求
	for {
		conn, err := listener.Accept()
		//fmt.Println("Accept a dail")
		if err != nil {
			log.WithError(err).Error("Accept is falied")
			return
		}
		//启动客户端处理请求Go程
		go HandleConn(conn)
	}

}
