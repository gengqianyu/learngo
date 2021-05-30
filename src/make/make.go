package main

import "fmt"

func main() {
	//make 的作用是初始化内置的数据结构，也就是我们在前面提到的切片、哈希表和 Channel；
	//new 的作用是根据传入的类型分配一片内存空间并返回指向这片内存空间的指针；
	slice := make([]int, 0, 100)   //slice 是一个包含 data、cap 和 len 的结构体 reflect.SliceHeader；
	hash := make(map[int]bool, 10) //hash 是一个指向 runtime.hmap 结构体的指针；
	ch := make(chan int, 5)        //ch 是一个指向 runtime.hchan 结构体的指针；

	fmt.Println(len(slice))
	fmt.Println(len(hash))
	fmt.Println(len(ch))
}
