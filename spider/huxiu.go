package spider

import (
	"encoding/json"
	"os"
	"strconv"

	"github.com/valyala/fasthttp"
)

func init() {
	spiderName := os.Getenv("SPIDER")
	if spiderName == "" || spiderName == "huxiu" {
		spiderManager.list = append(spiderManager.list, huxiuSpider)
	}
}

type huxiuJson struct {
	Data struct {
		DataList []struct {
			Aid   string `json:"aid"`
			Title string `json:"title"`
			Time  string `json:"dateline"`
		} `json:"dataList"`
	} `json:"data"`
}

func huxiuSpider() []NewsItem {
	url := "https://article-api.huxiu.com/web/article/articleList"
	newsItems := make([]NewsItem, 0)
	args := &fasthttp.Args{}
	args.Add("platform", "www")
	args.Add("pagesize", "22")
	status, resp, err := fasthttp.Post(nil, url, args)
	if err != nil || status != fasthttp.StatusOK {
		return newsItems
	}
	huxiuJsonStruct := huxiuJson{}
	json.Unmarshal(resp, &huxiuJsonStruct)
	for _, item := range huxiuJsonStruct.Data.DataList {
		if IsNeedFilter(item.Title) {
			continue
		}
		time, _ := strconv.Atoi(item.Time)
		newsItems = append(newsItems, NewsItem{
			Title:  item.Title,
			Link:   "https://www.huxiu.com/article/" + item.Aid + ".html",
			Origin: "虎嗅",
			Time:   int64(time),
		})
	}
	return newsItems
}
