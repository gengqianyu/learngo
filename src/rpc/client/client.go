package main

import (
	"fmt"
	"net"
	"net/rpc/jsonrpc"
	"rpc/services"
)

func main() {
	//拨号连接 tcp 1234 获取一个 tcp 连接
	conn, err := net.Dial("tcp", ":1234")
	defer conn.Close()

	if err != nil {
		panic(err)
	}

	//用这个四层 tcp connection 创建一个 五层 rpc 应用客户端
	client := jsonrpc.NewClient(conn)
	defer client.Close()
	args := services.Args{
		A: 10,
		B: 2,
	}
	var reply float64
	//执行 rpc 远程调用
	err = client.Call("Demo.Div", args, &reply)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(reply)
	}

}
