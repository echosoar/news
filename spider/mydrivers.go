package spider

import (
	"os"
	"regexp"

	"github.com/echosoar/news/utils"
	"github.com/valyala/fasthttp"
)

func init() {
	spiderName := os.Getenv("SPIDER")
	if spiderName == "" || spiderName == "mydr" {
		spiderManager.list = append(spiderManager.list, mydrSpider)
	}
}

func mydrSpider() []NewsItem {
	url := "https://news.mydrivers.com/dt.shtml"
	newsItems := make([]NewsItem, 0)
	status, resp, err := fasthttp.Get(nil, url)
	if err != nil || status != fasthttp.StatusOK {
		return newsItems
	}
	r, _ := regexp.Compile("[\n\r]")
	text := string(r.ReplaceAll(resp, []byte{}))
	reg, _ := regexp.Compile(`<a href="([^"]*?)"><span class="titl">([^<]*?)</span><span class="t today">(\d+:\d+)</span></a>`)
	res := reg.FindAllStringSubmatch(text, -1)
	newsItems = make([]NewsItem, 0, len(res))
	for _, matchedItem := range res {
		newsItems = append(newsItems, NewsItem{
			Filter: IsNeedFilter(matchedItem[2], []string{}),
			Title:  utils.FormatTitle(matchedItem[2]),
			Link:   matchedItem[1],
			Origin: "快科技",
			Time:   utils.FormatTimeHMToUnix(matchedItem[3]),
		})
	}
	return newsItems
}
