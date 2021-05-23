package rpcsupport

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

func RunServer(receiver interface{}, host string) error {
	//注册 rpc 服务
	err := rpc.Register(receiver)
	if err != nil {
		return err
	}
	//创建监听者
	listener, err := net.Listen("tcp", host)
	if err != nil {
		return err
	}
	log.Printf("listening on %s", host)
	for {
		// 接收客户端连接
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("rpc err:%v", err)
			continue
		}
		// 开一个 goroutine，用 jsonrpc 去处理连接
		go jsonrpc.ServeConn(conn)
	}
	return nil
}

func FactoryClient(host string) (*rpc.Client, error) {
	// 创建一个 tcp 连接
	conn, err := net.Dial("tcp", host)
	if err != nil {
		return nil, err
	}
	//用四层 tcp 连接创建一个 五层 rpc 应用客户端
	client := jsonrpc.NewClient(conn)
	return client, err
}
