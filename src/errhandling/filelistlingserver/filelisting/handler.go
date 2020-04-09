package filelisting

import (
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const prefix = "/list/"

// go 语言定义一个string类型也能实现接口，很灵活
type userError string

func (e userError) Error() string {
	return e.Message()
}

// 实现一个Message（） 返回就是它本身字符串
func (e userError) Message() string {
	return string(e)
}

// 处理文件列表函数
func HandleFileList(writer http.ResponseWriter, request *http.Request) error {
	if strings.Index(request.URL.Path, prefix) != 0 {
		return userError("path must start with " + prefix)
	}

	// 获取文件路径
	path := request.URL.Path[len(prefix):]
	// 打一个文件
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	// 关闭文件
	defer file.Close()

	// 获取文件内容
	all, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	// 将文件内容写入 http.ResponseWriter中去显示
	writer.Write(all)
	return nil
}
