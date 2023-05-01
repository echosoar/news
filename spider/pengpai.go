package spider

import (
	"os"
	"regexp"

	"github.com/echosoar/news/utils"
	"github.com/valyala/fasthttp"
)

func init() {
	spiderName := os.Getenv("SPIDER")
	if spiderName == "" || spiderName == "pengpai" {
		spiderManager.list = append(spiderManager.list, pengpaiSpider)
	}
}

// 澎湃新闻
func pengpaiSpider() []NewsItem {
	url := "https://m.thepaper.cn/list_page.jsp?&nodeid=25949&isList=1&pageidx=1"
	newsItems := make([]NewsItem, 0)
	status, resp, err := fasthttp.Get(nil, url)
	if err != nil || status != fasthttp.StatusOK {
		return newsItems
	}
	r, _ := regexp.Compile("[\n\r]")
	text := string(r.ReplaceAll(resp, []byte{}))
	reg, _ := regexp.Compile("<a href=\"(newsDetail_.*?)\">(.*?)</a>.*?<a href=\"list.*?</a></span>\\s*<span>(.*?)</span>")
	res := reg.FindAllStringSubmatch(text, -1)
	newsItems = make([]NewsItem, 0, len(res))
	for _, matchedItem := range res {
		if IsNeedFilter(matchedItem[2]) {
			continue
		}
		newsItems = append(newsItems, NewsItem{
			Title:  matchedItem[2],
			Link:   "https://www.thepaper.cn/" + matchedItem[1],
			Origin: "澎湃",
			Time:   utils.FormatTimeAgo(matchedItem[3]),
		})
	}
	return newsItems
}
