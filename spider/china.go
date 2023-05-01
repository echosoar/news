package spider

import (
	"os"
	"regexp"

	"github.com/echosoar/news/utils"
	"github.com/valyala/fasthttp"
)

func init() {
	spiderName := os.Getenv("SPIDER")
	if spiderName == "" || spiderName == "china" {
		spiderManager.list = append(spiderManager.list, chinaSpider)
	}
}

func chinaSpider() []NewsItem {
	newsItems := make([]NewsItem, 0)
	urls := []string{
		"http://news.china.com.cn/node_7247300.htm",
	}

	for _, url := range urls {
		status, resp, err := fasthttp.Get(nil, url)
		if err != nil || status != fasthttp.StatusOK {
			continue
		}
		r, _ := regexp.Compile("[\n\r]")
		text := string(r.ReplaceAll(resp, []byte{}))
		reg, _ := regexp.Compile("<h2><a\\s*href=\"([^\"]*?)\">([^<]*?)</a></h2>\\s*<div class=\"time\">(.*?)</div>")
		res := reg.FindAllStringSubmatch(text, -1)
		newsItems = make([]NewsItem, 0, len(res))
		for _, matchedItem := range res {
			if IsNeedFilter(matchedItem[2]) {
				continue
			}
			newsItems = append(newsItems, NewsItem{
				Title:  matchedItem[2],
				Link:   matchedItem[1],
				Origin: "中国网",
				Time:   utils.FormatTimeYMDToUnix(matchedItem[3]),
			})
		}
	}

	return newsItems
}
