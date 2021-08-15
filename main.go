package main

import (
	"flag"

	log "github.com/sirupsen/logrus"
)

func main() {
	isServer := flag.Bool("isserver", false, "specify will this program act as a server")
	isCilent := flag.Bool("isclient", false, "specify will this program act as a client")
	flag.Parse()
	if !flag.Parsed() {
		flag.Usage()
		return
	}

	if *isServer {
		log.Info("Server Start Up")
		ServerMain()
	}
	if *isCilent {
		log.Info("Client Start Up")
		CilenMain()
	}

}
