package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net"
	"sync/atomic"
)

var counter atomic.Int32

type Request struct {
	Method string   `json:"method"`
	Number *float64 `json:"number"`
}

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
		counter.Add(1)
		fmt.Println(counter.Load())
		go handler(c)
	}
}

func handler(c net.Conn) {
	defer counter.Add(-1)
	defer c.Close()

	s := bufio.NewScanner(c)
	for s.Scan() {
		msg := s.Bytes()
		req := Request{}
		err := json.Unmarshal(msg, &req)
		if err != nil || req.Number == nil || req.Method != "isPrime" {
			log.Print(string(msg), err)
			c.Write([]byte("momo"))
			break
		}

		res := []byte(fmt.Sprintf("{\"method\":\"isPrime\",\"prime\":%t}\n", big.NewInt(int64(*req.Number)).ProbablyPrime(0)))
		_, err = c.Write(res)
		if err != nil {
			log.Print(err)
		}
	}
}
