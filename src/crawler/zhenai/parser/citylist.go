/*
defined an parser city list (城市列表解析器)

return request and item
*/
package parser

import (
	"crawler/engine"
	"crawler_distributed/config"
	"regexp"
)

const regexpA = `<a href="(http://www.zhenai.com/zhenghun/[a-z0-9]+)"[^>]*>([\p{Han}]+)</a>`

func CityList(contents []byte, _ string) engine.ParserResult {
	result := engine.ParserResult{}

	reg := regexp.MustCompile(regexpA)
	matches := reg.FindAllSubmatch(contents, -1)
	for _, match := range matches {
		result.Requests = append(result.Requests, engine.Request{
			Url:    string(match[1]),
			Parser: engine.NewFuncParser(City, config.ParserCity),
		})
	}
	return result
}
