package spider

import (
	"os"
	"regexp"
	"strings"

	"github.com/echosoar/news/utils"
	"github.com/valyala/fasthttp"
)

func init() {
	spiderName := os.Getenv("SPIDER")
	if spiderName == "" || spiderName == "zaobao" {
		spiderManager.list = append(spiderManager.list, zaobaoSpider)
	}
}

func zaobaoSpider() []NewsItem {
	newsItems := make([]NewsItem, 0)
	urls := []string{
		"https://www.zaobao.com/realtime/world",
		"https://www.zaobao.com/realtime/china",
	}

	for _, url := range urls {
		status, resp, err := fasthttp.Get(nil, url)
		if err != nil || status != fasthttp.StatusOK {
			continue
		}
		r, _ := regexp.Compile("[\n\r]")
		text := string(r.ReplaceAll(resp, []byte{}))

		reg, _ := regexp.Compile("<a class=\"article-type-link\"\\s*href=\"(/realtime[^\"]*?)\"><h2>(.*?)</h2></a>\\s+<div.*?-date\">\\s*([^>]*?)\\s*</span>")
		res := reg.FindAllStringSubmatch(text, -1)
		newsItems = make([]NewsItem, 0, len(res))
		for _, matchedItem := range res {
			if IsNeedFilter(matchedItem[2], []string{}) {
				continue
			}
			var time int64
			if strings.HasSuffix(matchedItem[3], "前") {
				time = utils.FormatTimeAgo(matchedItem[3])
			} else {
				time = utils.FormatTimeByFormatToUnix(matchedItem[3], "02/01/2006")
			}

			newsItems = append(newsItems, NewsItem{
				Title:  utils.FormatTitle(matchedItem[2]),
				Link:   "https://www.zaobao.com" + matchedItem[1],
				Origin: "联合早报",
				Time:   time,
			})
		}
	}

	return newsItems
}
