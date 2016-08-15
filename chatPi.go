package main

import (
	"net"
	"fmt"
	"time"
	"github.com/aki237/chatPi/user"
)

func log(a ...interface{}) {
	fmt.Println(time.Now().String()," ",a)
}

func main() {
	ip := "10.9.10.148"
	log("...chatPi...")
	log("starting server")
	user.NewChat(ip)
	ln, err := net.Listen("tcp", ip+":6672")

	if err != nil {
		log(err)
		return
	}

	// run loop forever (or until ctrl-c)
	for {
		conn, _ := ln.Accept()
		newconn := chatConn{conn}
		go newconn.Serve()
	}
}
