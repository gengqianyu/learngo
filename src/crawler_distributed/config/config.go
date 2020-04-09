package config

const (
	// parser name
	ParserCityList = "CityList"
	ParserCity     = "City"
	ParserProfile  = "Profile"
	ParserNil      = "NilParser"
	// elasticSearch
	ElasticIndex = "dating_profile"
	// rpc service
	ItemServiceMethod   = "ItemService.Save"
	WorkerServiceMethod = "CrawlService.Process"
	// rate limit
	Qps = 20
)
