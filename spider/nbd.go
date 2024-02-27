package spider

import (
	"os"
	"regexp"

	"github.com/echosoar/news/utils"
	"github.com/valyala/fasthttp"
)

func init() {
	spiderName := os.Getenv("SPIDER")
	if spiderName == "" || spiderName == "nbd" {
		spiderManager.list = append(spiderManager.list, nbdSpider)
	}
}

// 每日财经新闻
func nbdSpider() []NewsItem {
	url := "https://www.nbd.com.cn/columns/3/"
	newsItems := make([]NewsItem, 0)
	status, resp, err := fasthttp.Get(nil, url)
	if err != nil || status != fasthttp.StatusOK {
		return newsItems
	}
	r, _ := regexp.Compile("[\n\r]")
	text := string(r.ReplaceAll(resp, []byte{}))
	today := utils.FormatNow("2006-01-02")
	reg, _ := regexp.Compile("<p class=\"u-channeltime\">\\s*" + today + "\\s*</p>(.*?)</ul>")
	res := reg.FindAllStringSubmatch(text, -1)
	if len(res) == 0 {
		return newsItems
	}
	listHtml := res[0][0]
	itemReg, _ := regexp.Compile("<li class=\"u-news-title\">.*?<a href=\"(.*?)\".*?>(.*?)</a>\\s*<span>\\s*(.*?)\\s*</span>\\s*</li>")
	itemRes := itemReg.FindAllStringSubmatch(listHtml, -1)
	for _, matchedItem := range itemRes {
		newsItems = append(newsItems, NewsItem{
			Filter: IsNeedFilter(matchedItem[2], []string{"每经"}),
			Title:  utils.FormatTitle(matchedItem[2]),
			Link:   matchedItem[1],
			Origin: "每经网",
			Time:   utils.FormatTimeYMDHMSToUnix(today + " " + matchedItem[3]),
		})
	}
	return newsItems
}
