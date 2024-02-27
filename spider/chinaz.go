package spider

import (
	"os"
	"regexp"

	"github.com/echosoar/news/utils"
	"github.com/valyala/fasthttp"
)

func init() {
	spiderName := os.Getenv("SPIDER")
	if spiderName == "" || spiderName == "chinaz" {
		spiderManager.list = append(spiderManager.list, chinazSpider)
	}
}

func chinazSpider() []NewsItem {
	url := "https://www.chinaz.com/news/"
	newsItems := make([]NewsItem, 0)
	status, resp, err := fasthttp.Get(nil, url)
	if err != nil || status != fasthttp.StatusOK {
		return newsItems
	}
	r, _ := regexp.Compile("[\n\r]")
	text := string(r.ReplaceAll(resp, []byte{}))
	reg, _ := regexp.Compile(`<div class="info">\s*<div class="info-limit">\s*<h3><a href="([^"]*?)" target="_blank">([^<]*?)</a></h3>.*?<div class="time" title="([^"]*?)"`)
	res := reg.FindAllStringSubmatch(text, -1)
	newsItems = make([]NewsItem, 0, len(res))
	for _, matchedItem := range res {
		newsItems = append(newsItems, NewsItem{
			Title:  utils.FormatTitle(matchedItem[2]),
			Link:   "https:" + matchedItem[1],
			Origin: "站长之家",
			Filter: IsNeedFilter(matchedItem[2], []string{}),
			Time:   utils.FormatTimeYMDHMSToUnix(matchedItem[3]),
		})
	}
	return newsItems
}
