package fib

// 定义一个类型
type fib func() int

//定义一个斐波那契数列生成器
func Fibonacci() fib {
	a, b := 0, 1
	return func() int {
		a, b = b, a+b
		return a
	}
}
