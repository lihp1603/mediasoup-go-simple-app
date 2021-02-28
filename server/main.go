package main

import (
	"flag"
	"fmt"
)

func main() {
	port := flag.String("port", "8000", "set http port")
	ip := flag.String("ip", "127.0.0.1", "set http ip")
	cert := flag.String("cert", "/etc/certs/fullchain.pem", "set https ip")
	key := flag.String("key", "/etc/certs/privkey.pem", "set https ip")
	flag.Parse()

	fmt.Printf("ip:%s port :%d", *ip, *port)

	mediaServer := NewMediaServer()
	signalServer := NewSignalServer(mediaServer)
	signalServer.Start(*ip, *port, *cert, *key)

	wait := make(chan struct{})
	<-wait
}
