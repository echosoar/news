package spider

import (
	"os"
	"regexp"

	"github.com/echosoar/news/utils"
	"github.com/valyala/fasthttp"
)

func init() {
	spiderName := os.Getenv("SPIDER")
	if spiderName == "" || spiderName == "ifeng" {
		spiderManager.list = append(spiderManager.list, ifengSpider)
	}
}

func ifengSpider() []NewsItem {
	url := "https://news.ifeng.com/"
	newsItems := make([]NewsItem, 0)
	status, resp, err := fasthttp.Get(nil, url)
	if err != nil || status != fasthttp.StatusOK {
		return newsItems
	}
	r, _ := regexp.Compile("[\n\r]")
	text := string(r.ReplaceAll(resp, []byte{}))
	reg, _ := regexp.Compile("<a href=\"([^\"]*?)\" title=\"([^\"]*?)\".*?><time .*?(\\d+:\\d+)</time>")
	res := reg.FindAllStringSubmatch(text, -1)
	newsItems = make([]NewsItem, 0, len(res))
	for _, matchedItem := range res {
		newsItems = append(newsItems, NewsItem{
			Filter: IsNeedFilter(matchedItem[2], []string{"凤凰"}),
			Title:  utils.FormatTitle(matchedItem[2]),
			Link:   "https:" + matchedItem[1],
			Origin: "凤凰网",
			Time:   utils.FormatTimeHMToUnix(matchedItem[3]),
		})
	}
	return newsItems
}
