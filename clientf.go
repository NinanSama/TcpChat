package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

func ReadFromServer(conn net.Conn) {
	for {
		buf := [512]byte{}
		n, err := conn.Read(buf[:])
		if err != nil {
			log.WithError(err).Error("Read from server failed")
			return
		}
		fmt.Println(string(buf[:n]))
	}
}

func WriteToServer(conn net.Conn) {
	inputReader := bufio.NewReader(os.Stdin)
	//tempnum := 0
	for {
		input, _ := inputReader.ReadString('\n')
		inputInfo := strings.TrimSpace(input)
		_, err := conn.Write([]byte(inputInfo))
		//fmt.Println("Start send")
		if err != nil {
			log.WithError(err).Error("Cilent Send Msg to Server Failed")
			return
		}
		//tempnum++
		//fmt.Println(tempnum)

	}
}
