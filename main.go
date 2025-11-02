package main

import (
	"bytes"
	"cmp"
	"encoding/binary"
	"io"
	"log"
	"math/rand"
	"net"
	"slices"
	"sync"
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
		go handler(c)
	}
}

func handler(c net.Conn) {
	defer c.Close()
	id := rand.Int()

	prices := [][2]int32{}
	mu := sync.Mutex{}

	for {
		msg := make([]byte, 9)
		_, err := io.ReadFull(c, msg)
		if err != nil {
			log.Println("read ", err, id)
			return
		}

		q := msg[0]
		var left, right int32
		err = binary.Read(bytes.NewReader(msg[1:5]), binary.BigEndian, &left)
		if err != nil {
			log.Println("left ", err, id)
			continue
		}
		err = binary.Read(bytes.NewReader(msg[5:]), binary.BigEndian, &right)
		if err != nil {
			log.Println("right ", err, id)
			return
		}
		log.Println("got message ", string(q), left, right, id)

		switch q {
		case 'I':
			mu.Lock()
			prices = append(prices, [2]int32{left, right})
			mu.Unlock()
		case 'Q':
			mu.Lock()
			slices.SortFunc(prices, func(a, b [2]int32) int { return cmp.Compare(a[0], b[0]) })
			start := slices.IndexFunc(prices, func(v [2]int32) bool { return v[0] >= left && v[0] <= right })
			if start == -1 {
				c.Write([]byte{0, 0, 0, 0})
				log.Println("not found ", err, left, right, id)
				mu.Unlock()
				continue
			}

			counter := 0
			sum := int64(0)
			for len(prices) > start && prices[start][0] <= right {
				counter++
				sum += int64(prices[start][1])
				start++
			}
			if counter == 0 {
				c.Write([]byte{0, 0, 0, 0})
				log.Println("zero counter ", err, left, right, id)
				mu.Unlock()
				continue
			}
			mean := int32(sum / int64(counter))
			binary.Write(c, binary.BigEndian, mean)
			log.Println("found ", err, left, right, mean, id)
			mu.Unlock()
		default:
			c.Write([]byte{0, 0, 0, 0})
			log.Println("wrong char ", string(q), id)
		}
	}
}
