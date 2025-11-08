package main

import (
	"bufio"
	"bytes"
	"log"
	"net"
	"regexp"
)

var boguscoinAddresses = regexp.MustCompile(`^7[a-zA-Z0-9]{25,34}$`)

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
		go handler(c)
	}
}

func handler(c net.Conn) {
	defer c.Close()
	upstream, err := net.Dial("tcp", "chat.protohackers.com:16963")
	if err != nil {
		log.Println("upstream", err)
		return
	}
	defer upstream.Close()

	go func() {
		scanner := bufio.NewScanner(upstream)
		for scanner.Scan() {
			msg := scanner.Bytes()

			tokens := bytes.Split(msg, []byte(" "))
			for i, token := range tokens {
				if boguscoinAddresses.Match([]byte(token)) {
					tokens[i] = []byte("7YWHMfk9JZe0LM0g1ZauHuiSxhI")
				}
			}

			msg = bytes.Join(tokens, []byte(" "))
			msg = append(msg, '\n')

			_, err = c.Write(msg)
			if err != nil {
				log.Println("incoming write", err)
				continue
			}
		}
		if scanner.Err() != nil {
			log.Println("scanner.Err", scanner.Err())
			return
		}
		upstream.Close()
	}()

	reader := bufio.NewReader(c)
	for {
		msg, err := reader.ReadBytes('\n')
		if err != nil {
			log.Println("reader", string(msg), err)
			return
		}
		if len(msg) == 0 {
			log.Println("empty msg")
			continue
		}
		if !bytes.HasSuffix(msg, []byte{'\n'}) {
			log.Println("no new line")
			continue
		}
		msg = msg[:len(msg)-1]

		tokens := bytes.Split(msg, []byte(" "))
		for i, token := range tokens {
			if boguscoinAddresses.Match([]byte(token)) {
				tokens[i] = []byte("7YWHMfk9JZe0LM0g1ZauHuiSxhI")
			}
		}

		msg = bytes.Join(tokens, []byte(" "))
		msg = append(msg, '\n')

		_, err = upstream.Write(msg)
		if err != nil {
			log.Println("outgoing write", err)
			continue
		}
	}
}
