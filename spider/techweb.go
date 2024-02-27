package spider

import (
	"os"
	"regexp"

	"github.com/echosoar/news/utils"
	"github.com/valyala/fasthttp"
)

func init() {
	spiderName := os.Getenv("SPIDER")
	if spiderName == "" || spiderName == "techWeb" {
		spiderManager.list = append(spiderManager.list, techWebSpider)
	}
}

func techWebSpider() []NewsItem {
	url := "http://www.techweb.com.cn/roll/"
	newsItems := make([]NewsItem, 0)
	req := fasthttp.AcquireRequest()
	req.Header.Set("Host", "www.techweb.com.cn")
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
	reg, _ := regexp.Compile(`<span class="tit"><a href="([^"]*)"\s+target="_blank" title="([^"]*)">.*?<span class="source">\s*TechWeb.com.cn\s*</span>\s*<span class="time">\s*(.*?):\d+\s*</span>`)
	res := reg.FindAllStringSubmatch(text, -1)
	newsItems = make([]NewsItem, 0, len(res))
	for _, matchedItem := range res {
		newsItems = append(newsItems, NewsItem{
			Filter: IsNeedFilter(matchedItem[2], []string{}),
			Title:  utils.FormatTitle(matchedItem[2]),
			Link:   matchedItem[1],
			Origin: "TechWeb",
			Time:   utils.FormatTimemdToUnix(matchedItem[3]),
		})
	}
	return newsItems
}
