package spider

import (
	"os"
	"regexp"
	"time"

	"github.com/valyala/fasthttp"
)

func init() {
	spiderName := os.Getenv("SPIDER")
	if spiderName == "" || spiderName == "bbc" {
		spiderManager.list = append(spiderManager.list, bbcSpider)
	}
}

func bbcSpider() []NewsItem {
	newsItems := make([]NewsItem, 0)
	url := "https://feeds.bbci.co.uk/zhongwen/simp/rss.xml"
	status, resp, err := fasthttp.Get(nil, url)
	if err != nil || status != fasthttp.StatusOK {
		return newsItems
	}

	r, _ := regexp.Compile("[\n\r]")
	text := string(r.ReplaceAll(resp, []byte{}))

	reg, _ := regexp.Compile("<item>\\s*<title>.*?CDATA\\[\\s*(.*?)\\s*\\]\\]>.*?<link>(.*?)</link>.*?<pubDate>(.*?)</pubDate>")
	res := reg.FindAllStringSubmatch(text, -1)

	for _, matchedItem := range res {
		if IsNeedFilter(matchedItem[1]) {
			continue
		}
		t, _ := time.Parse(time.RFC1123, matchedItem[3])
		newsItems = append(newsItems, NewsItem{
			Title:  matchedItem[1],
			Link:   matchedItem[2],
			Origin: "BBC",
			Time:   t.Unix(),
		})
	}

	return newsItems
}
