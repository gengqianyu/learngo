package parser

import (
	"io/ioutil"
	"testing"
)

func TestCityList(t *testing.T) {

	tests := []struct {
		file          string
		expectedCount int
		expectedUrls  []string
	}{
		{
			"cityList_test_data.html",
			470,
			[]string{
				"http://www.zhenai.com/zhenghun/aba",
				"http://www.zhenai.com/zhenghun/akesu",
				"http://www.zhenai.com/zhenghun/alashanmeng",
			},
		},
	}

	for _, test := range tests {
		// io直接读文件
		content, err := ioutil.ReadFile(test.file)
		if err != nil {
			panic(err)
		}
		result := CityList(content, "http://www.zhenai.com/zhenghun")

		if count := len(result.Requests); count != test.expectedCount {
			t.Errorf("expected count:%d,got %d", test.expectedCount, count)
		}

		for i, u := range test.expectedUrls {
			// 如果预期url和实际不符不通过
			if url := result.Requests[i].Url; u != url {
				t.Errorf("expected url %s,got %s", u, url)
			}
		}
	}

}
