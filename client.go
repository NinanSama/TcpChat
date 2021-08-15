package main

import (
	"net"

	log "github.com/sirupsen/logrus"
)

func CilenMain() {
	conn, err := net.Dial("tcp", "127.0.0.1:20000")
	if err != nil {
		log.WithError(err).Error("Dial failed")
		return
	}

	defer conn.Close()

	go ReadFromServer(conn)

	go WriteToServer(conn)

	select {}

}
