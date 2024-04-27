package main

import (
	"io"
	"log"
	"net"
)

func main() {
	s, err := net.Listen("tcp", "0.0.0.0:6969")
	if err != nil {
		log.Panic(err)
	}
	defer s.Close()

	for {
		c, err := s.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		handler(c)
	}

}

func handler(c net.Conn) {
	msg, err := io.ReadAll(c)
	if err != nil {
		log.Print(c.RemoteAddr().String(), err)
	}
	c.Write(msg)
	err = c.Close()
	if err != nil {
		log.Print(c.RemoteAddr().String(), err)
	}
}
