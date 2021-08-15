package main

import (
	"net"
	"strings"

	log "github.com/sirupsen/logrus"
)

type cln struct {
	C    chan string
	Name string
	Addr string
}

var OnlineClnMap map[string]cln

var Messagechan = make(chan string)

func Manager() {
	OnlineClnMap = make(map[string]cln)
	for {
		Msg := <-Messagechan
		temp := Msg[strings.Index(Msg, "[")+1 : strings.Index(Msg, "]")]
		for _, value := range OnlineClnMap {
			if value.Name != temp {
				value.C <- Msg
			}
		}
	}
}

func HandleConn(conn net.Conn) {
	defer conn.Close()
	//创建新连接用户信息的结构体,默认用户名是IP+port
	ccln := cln{}
	ccln.C = make(chan string)
	ccln.Addr = conn.RemoteAddr().String()
	ccln.Name = conn.RemoteAddr().String()

	//fmt.Println(ccln.Addr)
	//fmt.Println(ccln.Name)

	//将用户写进OnlinMap中
	OnlineClnMap[conn.RemoteAddr().String()] = ccln

	//给全局通道发送用户登录信息
	Messagechan <- "[" + ccln.Name + "]" + "  is login"

	//给每个客户端配置一个专用接受信息Go程
	go WritrToCilent(ccln, conn)

	go ReadFromClient(ccln, conn)

	//fmt.Println("before select")
	select {}

}

func WritrToCilent(ccln cln, conn net.Conn) {
	for msg := range ccln.C {
		conn.Write([]byte(msg + "\n"))
		//fmt.Println(string([]byte(msg + "\n")))
		//加换行符防止阻塞，缓存区默认不刷新
	}
}

func ReadFromClient(ccln cln, conn net.Conn) {
	buf := [512]byte{}
	for {
		n, err := conn.Read(buf[:])
		if err != nil {
			log.WithError(err).Error("Server Read the cilent's msg failed")
			return
		}
		Messagechan <- "[" + ccln.Name + "]:" + string(buf[:n])
	}
}
