/*
安装 gopm 工具再安装go imports
install gopm
go get -v github.com/gpmgo/gopm

install goimports
gopm -g(把包安装到go path 目录下) -u(已安装更新包使用) -v golang.org/x/tools/cmd/goimports
再 gopath 目录下执行 go install golang.org/x/tools/cmd/goimports 将 goimports 安装到 bin 目录
setting->tools->file watchers click + add
*/
package main

import (
	"fmt"
	"tree"
)

//用组合的方式扩展包
type myTreeNode struct {
	node *tree.Node
}

func (myNode *myTreeNode) postOrder() {
	if myNode == nil || myNode.node == nil {
		return
	}
	// 只有包装成myTreeNode结构体 才能调用postOrder方法
	left := myTreeNode{myNode.node.Left}
	right := myTreeNode{myNode.node.Right}
	left.postOrder()
	right.postOrder()
	myNode.node.PrintValue()
}

func main() {
	// 创建结构体1
	var root tree.Node
	// 创建结构体2
	//root := tree.Node{}
	// 无论是地址还是结构本身，一律使用.来访问成员
	root = tree.Node{Value: 3}
	root.Left = &tree.Node{}
	//root.Right = &tree.Node{5, nil, nil}
	root.Right = &tree.Node{Value: 5}

	// 使用new函数创建结构体
	root.Right.Left = new(tree.Node)
	// 使用工厂模式创建
	root.Left.Right = tree.NodeFactory(2)
	root.Right.Left.SetValue(4)
	// 遍历一个树形结构体
	root.Traverse()
	//统计二叉树节点个数
	nodeCount := root.CountNode()
	fmt.Println("Node count:", nodeCount)

	fmt.Println()
	//
	//myRoot := myTreeNode{&root}
	//myRoot.postOrder()

	maxNode := root.GetMaxNode()
	fmt.Println("the max node value is ", maxNode)

}
