package main

import (
	"fmt"
	"net"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

//用户结构体
type cln struct {
	C    chan string
	Name string
	Addr string
}

//全局用户在线map
var OnlineClnMap map[string]cln

//用户消息中间传输通道
var Messagechan = make(chan string)

//用户退出检测通道
var isQuite = make(chan bool)

//服务器管理消息广播进程
func Manager() {
	//实例化全局用户在线map
	OnlineClnMap = make(map[string]cln)
	for {
		Msg := <-Messagechan
		//通过标识符截取需要的信息
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
	OnlineClnMap[ccln.Addr] = ccln

	//给每个客户端配置一个专用接受信息Go程
	go WritrToCilent(ccln, conn)

	go ReadFromClient(ccln, conn)

	//给全局通道发送用户登录信息
	Messagechan <- "[" + ccln.Name + "]" + "  is login"

	//fmt.Println("before select")
	select {
	case <-isQuite:
		tempname := OnlineClnMap[ccln.Addr].Name
		delete(OnlineClnMap, ccln.Addr)
		Messagechan <- "[" + tempname + "]" + "  is logout"
		return
	case <-time.After(time.Second * 20):
		tempname := OnlineClnMap[ccln.Addr].Name
		delete(OnlineClnMap, ccln.Addr)
		Messagechan <- "[" + tempname + "]" + "  is logout"
		return

	}

}

func WritrToCilent(ccln cln, conn net.Conn) {
	for msg := range ccln.C {
		conn.Write([]byte(msg))
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
		if n == 1 {
			isQuite <- true
			fmt.Printf("检测到用户%s退出", ccln.Name)
			return
		}

		msg := string(buf[:n])

		if string(msg) == "who" && len(msg) == 3 {
			conn.Write([]byte("Online User : "))
			for _, user := range OnlineClnMap {
				conn.Write([]byte("[" + user.Addr + "]: " + string(user.Name) + "\n"))
			}
		} else if len(msg) >= 8 && msg[:6] == "rename" {
			NewName := strings.Split(msg, "|")[1]
			ccln.Name = NewName
			OnlineClnMap[ccln.Addr] = ccln
			conn.Write([]byte("Rename Successd"))

		} else {
			Messagechan <- "[" + ccln.Name + "]:" + string(buf[:n])
		}
	}
}
