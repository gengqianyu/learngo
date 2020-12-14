package main

import "fmt"

func main() {
	//  创建一个slice 1
	var s1 []int
	if s1 == nil {
		fmt.Println("zero value for slice is nil")
	}
	// 打印0到100所有奇数
	for i := 0; i < 50; i++ {
		//printSlice(s1)
		s1 = append(s1, 2*i+1)
	}
	fmt.Println(s1)

	//  创建切片2
	s2 := []int{3, 5, 7}
	printSlice(s2)

	// 创建切片3 第一个参数为切片类型 第二个参数为len 第三个参数为cap
	s3 := make([]int, 16)
	s4 := make([]int, 16, 32)
	if s3 != nil {
		fmt.Println("s3")
		printSlice(s3)
	}
	printSlice(s4)

	// 拷贝一个slice
	copy(s3, s2)
	printSlice(s3)
	// 删除一个slice元素 把第4个元素删除 第4个元素是0
	s3 = append(s3[:3], s3[4:]...)
	fmt.Println("删除一个元素")
	printSlice(s3)
	// slice 取出/删除第一个元素
	front := s3[0]
	s3 = s3[1:]
	//  slice 取出/删除最后一个元素
	tail := s3[len(s3)-1]
	s3 = s3[:len(s3)-1]
	fmt.Println("front=", front, "tail=", tail)
	printSlice(s3)

}

func printSlice(s []int) {
	fmt.Printf("value=%v,len=%d,cap=%d \n", s, len(s), cap(s))
}
