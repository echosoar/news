package spider

import (
	"os"
	"regexp"

	"github.com/echosoar/news/utils"
	"github.com/valyala/fasthttp"
)

func init() {
	spiderName := os.Getenv("SPIDER")
	if spiderName == "" || spiderName == "guancha" {
		spiderManager.list = append(spiderManager.list, guanchaSpider)
	}
}

func guanchaSpider() []NewsItem {
	newsItems := make([]NewsItem, 0)
	url := "https://www.guancha.cn/GuanChaZheTouTiao/list_1.shtml"
	status, resp, err := fasthttp.Get(nil, url)
	if err != nil || status != fasthttp.StatusOK {
		return newsItems
	}
	r, _ := regexp.Compile("[\n\r]")
	text := string(r.ReplaceAll(resp, []byte{}))
	reg, _ := regexp.Compile("\"content-headline\".*?href=\"(.*?)\" target=\"_blank\">(.*?)<.*?\"interact-comment\".*?<span>(.*?)</span>")
	res := reg.FindAllStringSubmatch(text, -1)
	newsItems = make([]NewsItem, 0, len(res))
	for _, matchedItem := range res {
		if IsNeedFilter(matchedItem[2]) {
			continue
		}
		newsItems = append(newsItems, NewsItem{
			Title:  matchedItem[2],
			Link:   "https://www.guancha.cn" + matchedItem[1],
			Origin: "观察者",
			Time:   utils.FormatTimeYMDHMSToUnix(matchedItem[3]),
		})
	}

	return newsItems
}
