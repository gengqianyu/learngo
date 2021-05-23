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
	if err != nil {
		panic(err)
	}

	//用这个四层 tcp 连接创建一个 五层应用 rpc 客户端
	client := jsonrpc.NewClient(conn)
	defer client.Close()
	args := services.Args{
		A: 10,
		B: 2,
	}
	var result float64
	//执行 rpc 远程调用
	err = client.Call("Demo.Div", args, &result)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(result)
	}

}
