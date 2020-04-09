package worker

import (
	"crawler/engine"
	"crawler/zhenai/parser"
	"crawler_distributed/config"
	"errors"
	"fmt"
	"log"
)

// defined worker Parser
type SerializedParser struct {
	Name string // function name
	Args interface{}
}

//defined worker request
type Request struct {
	Url    string
	Parser SerializedParser
}

//defined worker ParserResult
type ParserResult struct {
	Items    []engine.Item
	Requests []Request
}

// engine request to worker request
func SerializeRequest(request engine.Request) Request {
	name, args := request.Parser.Serialize()
	return Request{
		Url: request.Url,
		Parser: SerializedParser{
			Name: name,
			Args: args,
		},
	}
}

// worker request to engine request
func DeserializeRequest(request Request) (engine.Request, error) {
	parser, err := DeserializeParser(request.Parser)
	if err != nil {
		return engine.Request{}, err
	}
	return engine.Request{
		Url:    request.Url,
		Parser: parser,
	}, nil
}

func DeserializeParser(serializeParser SerializedParser) (engine.Parser, error) {
	switch serializeParser.Name {
	case config.ParserCityList:
		return engine.NewFuncParser(parser.CityList, serializeParser.Name), nil
	case config.ParserCity:
		return engine.NewFuncParser(parser.City, serializeParser.Name), nil
	case config.ParserProfile:
		if userName, ok := serializeParser.Args.(string); ok {
			return parser.NewProfileParser(parser.Profile, serializeParser.Name, userName), nil
		} else {
			return nil, fmt.Errorf("invalid arg:%v", serializeParser.Args)
		}
	case config.ParserNil:
		return engine.NilParser{}, nil
	default:
		return nil, errors.New("unknown parser name")
	}
}

//engine parserResult to worker parserResult
func SerializeParserResult(result engine.ParserResult) ParserResult {
	var requests []Request
	for _, request := range result.Requests {
		requests = append(requests, SerializeRequest(request))
	}
	return ParserResult{
		Items:    result.Items,
		Requests: requests,
	}
}

// worker parserResult to engine parserResult
func DeserializeParserResult(result ParserResult) engine.ParserResult {
	var requests []engine.Request
	for _, request := range result.Requests {
		engineRequest, err := DeserializeRequest(request)
		if err != nil {
			log.Printf("error deserialize request:%v", err)
			continue
		}
		requests = append(requests, engineRequest)
	}
	return engine.ParserResult{
		Requests: requests,
		Items:    result.Items,
	}
}

//type CrawlService struct{}

//因为 engine.Request内的Parser属性不是一个网络可传递的对象，
// Parser是一个实现了Parser接口的结构体。里面还还结构体的属性还包括函数，
// 所以不能传递
//func (s CrawlService) Process(request engine.Request, result engine.ParserResult) error {
//	return nil
//}

//func (s CrawlService) Process(request Request, result ParserResult) error {
//	return nil
//}
