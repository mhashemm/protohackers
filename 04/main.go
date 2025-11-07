package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

func main() {
	db := map[string]string{}

	addr, err := net.ResolveUDPAddr("udp", "0.0.0.0:6969")
	if err != nil {
		panic(err)
	}
	s, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Panic(err)
	}
	defer s.Close()
	for {
		packet := make([]byte, 1000)
		n, c, err := s.ReadFromUDP(packet)
		if err != nil {
			log.Println("ReadFromUDP", err)
			continue
		}
		handler(s, c, string(packet[:n]), db)
	}
}

func handler(s *net.UDPConn, c *net.UDPAddr, msg string, db map[string]string) {
	key, value, isInsert := strings.Cut(msg, "=")
	if isInsert {
		if key != "version" {
			db[key] = value
		}
		return
	}

	if key == "version" {
		res := fmt.Appendf(nil, "version=Ken's Key-Value Store 1.0")
		s.WriteToUDP(res, c)
	}

	res := fmt.Appendf(nil, "%s=%s", key, db[key])
	s.WriteToUDP(res, c)
}
