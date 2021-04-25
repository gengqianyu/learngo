package main

import (
	"fmt"
	"math"
	"math/cmplx"
)

var (
	aa = 3
	ss = "kkk"
	bb = true
)

func variableZeroValue() {
	var a int
	var s string
	fmt.Printf("%d %q\n", a, s)
}

func variableInitialValue() {
	var a, b int = 3, 4
	var s string = "abc"
	fmt.Println(a, b, s)
}

func variableTypeDeduction() {
	var a, b, c, s = 3, 4, true, "def"
	fmt.Println(a, b, c, s)
}

func variableShorter() {
	a, b, c, s := 3, 4, true, "def"
	b = 5
	fmt.Println(a, b, c, s)
}

func euler() {
	fmt.Printf("%.3f\n",
		cmplx.Exp(1i*math.Pi)+1)
}

// 三角形
func triangle() {
	var a, b int = 3, 4
	fmt.Println(calcTriangle(a, b))
}

// 计算三角形
func calcTriangle(a, b int) int {
	var c int
	c = int(math.Sqrt(float64(a*a + b*b)))
	return c
}

func setConst() {
	const (
		filename = "abc.txt"
		a, b     = 3, 4
	)
	var c int
	c = int(math.Sqrt(a*a + b*b))
	fmt.Println(filename, c)
}

//枚举
// iota 是 golang 语言的常量计数器,只能在常量的表达式中使用。
//iota 在 const关键字出现时将被重置为 0(const内部的第一行之前)，
//const 中每新增一行常量声明将使 iota 计数一次( iota 可理解为 const 语句块中的行索引)。

func enums() {
	const (
		cpp = iota
		_
		python
		golang
		javascript
	)

	const (
		// 1	1	1	1	1	1	1	1
		// 128	64	32	16	8	4	2	1
		// 1 向左移动 (10 * iota) 位
		b = 1 << (10 * iota)
		kb
		mb
		gb
		tb
		pb
	)

	const (
		EOF = -(iota + 1)
		Ident
		Int
		Float
		Char
		String
		RawString
		Comment
	)
	const (
		ScanIdents     = 1 << -Ident
		ScanInts       = 1 << -Int
		ScanFloats     = 1 << -Float // includes Ints
		ScanChars      = 1 << -Char
		ScanStrings    = 1 << -String
		ScanRawStrings = 1 << -RawString
		ScanComments   = 1 << -Comment
		//SkipComments   = 1 << -skipComment // if set with ScanComments, comments become white space
		// GoTokens       = ScanIdents | ScanFloats | ScanChars | ScanStrings | ScanRawStrings | ScanComments | SkipComments
	)

	fmt.Println(cpp, python, golang, javascript)
	fmt.Println(cpp, javascript, python, golang)
	fmt.Println(b, kb, mb, gb, tb, pb)
	fmt.Println(EOF, Ident, Int, Float, Char, String, RawString, Comment)
	fmt.Println(ScanIdents, ScanInts, ScanFloats, ScanChars, ScanStrings, ScanRawStrings, ScanComments)
}

func main() {
	fmt.Println("Hello world")
	variableZeroValue()
	variableInitialValue()
	variableTypeDeduction()
	variableShorter()
	fmt.Println(aa, ss, bb)

	euler()
	triangle()
	setConst()
	fmt.Println()
	enums()
	s := NewSlice()

	defer s.Add(1).Add(2)
	s.Add(3)
	fmt.Println(s)
}

type Slice []int

func NewSlice() Slice {
	return make(Slice, 0)
}

func (s *Slice) Add(element int) *Slice {
	*s = append(*s, element)
	fmt.Println(element)
	return s

}
