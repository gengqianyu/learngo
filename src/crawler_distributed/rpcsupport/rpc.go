package rpcsupport

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

func RunServer(receiver interface{}, host string) error {
	err := rpc.Register(receiver)
	if err != nil {
		return err
	}
	listener, err := net.Listen("tcp", host)
	if err != nil {
		return err
	}
	log.Printf("listening on %s", host)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("rpc err:%v", err)
			continue
		}
		// goroutine
		go jsonrpc.ServeConn(conn)
	}
	return nil
}

func FactoryClient(host string) (*rpc.Client, error) {
	conn, err := net.Dial("tcp", host)
	if err != nil {
		return nil, err
	}
	client := jsonrpc.NewClient(conn)
	return client, err
}
