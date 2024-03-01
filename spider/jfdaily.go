package spider

import (
	"encoding/json"
	"os"
	"strconv"
	"time"

	"github.com/echosoar/news/utils"
	"github.com/valyala/fasthttp"
)

func init() {
	spiderName := os.Getenv("SPIDER")
	if spiderName == "" || spiderName == "jfdaily" {
		spiderManager.list = append(spiderManager.list, jfdailySpider)
	}
}

type jfdailyJson struct {
	DataList []struct {
		Id    int    `json:"id"`
		Sid   int    `json:"sectionid"`
		Title string `json:"title"`
		Time  int64  `json:"publishtime"`
	} `json:"object"`
}

func jfdailySpider() []NewsItem {
	timeStr := strconv.Itoa(int(time.Now().Unix())) + "000"
	url := "https://www.jfdaily.com/news/homeMoreNews?v=" + timeStr
	newsItems := make([]NewsItem, 0)
	args := &fasthttp.Args{}
	args.Add("page", "1")
	args.Add("lastpublishtime", timeStr)
	status, resp, err := fasthttp.Post(nil, url, args)
	if err != nil || status != fasthttp.StatusOK {
		return newsItems
	}
	jsonStruct := jfdailyJson{}
	json.Unmarshal(resp, &jsonStruct)
	for _, item := range jsonStruct.DataList {
		newsItems = append(newsItems, NewsItem{
			Filter: IsNeedFilter(item.Title, []string{}),
			Title:  utils.FormatTitle(item.Title),
			Link:   "https://www.jfdaily.com/staticsg/res/html/web/newsDetail.html?id=" + strconv.Itoa(item.Id) + "&sid=" + strconv.Itoa(item.Sid),
			Origin: "上观",
			Time:   item.Time / 1000,
		})
	}
	return newsItems
}
