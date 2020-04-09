package worker

import "crawler/engine"

type CrawlService struct{}

func (s *CrawlService) Process(request Request, parserResult *ParserResult) error {
	engineRequest, err := DeserializeRequest(request)
	if err != nil {
		return err
	}
	engineParserResult, err := engine.Worker(engineRequest)
	if err != nil {
		return err
	}
	//给指针指向的地址赋值
	*parserResult = SerializeParserResult(engineParserResult)
	return nil
}
