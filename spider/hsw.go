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
	if spiderName == "" || spiderName == "hsw" {
		spiderManager.list = append(spiderManager.list, hswSpider)
	}
}

func hswSpider() []NewsItem {
	url := "http://news.hsw.cn/shhot/"
	newsItems := make([]NewsItem, 0)
	req := fasthttp.AcquireRequest()
	// req.Header.Set("Host", "www.techweb.com.cn")
	// req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
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
	reg, _ := regexp.Compile(`<h3>\s*<a href="([^"]*)">([^<]*)</a>\s*</h3>\s*</div>\s*<div class="news_tag">\s*<span class="time">今天</span>`)
	res := reg.FindAllStringSubmatch(text, -1)
	newsItems = make([]NewsItem, 0, len(res))
	for _, matchedItem := range res {
		newsItems = append(newsItems, NewsItem{
			Title:  utils.FormatTitle(matchedItem[2]),
			Link:   matchedItem[1],
			Origin: "华商网",
			Filter: IsNeedFilter(matchedItem[2], []string{}),
			Time:   time.Now().Unix() - 30*int64(time.Minute.Seconds()), // 偏移30分钟
		})
	}
	return newsItems
}
