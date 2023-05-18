package spider

import (
	"encoding/json"
	"os"

	"github.com/echosoar/news/utils"
	"github.com/valyala/fasthttp"
)

func init() {
	spiderName := os.Getenv("SPIDER")
	if spiderName == "" || spiderName == "chengdu" {
		spiderManager.list = append(spiderManager.list, chengduSpider)
	}
}

type ChengduJson struct {
	Lists []struct {
		Url   string `json:"url"`
		Title string `json:"title"`
		Time  string `json:"times"`
	} `json:"lists"`
}

func chengduSpider() []NewsItem {
	url := "http://wap.chengdu.cn/cmstopapi/api_news.php?size=30&page=1"
	newsItems := make([]NewsItem, 0)
	status, resp, err := fasthttp.Get(nil, url)
	if err != nil || status != fasthttp.StatusOK {
		return newsItems
	}
	chengduJsonStruct := ChengduJson{}
	json.Unmarshal(resp, &chengduJsonStruct)
	for _, item := range chengduJsonStruct.Lists {
		if IsNeedFilter(item.Title, []string{}) {
			continue
		}
		newsItems = append(newsItems, NewsItem{
			Title:  utils.FormatTitle(item.Title),
			Link:   item.Url,
			Origin: "红星新闻",
			Time:   utils.FormatTimeAgo(item.Time),
		})
	}
	return newsItems
}
