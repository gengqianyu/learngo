package parser

import (
	"crawler/Model"
	"crawler/engine"
	"regexp"
	"strconv"
	"strings"
)

const informationRe = `<div class="des f-cl" data-v-3c42fade>([^<]+)</div>`

var ageRe = regexp.MustCompile(`([0-9]+)[\p{Han}]+`)

var heightRe = regexp.MustCompile(`([0-9]+)\w+`)

var basicRe = regexp.MustCompile("basicInfo\":\\[(.+)],\"detailInfo")
var detailRe = regexp.MustCompile("detailInfo\":\\[(.+)].+educationString")

var urlRe = regexp.MustCompile(`album.zhenai.com/u/(.+)`)

func Profile(contents []byte, url string, name string) engine.ParserResult {
	profile := Model.Profile{}
	profile.Name = name

	result := engine.ParserResult{}

	reg := regexp.MustCompile(informationRe)
	matches := reg.FindAllSubmatch(contents, -1)
	for _, match := range matches {
		info := strings.Split(string(match[1]), "|")
		age := strings.TrimSpace(info[1])
		matchAge := ageRe.FindStringSubmatch(age)
		profile.Age = StrToInt(matchAge[1])

		Height := strings.TrimSpace(info[4])
		matchHeight := heightRe.FindStringSubmatch(Height)
		profile.Height = StrToInt(matchHeight[1])

		profile.Birthplace = strings.TrimSpace(info[0])
		profile.Income = strings.TrimSpace(info[5])
		profile.Marriage = strings.TrimSpace(info[3])
		profile.Occupation = strings.TrimSpace(info[2])
	}

	basicInfo := basicRe.FindSubmatch(contents)
	var basicInfos []string
	if len(basicInfo) > 1 {
		basicInfos = strings.Split(string(basicInfo[1]), ",")
	}
	for _, info := range basicInfos {
		profile.BasicInfo = profile.BasicInfo + strings.Trim(info, `"`) + ","
	}

	detailInfo := detailRe.FindSubmatch(contents)
	var detailInfos []string
	if len(detailInfo) > 1 {
		detailInfos = strings.Split(string(detailInfo[1]), ",")
	}
	for _, info := range detailInfos {
		profile.DetailInfo = profile.DetailInfo + strings.Trim(info, `"`) + ","
	}

	urls := urlRe.FindStringSubmatch(url)
	//log.Println(url)
	result.Items = []engine.Item{
		{
			Url:     url,
			Id:      urls[1],
			Type:    "zhenai",
			Payload: profile,
		},
	}
	return result
}

func StrToInt(str string) int {
	i, error := strconv.Atoi(str)
	if error == nil {
		return i
	} else {
		return 0
	}
}

type ProfileFunc func([]byte, string, string) engine.ParserResult

type ProfileParser struct {
	parser   ProfileFunc
	name     string
	UserName string
}

func (p *ProfileParser) Parse(contents []byte, Url string) engine.ParserResult {
	return p.parser(contents, Url, p.UserName)
}

func (p *ProfileParser) Serialize() (name string, args interface{}) {
	return p.name, p.UserName
}

func NewProfileParser(p ProfileFunc, name string, userName string) *ProfileParser {
	return &ProfileParser{
		parser:   p,
		name:     name,
		UserName: userName,
	}
}
