package spider

import (
	"os"
	"regexp"

	"github.com/echosoar/news/utils"
	"github.com/valyala/fasthttp"
)

func init() {
	spiderName := os.Getenv("SPIDER")
	if spiderName == "" || spiderName == "stdaily" {
		spiderManager.list = append(spiderManager.list, stdailySpider)
	}
}

func stdailySpider() []NewsItem {
	url := "http://www.stdaily.com/index/kejixinwen/kejixinwen.shtml"
	newsItems := make([]NewsItem, 0)
	status, resp, err := fasthttp.Get(nil, url)
	if err != nil || status != fasthttp.StatusOK {
		return newsItems
	}
	r, _ := regexp.Compile("[\n\r]")
	text := string(r.ReplaceAll(resp, []byte{}))
	reg, _ := regexp.Compile("<h3>\\s*<a href=\"(.*?)\"[^>]*>(.*?)</a>\\s*</h3>.*?<span>\\s*(\\d[\\d\\s-:]*\\d)\\s*</span>")
	res := reg.FindAllStringSubmatch(text, -1)
	newsItems = make([]NewsItem, 0, len(res))
	for _, matchedItem := range res {
		newsItems = append(newsItems, NewsItem{
			Filter: IsNeedFilter(matchedItem[2], []string{}),
			Title:  utils.FormatTitle(matchedItem[2]),
			Link:   "http://www.stdaily.com" + matchedItem[1],
			Origin: "科技日报",
			Time:   utils.FormatTimeYMDHMSToUnix(matchedItem[3]),
		})
	}
	return newsItems
}
