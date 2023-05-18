package spider

import (
	"os"
	"regexp"

	"github.com/echosoar/news/utils"
	"github.com/valyala/fasthttp"
)

func init() {
	spiderName := os.Getenv("SPIDER")
	if spiderName == "" || spiderName == "163dy" {
		spiderManager.list = append(spiderManager.list, dy163Spider)
	}
}

func dy163Spider() []NewsItem {
	url := "https://dy.163.com/"
	newsItems := make([]NewsItem, 0)
	status, resp, err := fasthttp.Get(nil, url)
	if err != nil || status != fasthttp.StatusOK {
		return newsItems
	}
	r, _ := regexp.Compile("[\n\r]")
	text := string(r.ReplaceAll(resp, []byte{}))
	reg, _ := regexp.Compile("<h3>\\s*<a href=\"([^\"]*)\" title=\"([^\"]*)\">[^>]*</a>\\s*</h3>.*?<div class=\"post_recommend_time\">(.*?)</div>")
	res := reg.FindAllStringSubmatch(text, -1)
	newsItems = make([]NewsItem, 0, len(res))
	for _, matchedItem := range res {
		if IsNeedFilter(matchedItem[2], []string{"网易"}) {
			continue
		}
		newsItems = append(newsItems, NewsItem{
			Title:  utils.FormatTitle(matchedItem[2]),
			Link:   matchedItem[1],
			Origin: "网易号",
			Time:   utils.FormatTimeYMDHMSToUnix(matchedItem[3]),
		})
	}
	return newsItems
}
