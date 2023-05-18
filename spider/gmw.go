package spider

import (
	"os"
	"regexp"

	"github.com/echosoar/news/utils"
	"github.com/valyala/fasthttp"
)

func init() {
	spiderName := os.Getenv("SPIDER")
	if spiderName == "" || spiderName == "gmw" {
		spiderManager.list = append(spiderManager.list, gmwSpider)
	}
}

func gmwSpider() []NewsItem {
	url := "https://politics.gmw.cn/node_9840.htm"
	newsItems := make([]NewsItem, 0)
	status, resp, err := fasthttp.Get(nil, url)
	if err != nil || status != fasthttp.StatusOK {
		return newsItems
	}
	r, _ := regexp.Compile("[\n\r]")
	text := string(r.ReplaceAll(resp, []byte{}))
	reg, _ := regexp.Compile("<li><a href=([^>]*?) target=_blank>(.*?)</a><span class=\"channel-newsTime\">(.*?)</span></li>")
	res := reg.FindAllStringSubmatch(text, -1)
	newsItems = make([]NewsItem, 0, len(res))
	for _, matchedItem := range res {
		if IsNeedFilter(matchedItem[2], []string{}) {
			continue
		}
		newsItems = append(newsItems, NewsItem{
			Title:  utils.FormatTitle(matchedItem[2]),
			Link:   "https://politics.gmw.cn/" + matchedItem[1],
			Origin: "光明网",
			Time:   utils.GetYMDUnixTime(matchedItem[3]),
		})
	}
	return newsItems
}
