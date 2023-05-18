package spider

import (
	"os"
	"regexp"

	"github.com/echosoar/news/utils"
	"github.com/valyala/fasthttp"
)

func init() {
	spiderName := os.Getenv("SPIDER")
	if spiderName == "" || spiderName == "chinanews" {
		spiderManager.list = append(spiderManager.list, chinaNewsSpider)
	}
}

func chinaNewsSpider() []NewsItem {
	url := "http://www.chinanews.com/scroll-news/news1.html"
	newsItems := make([]NewsItem, 0)
	status, resp, err := fasthttp.Get(nil, url)
	if err != nil || status != fasthttp.StatusOK {
		return newsItems
	}
	r, _ := regexp.Compile("[\n\r]")
	text := string(r.ReplaceAll(resp, []byte{}))
	reg, _ := regexp.Compile("<div class=\"dd_bt\"><a href=\"(.*?)\">(.*?)</a></div><div class=\"dd_time\">(.*?)</div>")
	res := reg.FindAllStringSubmatch(text, -1)
	newsItems = make([]NewsItem, 0, len(res))
	for _, matchedItem := range res {
		if IsNeedFilter(matchedItem[2], []string{}) {
			continue
		}
		newsItems = append(newsItems, NewsItem{
			Title:  utils.FormatTitle(matchedItem[2]),
			Link:   "http://www.chinanews.com/" + matchedItem[1],
			Origin: "中新网",
			Time:   utils.FormatTimemdToUnix(matchedItem[3]),
		})
	}
	return newsItems
}
