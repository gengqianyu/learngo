package main

import (
	"fmt"
	"net"
	"net/rpc/jsonrpc"
	rpcdemo "rpc"
)

func main() {
	//拨号连接 tcp 1234 获取一个tcp连接
	conn, err := net.Dial("tcp", ":1234")
	if err != nil {
		panic(err)
	}

	//用这个tcp连接创建一个jsonrpc客户端
	client := jsonrpc.NewClient(conn)
	args := rpcdemo.Args{
		A: 10,
		B: 0,
	}
	var result float64
	err = client.Call("DemoService.Div", args, &result)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(result)
	}

}
