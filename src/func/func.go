package main

/*
函数传参
传值vs传引用
传值：函数调用时会对参数进行拷贝，被调用方和调用方两者持有不相关的两份数据；
传引用：函数调用时会传递参数的指针，被调用方和调用方两者持有相同的数据，任意一方做出的修改都会影响另一方。
不同语言会选择不同的方式传递参数，Go 语言选择了传值的方式，无论是传递基本类型、结构体还是指针，都会对传递的参数进行拷贝。
*/
import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"reflect"
	"runtime"
	"strconv"
)

func main() {
	fmt.Println(eval(3, 4, "*"))
	if result, error := eval(3, 4, "x"); error != nil {
		fmt.Println("error:", error)
	} else {
		fmt.Println(result)
	}

	q, r := div(3, 4)
	fmt.Println(q, r)
	//  函数式编程
	fmt.Println(apply(pow, 3, 4))
	fmt.Println(apply(func(a, b int) int {
		return int(math.Pow(float64(a), float64(b)))
	}, 3, 4))
	// 调用可变参数
	h := sum(1, 3, 4, 5, 6, 7, 8)
	fmt.Println(h)
	// 指针参数 通过指针改变变量的值
	a, b := 3, 4
	swap(&a, &b)
	fmt.Println(a, b)
	// 通过函数变化
	a, b = swap2(a, b)
	fmt.Println(a, b)

	// 十进制转二进制
	fmt.Println(convertToBin(18))

	//打印文件内容
	printFile("abc.txt")
	printFile2("abc.txt")

}

// 一般函数
func eval(a, b int, op string) (int, error) {
	switch op {
	case "+":
		return a + b, nil
	case "_":
		return a - b, nil
	case "*":
		return a * b, nil
	case "/":
		q, _ := div(a, b)
		return q, nil
	default:
		//panic("unsupported operation:" + op)
		return 0, fmt.Errorf("unsupported operation:%s", op)
	}
}

// 多返回值函数
func div(a, b int) (q, r int) {
	return a / b, a % b
}

//函数式编程
func apply(op func(a, b int) int, a, b int) int {
	//reflect.ValueOf(op).Pointer() 返回一个uintptr类型
	fmt.Printf("Calling %s with %d,%d\n", runtime.FuncForPC(reflect.ValueOf(op).Pointer()).Name(), a, b)
	return op(a, b)
}

// 重写pow 求幂
func pow(a, b int) int {
	return int(math.Pow(float64(a), float64(b)))
}

// 可变参数
func sum(numbers ...int) int {
	fmt.Printf("类型：%T,值%v\n", numbers, numbers)
	s := 0
	for index := range numbers {
		s += numbers[index]
	}
	return s
}

// 指针变量
func swap(a, b *int) {
	*a, *b = *b, *a
}

func swap2(a, b int) (int, int) {
	return b, a
}

//整数转成二进制 对2取模，商再对二取模 把所有莫反向拼到一起
func convertToBin(v int) string {
	result := ""
	for ; v > 0; v /= 2 {
		// 计算出的莫为整形，转为字符型
		result = strconv.Itoa(v%2) + result
	}
	return result
}

// 打印一个文件内容
func printFile(filename string) {
	// 打开一个文件
	if file, e := os.Open(filename); e != nil {
		panic(e)
	} else {
		// 生成一个扫描器
		scanner := bufio.NewScanner(file)
		fmt.Printf("scanner type:%T value:%v \n", scanner, scanner)
		// 逐行扫描打印文件内容
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}

}

//打印文件2
func printFile2(filename string) {
	if file, e := ioutil.ReadFile(filename); e != nil {
		panic(e)
	} else {
		fmt.Printf("type:%T,value:%s", file, file)
		fmt.Println("----------------------")
		fmt.Println(string(file))
		for _, value := range file {
			fmt.Printf("%s \n", string(value))
		}
	}
}
