package main

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"rpc/services"
)

func main() {
	// register a service
	rpc.Register(services.Demo{})
	// start up an tcp service and listen 1234 port
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		panic(err)
	}
	//accept connect from listener and execute jsonRpc method
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("accept error:%v", err)
			continue
		}
		// 不能在这里当场做，会阻塞影响效率 所以得开 goroutine
		go jsonrpc.ServeConn(conn)
	}

}
