package parser

import (
	"crawler/engine"
	"crawler_distributed/config"
	"regexp"
)

var profileRe = regexp.MustCompile(`<a href="(http://album.zhenai.com/u/[0-9]+)"[^>]*>([\p{Han}]+)</a>`)
var cityRe = regexp.MustCompile(`href="(http://www.zhenai.com/zhenghun/[^"]+)"`)

func City(contents []byte, url string) engine.ParserResult {
	result := engine.ParserResult{}
	matches := profileRe.FindAllSubmatch(contents, -1)
	for _, match := range matches {
		//name := string(match[2])
		result.Requests = append(result.Requests, engine.Request{
			Url:    string(match[1]),
			Parser: NewProfileParser(Profile, config.ParserProfile, string(match[2])), //函数调用会传值赋值参数，所以不用上面的name了
		})
	}
	// 匹配城市列表页的城市连接
	matches = cityRe.FindAllSubmatch(contents, -1)
	for _, match := range matches {
		result.Requests = append(result.Requests, engine.Request{
			Url:    string(match[1]),
			Parser: engine.NewFuncParser(City, config.ParserCity),
		})
	}

	return result
}

func GetParserProfile(name string) engine.ParserFunc {
	return func(bytes []byte, url string) engine.ParserResult {
		return Profile(bytes, url, name)
	}
}
