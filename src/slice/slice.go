package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
)

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

	s3 := arr[2:]
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

	var a = []int{1, 2, 3}
	a = append([]int{0}, a...)          // 在开头添加1个元素
	a = append([]int{-3, -2, -1}, a...) // 在开头添加1个切片
	i := 1
	a = append(a[:i], append([]int{0}, a[i:]...)...)       // 在第i个位置插入0
	a = append(a[:i], append([]int{4, 5, 6}, a[i:]...)...) // 在第i个位置插入切片

	a = append(a[:0], a[1:]...) // 删除开头1个元素
	//a = append(a[:0], a[N:]...) // 删除开头N个元素
	a = append(a[:i], a[i+1:]...) // 删除中间1个元素
	//a = append(a[:i], a[i+N:]...) // 删除中间N个元素
}

func updateSlice(s []int) {
	s[0] = 100
}

//删除[]byte中的空格
func TrimSpace(s []byte) []byte {
	b := s[:0]
	for _, x := range s {
		if x != ' ' {
			b = append(b, x)
		}
	}
	return b
}

func Filter(s []byte, fn func(x byte) bool) []byte {
	b := s[:0]
	for _, x := range s {
		if !fn(x) {
			b = append(b, x)
		}
	}
	return b
}

//这段代码返回的[]byte指向保存整个文件的数组。因为切片引用了整个原始数组，导致自动垃圾回收器不能及时释放底层数组的空间。
//一个小的需求可能导致需要长时间保存整个文件数据。这虽然这并不是传统意义上的内存泄漏，但是可能会拖慢系统的整体性能
func FindPhoneNumber(filename string) []byte {
	b, _ := ioutil.ReadFile(filename)
	return regexp.MustCompile("[0-9]+").Find(b)
}

//数据的传值是Go语言编程的一个哲学，虽然传值有一定的代价，但是换取好处是切断了对原始数据的依赖
func FindPhoneNumber2(filename string) []byte {
	b, _ := ioutil.ReadFile(filename)
	b = regexp.MustCompile("[0-9]+").Find(b)
	return append([]byte{}, b...)
}
