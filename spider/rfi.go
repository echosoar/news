package spider

import (
	"os"
	"regexp"
	"time"

	"github.com/valyala/fasthttp"
)

func init() {
	spiderName := os.Getenv("SPIDER")
	if spiderName == "" || spiderName == "rfi" {
		spiderManager.list = append(spiderManager.list, rfiSpider)
	}
}

func rfiSpider() []NewsItem {
	newsItems := make([]NewsItem, 0)
	urls := []string{
		"https://www.rfi.fr/cn/%E6%BB%9A%E5%8A%A8%E6%96%B0%E9%97%BB/",
	}

	for _, url := range urls {
		status, resp, err := fasthttp.Get(nil, url)
		if err != nil || status != fasthttp.StatusOK {
			continue
		}
		r, _ := regexp.Compile("[\n\r]")
		text := string(r.ReplaceAll(resp, []byte{}))

		reg, _ := regexp.Compile("<a href=\"([^\"]*)\" data-article-item-link>.*?</a>.*?<time datetime=\"([^\"]*)\">.*?<p [^>]*>(.*?)</p>\\s*</div>\\s*</a>")
		res := reg.FindAllStringSubmatch(text, -1)
		newsItems = make([]NewsItem, 0, len(res))
		loc, _ := time.LoadLocation("Europe/Berlin")
		for _, matchedItem := range res {
			if IsNeedFilter(matchedItem[3]) {
				continue
			}
			t, _ := time.ParseInLocation(time.RFC3339, matchedItem[2], loc)
			newsItems = append(newsItems, NewsItem{
				Title:  matchedItem[3],
				Link:   "https://www.rfi.fr" + matchedItem[1],
				Origin: "RFI",
				Time:   t.Unix(),
			})
		}
	}

	return newsItems
}
