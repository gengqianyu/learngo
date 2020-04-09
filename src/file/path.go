package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	var tree []string
	tree, err := getFileList(`src/`, &tree)
	if err != nil {
		panic(err)
	}

	for _, e := range tree {
		fmt.Println(e)
	}
}

//递归获取目录内容
//关键在于使用了指针
func getFileList(path string, dir *[]string) ([]string, error) {
	fileInfo, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, file := range fileInfo {
		if file.IsDir() {
			//这里注意一定要重新赋值变量
			subPath := path + file.Name() + "/"
			//fmt.Println(subPath)
			*dir = append(*dir, subPath)
			//这里的dir 其实也是值传递，go语言所有的函数传参都是传值。奥妙在于复制指针，还是指向的同一个slice。
			//所以每次递归都能改变原始slice的值
			getFileList(subPath, dir)
		} else {
			//fmt.Println(file.Name())
			name := file.Name()
			*dir = append(*dir, name)
		}

	}
	return *dir, nil
}
