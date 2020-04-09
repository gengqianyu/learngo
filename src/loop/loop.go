package main

import (
	"bufio"
	"fmt"
	goScanner "go/scanner"
	"go/token"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"text/scanner"
)

//十进制转二进制
func convertToBin(n int) string {
	result := ""
	for ; n > 0; n /= 2 {
		lsb := n % 2
		result = strconv.Itoa(lsb) + result
	}
	return result
}

//按行打印文件内容 os.Open的用法
func printFile(filename string) {
	// 打开一个文件
	file, e := os.Open(filename)
	defer file.Close()

	if e != nil {
		panic(e)
	}
	printFileContents(file)
}

// bufio 用法 逐行扫描打印文件内容，
func printFileContents(reader io.Reader) {
	// 生成一个扫描器
	scanner := bufio.NewScanner(reader)
	fmt.Printf("scanner type:%T value:%v \n", scanner, scanner)
	// 逐行扫描打印文件内容，
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}

//按字符打印文件内容 ioutil 用法
func printFile2(filename string) {
	//将文件内容读进一个[]byte slice
	file, e := ioutil.ReadFile(filename)

	if e != nil {
		log.Fatal(e.Error())
	}

	fmt.Printf("type:%T,value:%s", file, file)
	fmt.Println("----------------------")
	fmt.Println(string(file))

	//遍历[]byte print element
	for _, value := range file {
		fmt.Printf("%s \n", string(value))
	}

}

// 按照token的方式打印字符 text/scanner 用法
func printFileForScanner(filename string) {
	// open file
	file, err := os.Open(filename)
	defer file.Close()

	if err != nil {
		log.Fatal(err.Error())
	}
	//将fd包装成reader
	reader := bufio.NewReader(file)

	//defend scanner
	var s scanner.Scanner
	s.Filename = file.Name()
	//file.Fd() 文件句柄
	//初始化 s
	s.Init(reader)

	//这是一个标准的for循环，当初始化一个tok，当满足tok！=scanner.EOF 执行for循环体，然后条件tok变化
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		fmt.Printf("postion:%v,tokenText:%s\n", s.Pos(), s.TokenText())
	}

}

// go/scanner
func printFileForGoScanner(filename string) {
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err.Error())
	}

	// defined scanner
	var s goScanner.Scanner

	// create token file set
	fileSet := token.NewFileSet()

	// init token file
	// AddFile添加一个具有给定文件名、基本偏移量和文件大小的新文件到文件集s并返回该文件。
	// 基本偏移量不能小于fileSet的Base ()，大小不能为负。在特殊情况下，如果提供了负基数，则使用FileSet的base()的当前值。
	file := fileSet.AddFile(filename /*文件名*/, fileSet.Base() /*file的初始偏移*/, len(src) /*文件的大小*/)
	// init scanner
	// Init通过将扫描器设置在src的开头来准备扫描器以标记文本src。 扫描仪使用fileSet文件获取位置信息，并为每行添加行信息。
	// 如果文件大小与src大小不匹配，则Init会引起panic。
	s.Init(file, src, nil, scanner.ScanComments)

	//note 下面的条件是token.EOF
	//扫描会扫描下一个token，并返回token位置，token及其文字字符串（如果适用）。 源结束点由token.EOF指示。
	for pos, tok, lit := s.Scan(); tok != token.EOF; pos, tok, lit = s.Scan() {
		fmt.Printf("postion:%s\t,token:%s\t,lit:%q\n", fileSet.Position(pos), tok, lit)
	}
}

func forever() {
	for {
		fmt.Println("abc")
	}
}

func main() {
	fmt.Println("convertToBin results:")
	fmt.Println(
		convertToBin(5),  // 101
		convertToBin(13), // 1101
		convertToBin(72387885),
		convertToBin(0),
	)

	fmt.Println("abc.txt contents:")
	printFile("abc.txt")
	fmt.Println("------------------------------")
	printFile2("abc.txt")
	fmt.Println("------------------------------")
	printFileForScanner("abc.txt")
	fmt.Println("------------------------------")
	printFileForGoScanner("abc.txt")
	fmt.Println("------------------------------")
	fmt.Println("printing a string:")
	//定义一个跨行的字符串
	s := `ssss "ccccc"
		gengqianyu
		123`
	// 把字符当成一个文件去使用
	printFileContents(strings.NewReader(s))

	// Uncomment to see it runs forever
	// forever()
}
