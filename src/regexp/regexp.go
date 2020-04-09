package main

import (
	"fmt"
	"regexp"
)

func main() {
	text := `my email is 18612314865@163.com
			 email is gengqianyu@gmail.com
				paopao@126.com.cn`

	// 判断是否匹配
	//if match, err := regexp.MatchString(`18612314865@163.com`, text); err != nil {
	//	fmt.Println(match)
	//}

	// 规定正则表达式
	re := regexp.MustCompile(`([\w]+)@([\w]+)\.([\w.]+)`)
	// 单匹配
	//match := re.FindString(text)
	//fmt.Println(match)
	// 多匹配 return []string
	match := re.FindAllString(text, -1)
	fmt.Println(match)
	// 获取所有子匹配 return [][]string
	matchSub := re.FindAllStringSubmatch(text, -1)
	for _, sub := range matchSub {
		fmt.Println(sub)
		for _, mh := range sub {
			fmt.Println(mh)
		}
	}
	fmt.Println(matchSub)

	text = `"basicInfo":["离异","37岁","天秤座(09.23-10.22)","165cm","65kg","工作地:阿坝马尔康","月收入:3-5千","其他职业","大学本科"]`
	basicRe := regexp.MustCompile("basicInfo\":\\[(.+)]")
	matches := basicRe.FindStringSubmatch(text)
	fmt.Printf("%s", matches[1])

	text = `detailInfo":["汉族","籍贯:四川眉山","体型:比较胖","不吸烟","社交场合会喝酒","住在单位宿舍","未买车","有孩子且住在一起","是否想要孩子:不想要孩子","何时结婚:时机成熟就结婚"],"educationString"`
	detailRe := regexp.MustCompile("detailInfo\":\\[(.+)].+educationString")
	matches = detailRe.FindStringSubmatch(text)
	fmt.Printf("%s", matches[1])
}
