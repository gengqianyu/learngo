package model

type Commands []map[string]interface{}

type Command struct {
	Label string
	Data  Commands
}

func CreateCommand() *Command {

	return &Command{
		Label: "go 命令手册",
		Data: []map[string]interface{}{
			{
				"cmd":     "go version",
				"example": "version",
				"flag": []string{
					"nil",
				},
				"desc": "打印Go版本号",
			},
			{
				"cmd":     "go env",
				"example": "go env",
				"flag": []string{
					"nil",
				},
				"desc": "打印Go环境信息",
			},
			{
				"cmd":     "go build [-o output] [-i] [build flags] [packages]",
				"example": "go build -o test -i main.go ",
				"flag": []string{
					"-o 指定输出文件名，只支持在编译单个 package 时使用",
					"-i 安装目标 package 所需的依赖",
					"-a 强行对所有涉及到的代码包进行重新构建，即使它们已经是最新的了",
					"-n 打印编译期间所用到的其它命令，但是并不真正执行它们",
					"-p n 指定编译过程中执行各任务的并行数量,默认等于CPU的逻辑核数",
					"-v 打印出那些被编译的代码包的名字",
					"-work 打印出编译时生成的临时工作目录的路径，并在编译结束时保留它。默认编译结束时会删除该目录。",
					"-x 打印编译期间所用到的其它命令",
				},
				"desc": "编译并安装包和依赖项",
			},
			{
				"cmd":     "go fmt [-n] [-x] [packages]",
				"example": "go fmt main.go",
				"flag": []string{
					"-n 打印编译期间所用到的其它命令，但是并不真正执行它们",
					"-x 标志在执行命令时打印命令",
				},
				"desc": "对 .go 源文件进行格式化",
			},
			{
				"cmd":     "go install [build flags] [packages]",
				"example": "go install -v -x main.go",
				"flag": []string{
					"参见build flag",
				},
				"desc": "编译并安装包和依赖项(默认程序被安装到bin目录下)",
			},
			{
				"cmd":     "go run [build flags] goFiles... [arguments...]",
				"example": "go run -v -x main.go -port 8080",
				"flag": []string{
					"参见build flag",
				},
				"desc": "编译并运行包含命名的Go源文件的主程序包,默认情况直接运行编译后的二进制文件。",
			},
			{
				"cmd":     "go clean [-i] [-r] [-n] [-x] [build flags] [packages]",
				"example": "",
				"flag": []string{
					"-i 删除由 go install 命令所安装的可执行目标文件",
					"-r 删除由 go install 命令所安装的可执行目标文件及其所有依赖的文件",
				},
				"desc": "从软件包源目录中删除目标文件",
			},
			{
				"cmd":     "go get [-d] [-f] [-fix] [-insecure] [-t] [-u] [build flags] [packages]",
				"example": "go get -u github.com/username/project@v2.0",
				"flag": []string{
					"-d 让命令程序只执行下载动作，而不执行安装动作",
					"-u 指示使用网络更新命名包及其依赖项。默认情况下，get使用网络来检查丢失的包，但不使用它来查找现有包的更新",
					"-t 指示还可以下载为指定软件包构建测试所需的软件包",
				},
				"desc": "下载并安装包和依赖项(添加版本号将会升级到指定的版本号version，如果在go module 模式下同时更改go.mod)",
			},
			{
				"cmd":     "go test [-c] [-i] [build and test flags] [packages] [flags for test binary]",
				"example": "go test .",
				"flag": []string{
					"-c 将测试二进制文件编译为pkg.test但不要运行它（其中pkg是软件包导入路径的最后一个元素）可以使用-o标志更改文件名",
					"-i 安装作为测试依赖项的软件包。不要运行测试",
					"-o 将测试二进制文件编译为命名文件。测试仍然运行（除非指定了-c或-i）",
				},
				"desc": "自动测试包",
			},
			{
				"cmd":     "go tool [-n] command [args...]",
				"example": "go tool cover -html=c.out （查看测试覆盖率）",
				"flag": []string{
					"-n 使工具打印将要执行但不执行的命令",
				},
				"desc": "运行指定go的工具",
			},
			{
				"cmd":     "go vet [-n] [-x] [packages]",
				"example": "go vet main.go",
				"flag": []string{
					"-n 使工具打印将要执行但不执行的命令",
					"-x 在执行命令时打印命令",
				},
				"desc": "Go源码静态检查工具",
			},
			{
				"cmd":     "go doc <pkg> <sym>[.<method>]",
				"example": "go doc encoding/json",
				"flag": []string{
					"-c 匹配标识符区分大小写",
					"-u 显示未导出的文档以及导出的符号和方法。",
				},
				"desc": "show documentation for package or symbol",
			},
			{
				"cmd":     "go list [-e] [-f format] [-json] [build flags] [packages]",
				"example": "go list -m -versions github.com/labstack/gommon",
				"flag": []string{
					"-f 标志指定列表的备用格式,默认就可以",
					"-json 标志使软件包数据以JSON格式而不是模板格式打印",
					"-m 使list列出模块而不是软件包",
				},
				"desc": "list packages",
			},
			{
				"cmd":     "go tidy",
				"example": "go tidy",
				"flag": []string{
					"nil",
				},
				"desc": "downland and update package",
			},
			{
				"cmd":     "go mod init package",
				"example": "go tidy",
				"flag": []string{
					"nil",
				},
				"desc": "初始化project go module（会 create go.mod file）",
			},
		},
	}
}

func (c *Command) getData() Commands {
	return c.Data
}
