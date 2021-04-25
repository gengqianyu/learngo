//值接收者vs指针接收者
//值接收者是go语言特有的
//(为什么呢，在其它语言的对象中均有，this，self指针对方法进行调用，
//go方法的定义都是定义在结构体外面，调用的时候就是用接收者去调用, 接收者可以是指针/值均可)
//值/指针接收者均可接收值和指针

package tree

import "fmt"

// 定义一个结构体 利用type,和struct关键字

type Node struct {
	Value       int
	Left, Right *Node
}

// 为结构体定义方法(因为go所有的都是传值，所以想修改结构体Value的值，就利用结构体的指针)

func (node Node) PrintValue() {
	fmt.Println(node.Value)
}

// 和普通函数不同的是，在方法前面定义了一个方法接收者
// 使用指针做为方法的接收者，不然因为值传递，无法修改结构体的值

func (node *Node) SetValue(Value int) {
	if node == nil {
		fmt.Println("node is nil!")
		return
	}
	node.Value = Value
}

// get max node value

func (node *Node) GetMaxNode() int {
	maxNode := 0
	nodesChannel := node.TraverseWithChannel()
	// range 去收channel 只返回一个返回值
	for node := range nodesChannel {
		if maxNode < node.Value {
			maxNode = node.Value
		}
	}
	return maxNode
}

// 节点工厂

func NodeFactory(Value int) *Node {
	return &Node{Value: Value}
}
