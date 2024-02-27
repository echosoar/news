package spider

import (
	"encoding/json"
	"os"
	"strconv"

	"github.com/echosoar/news/utils"
	"github.com/valyala/fasthttp"
)

func init() {
	spiderName := os.Getenv("SPIDER")
	if spiderName == "" || spiderName == "infzm" {
		spiderManager.list = append(spiderManager.list, infzmSpider)
	}
}

type infzmJson struct {
	Code int `json:"code"`
	Data struct {
		Contents []struct {
			ID      int    `json:"id"`
			Subject string `json:"subject"`
			Time    string `json:"publish_time"`
		} `json:"contents"`
	} `json:"data"`
}

func infzmSpider() []NewsItem {
	url := "http://www.infzm.com/contents?term_id=2&page=1&format=json"
	newsItems := make([]NewsItem, 0)
	status, resp, err := fasthttp.Get(nil, url)
	if err != nil || status != fasthttp.StatusOK {
		return newsItems
	}
	infzmJsonStruct := infzmJson{}
	json.Unmarshal(resp, &infzmJsonStruct)
	for _, item := range infzmJsonStruct.Data.Contents {
		newsItems = append(newsItems, NewsItem{
			Filter: IsNeedFilter(item.Subject, []string{}),
			Title:  utils.FormatTitle(item.Subject),
			Link:   "http://www.infzm.com/contents/" + strconv.Itoa(item.ID),
			Origin: "南方周末",
			Time:   utils.FormatTimeYMDHMSToUnix(item.Time),
		})
	}
	return newsItems
}
