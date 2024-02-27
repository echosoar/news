package spider

import (
	"os"
	"regexp"

	"github.com/echosoar/news/utils"
	"github.com/valyala/fasthttp"
)

func init() {
	spiderName := os.Getenv("SPIDER")
	if spiderName == "" || spiderName == "xinhua" {
		spiderManager.list = append(spiderManager.list, xinHuaSpider)
	}
}

func xinHuaSpider() []NewsItem {
	url := "http://m.news.cn/"
	newsItems := make([]NewsItem, 0)
	status, resp, err := fasthttp.Get(nil, url)
	if err != nil || status != fasthttp.StatusOK {
		return newsItems
	}
	r, _ := regexp.Compile("[\n\r]")
	text := string(r.ReplaceAll(resp, []byte{}))
	reg, _ := regexp.Compile("<a href=\"http://www.news.cn/politics/(\\d{4})-(\\d+)/(\\d+)/([^\"]*?)\" target=\"_blank\">([^<]*?)</a>")
	res := reg.FindAllStringSubmatch(text, -1)
	newsItems = make([]NewsItem, 0, len(res))
	for _, matchedItem := range res {
		newsItems = append(newsItems, NewsItem{
			Filter: IsNeedFilter(matchedItem[5], []string{"新华"}),
			Title:  utils.FormatTitle(matchedItem[5]),
			Link:   "http://www.news.cn/politics/" + matchedItem[1] + "-" + matchedItem[2] + "/" + matchedItem[3] + "/" + matchedItem[4],
			Origin: "新华网",
			Time:   utils.GetYMDUnixTime(matchedItem[1] + "-" + matchedItem[2] + "-" + matchedItem[3]),
		})
	}
	return newsItems
}
