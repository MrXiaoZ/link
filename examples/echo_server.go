package main

import (
	"flag"
	"fmt"
	"github.com/funny/binary"
	"github.com/funny/link"
	_ "github.com/funny/unitest"
)

var (
	addr  = flag.String("addr", ":10010", "echo server address")
	bench = flag.Bool("bench", false, "is benchmark server")
)

func main() {
	flag.Parse()

	server, err := link.Serve("tcp", *addr)
	if err != nil {
		panic(err)
	}
	fmt.Println("server start:", server.Listener().Addr().String())

	server.Serve(func(session *link.Session) {
		addr := session.Conn().RemoteAddr().String()
		log(addr, "connected")
		for {
			var msg Message
			if err := session.Receive(&msg); err != nil {
				break
			}
			log(addr, "say:", string(msg))
			session.Send(msg)
		}
		log(addr, "closed")
	})
}

type Message []byte

func (msg Message) Send(conn *binary.Writer) error {
	conn.WritePacket(msg, binary.SplitByUint16BE)
	return nil
}

func (msg *Message) Receive(conn *binary.Reader) error {
	*msg = conn.ReadPacket(binary.SplitByUint16BE)
	return nil
}

func log(v ...interface{}) {
	if !*bench {
		fmt.Println(v...)
	}
}
