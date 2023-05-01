package spider

import (
	"os"
	"regexp"

	"github.com/echosoar/news/utils"
	"github.com/valyala/fasthttp"
)

func init() {
	spiderName := os.Getenv("SPIDER")
	if spiderName == "" || spiderName == "ifanr" {
		spiderManager.list = append(spiderManager.list, ifanrSpider)
	}
}

func ifanrSpider() []NewsItem {
	url := "https://www.ifanr.com/"
	newsItems := make([]NewsItem, 0)
	status, resp, err := fasthttp.Get(nil, url)
	if err != nil || status != fasthttp.StatusOK {
		return newsItems
	}
	r, _ := regexp.Compile("[\n\r]")
	text := string(r.ReplaceAll(resp, []byte{}))
	reg, _ := regexp.Compile("<h3>\\s*<a class=\"js-title-transform\" href=\"https://www.ifanr.com/(\\d+)\" [^>]*ga-action=\"ToItemArticle\" ga-label=\"Article\" target=\"_blank\">(.*?)</a>.*?<time data-time=\"(.*?)\"")
	res := reg.FindAllStringSubmatch(text, -1)
	newsItems = make([]NewsItem, 0, len(res))
	for _, matchedItem := range res {
		if IsNeedFilter(matchedItem[2]) {
			continue
		}
		newsItems = append(newsItems, NewsItem{
			Title:  matchedItem[2],
			Link:   "https://www.ifanr.com/" + matchedItem[1],
			Origin: "爱范儿",
			Time:   utils.FormatTimeYMDHMSToUnix(matchedItem[3]),
		})
	}
	return newsItems
}
