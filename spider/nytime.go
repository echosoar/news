package spider

import (
	"os"
	"regexp"
	"time"

	"github.com/echosoar/news/utils"
	"github.com/valyala/fasthttp"
)

func init() {
	spiderName := os.Getenv("SPIDER")
	if spiderName == "" || spiderName == "nytime" {
		spiderManager.list = append(spiderManager.list, nytimeSpider)
	}
}

func nytimeSpider() []NewsItem {
	newsItems := make([]NewsItem, 0)
	url := "https://cn.nytimes.com/rss/"
	status, resp, err := fasthttp.Get(nil, url)
	if err != nil || status != fasthttp.StatusOK {
		return newsItems
	}

	r, _ := regexp.Compile("[\n\r]")
	text := string(r.ReplaceAll(resp, []byte{}))

	reg, _ := regexp.Compile("<item>\\s*<title>.*?CDATA\\[\\s*(.*?)\\s*\\]\\]>.*?<link>.*?CDATA\\[\\s*(.*?)\\s*\\]\\].*?<pubDate>(.*?)</pubDate>")
	res := reg.FindAllStringSubmatch(text, -1)

	for _, matchedItem := range res {
		if IsNeedFilter(matchedItem[1], []string{}) {
			continue
		}
		t, _ := time.Parse(time.RFC1123Z, matchedItem[3])
		newsItems = append(newsItems, NewsItem{
			Title:  utils.FormatTitle(matchedItem[1]),
			Link:   matchedItem[2],
			Origin: "纽约时报",
			Time:   t.Unix(),
		})
	}

	return newsItems
}
