/*
文档
用注释可以写文档
在测试中加入example可以写示例代码
使用go doc 生成文档，使用godoc -http :6060 起一个服务查看文档
*/
package queue

// 别名的形式扩展包
// An FIFO Queue
type Queue []int

// 为Queue slice切片结构定义方法
// pushes the element into the Queue.
func (q *Queue) Push(v int) {
	*q = append(*q, v)
}

// Pops element from head.
// 		e.g. q.Pop()
func (q *Queue) Pop() int {
	head := (*q)[0]
	*q = (*q)[1:]
	return head
}

// Return if the queue isEmpty or not.
func (q *Queue) IsEmpty() bool {
	return len(*q) == 0
}
