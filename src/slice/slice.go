package main

import "fmt"

func main() {
	// slice 是可以向后扩展，但不可以向前扩展

	// 定义一个数组
	arr := [...]int{1, 2, 3, 4, 5, 6, 7}
	//索引          0  1  2  3  4  5  6
	// 截取一个切片
	s1 := arr[2:6]
	//索引                0  1  2  3  4(不可用下标，容量占位)
	fmt.Printf("s1=%v,len(s1)=%d,cap(s1)=%d \n", s1, len(s1), cap(s1))
	// 截取一个切片
	s2 := s1[3:5]
	// 索引                        0  1
	fmt.Printf("s2=%v,len(s2)=%d,cap(s2)=%d \n", s2, len(s2), cap(s2))
	// 报错 s[i]取元素 索引i不可以超过len(s)，向后扩展不可以超过底层数组cap(s)
	//s3 := s1[3:6] //这里的6超过了cap(s1)
	// 索引                        0  1 超出索引x

	s3 := arr [2:]
	fmt.Println(s3)
	s4 := arr[:6]
	fmt.Println(s4)
	s5 := arr[:]
	fmt.Println(s5)
	// 引用类型
	fmt.Println("after updateSlice")
	updateSlice(s2)
	fmt.Println(s2)
	fmt.Println(arr)
	// reSlice
	fmt.Println("reSlice")
	fmt.Println(s1)
	s1 = s1[2:5]
	fmt.Println(s1)
	s1 = s1[2:3]
	fmt.Println(s1)

	//添加slice元素
	fmt.Printf("s4=%v,地址:%p ,容量:%d\n", s4, s4, cap(s4))
	s4 = append(s4, 8)
	fmt.Printf("s4=%v,地址:%p ,容量:%d\n", s4, s4, cap(s4))
	//slice append 元素时 如果超越cap，系统会重新分配更大的底层数组 添加9的时候超过cap所以重新分配了地址
	s4 = append(s4, 9)
	fmt.Printf("s4=%v,地址:%p ,容量:%d\n", s4, s4, cap(s4))
	s6 := append(s4, 10)
	s7 := append(s6, 11)
	s8 := append(s7, 12)
	fmt.Println("s6=", s6, "s7=", s7, "s8=", s8)
	// s6 s7 s8 都没引用arr 只有s4 引用了一次
	fmt.Println("arr=", arr)

}

func updateSlice(s []int) {
	s[0] = 100
}
