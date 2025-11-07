package main

import (
	"bufio"
	"fmt"
	"log"
	"maps"
	"net"
	"regexp"
	"slices"
	"strings"
	"sync"
)

var nameValidation = regexp.MustCompile(`^[a-zA-Z0-9]+$`)

func main() {
	clients := map[string]net.Conn{}
	mu := sync.RWMutex{}

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
		go handler(c, clients, &mu)
	}
}

func handler(c net.Conn, clients map[string]net.Conn, mu *sync.RWMutex) {
	defer c.Close()

	_, err := c.Write([]byte("Welcome to budgetchat! What shall I call you?\n"))
	if err != nil {
		log.Println("write welcome", err)
		return
	}

	scanner := bufio.NewScanner(c)
	firstMsg := scanner.Scan()
	if !firstMsg {
		log.Println("firstMsg", scanner.Err())
		return
	}
	name := scanner.Text()
	if !nameValidation.MatchString(name) {
		log.Println("invalid name", name)
		return
	}

	mu.Lock()
	_, exists := clients[name]
	if exists {
		log.Println("duplicate", name)
		return
	}
	clientNames := slices.Collect(maps.Keys(clients))
	clients[name] = c
	mu.Unlock()

	defer func() {
		mu.Lock()
		delete(clients, name)

		left := fmt.Appendf(nil, "* %s has left the room\n", name)
		for _, client := range clients {
			client.Write(left)
		}
		mu.Unlock()
	}()

	_, err = fmt.Fprintf(c, "* The room contains: %s\n", strings.Join(clientNames, ", "))
	if err != nil {
		log.Println("presence", name)
		return
	}

	mu.RLock()
	joined := fmt.Appendf(nil, "* %s has entered the room\n", name)
	for clientName, client := range clients {
		if clientName == name {
			continue
		}
		client.Write(joined)
	}
	mu.RUnlock()

	for scanner.Scan() {
		msg := scanner.Text()
		broadcastMsg := fmt.Appendf(nil, "[%s] %s\n", name, msg)

		mu.RLock()
		for clientName, client := range clients {
			if clientName == name {
				continue
			}
			client.Write(broadcastMsg)
		}
		mu.RUnlock()
	}

	if scanner.Err() != nil {
		log.Println("scanner.Err", scanner.Err())
		return
	}
}
