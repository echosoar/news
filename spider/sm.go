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
	if spiderName == "" || spiderName == "sm" {
		spiderManager.list = append(spiderManager.list, smSpider)
	}
}

// 神马热搜
func smSpider() []NewsItem {
	url := "https://tophub.today/n/n6YoVqDeZa"
	newsItems := make([]NewsItem, 0)
	req := fasthttp.AcquireRequest()
	req.Header.Set("Host", "tophub.today")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
	req.SetRequestURI(url)
	resp := fasthttp.AcquireResponse()
	err := fasthttp.Do(req, resp)

	status := resp.StatusCode()

	if err != nil || status != fasthttp.StatusOK {
		return newsItems
	}
	body := resp.Body()
	r, _ := regexp.Compile("[\n\r]")
	text := string(r.ReplaceAll(body, []byte{}))
	reg, _ := regexp.Compile("rel=\"nofollow\" itemid=\"\\d+\">([^<]*?)</a>")
	res := reg.FindAllStringSubmatch(text, -1)
	newsItems = make([]NewsItem, 0, len(res))
	for _, matchedItem := range res {
		newsItems = append(newsItems, NewsItem{
			Filter: IsNeedFilter(matchedItem[1], []string{}),
			Title:  utils.FormatTitle(matchedItem[1]),
			Link:   "https://m.sm.cn/s?q=" + matchedItem[1],
			Origin: "神马热搜",
			Time:   time.Now().Unix() - 30*int64(time.Minute.Seconds()), // 偏移30分钟
		})
	}
	return newsItems
}
