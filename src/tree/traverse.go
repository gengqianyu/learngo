/*
封装
名字一般使用CamelCase
首字母大写： public
首字母小写：private

包
每个目录一个包
main包包含一个可执行入口main func
为结构定义的方法必须放在同一个包内，可以是不同的文件

*/
package tree

import "fmt"

// 一个结构体方法可以放到不同文件中，只要在一个包内就可以被调用
// 看起来有点奇怪，其实很简单，就是Node结构体的Traverse方法内，
// 执行的是调用 Node结构体的另外一个方法TraverseFunc，去帮它做遍历
// 和php对象方法中调用其它方法是一样的
func (node *Node) Traverse() {
	// 注意 TraverseFunc的参数是一个函数类型，只要它遵循下面定义的 funcNode类型就行。
	// 它不管函数体内去实现什么
	node.TraverseFunc(func(n *Node) {
		n.PrintValue()
	})
	fmt.Println()
}

//定义一个函数类型
type funcNode func(*Node)

// 为结构体定义一个TraverseFunc方法，只负责遍历二叉树。
// 至于遍历到二叉树结点时，要干什么事情，那就得看f怎么实现了，比如下面实现一个统计
func (node *Node) TraverseFunc(f funcNode) {
	if node == nil {
		return
	}
	node.Left.TraverseFunc(f)
	f(node)
	node.Right.TraverseFunc(f)
}

//为node结构定义一个函数，统计二叉树结点个数
func (node *Node) CountNode() int {
	count := 0
	node.TraverseFunc(func(n *Node) {
		count++
		fmt.Printf("二叉树第%d个结点，value:%d\n", count, n.Value)
	})
	return count
}

func (node *Node) TraverseWithChannel() chan *Node {
	// create an channel
	out := make(chan *Node)
	// open goroutine to work
	go func() {
		node.TraverseFunc(func(n *Node) {
			// send an node to channel
			out <- n
		})
		// data is sent,close channel
		close(out)
	}()

	return out
}
