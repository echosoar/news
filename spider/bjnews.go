package spider

import (
	"encoding/json"
	"os"

	"github.com/echosoar/news/utils"
	"github.com/valyala/fasthttp"
)

func init() {
	spiderName := os.Getenv("SPIDER")
	if spiderName == "" || spiderName == "bjnews" {
		spiderManager.list = append(spiderManager.list, bjnewsSpider)
	}
}

type bjnewsJSONStruct struct {
	Data []struct {
		Title            string `json:"title"`
		PublishTimestamp int64  `json:"publish_timestamp"`
		DetailURL        struct {
			PcURL string `json:"pc_url"`
		} `json:"detail_url"`
	} `json:"data"`
}

func bjnewsSpider() []NewsItem {
	newsItems := make([]NewsItem, 0)
	url := "https://m.bjnews.com.cn/bwnew/index-tj?page=1&size=20&channel_id=101&wz_id=1"
	status, resp, err := fasthttp.Get(nil, url)
	if err != nil || status != fasthttp.StatusOK {
		return newsItems
	}

	bjnewsJSON := bjnewsJSONStruct{}
	json.Unmarshal(resp, &bjnewsJSON)

	for _, item := range bjnewsJSON.Data {
		if IsNeedFilter(item.Title, []string{"新京"}) {
			continue
		}
		newsItems = append(newsItems, NewsItem{
			Title:  utils.FormatTitle(item.Title),
			Link:   item.DetailURL.PcURL,
			Time:   item.PublishTimestamp,
			Origin: "新京报",
		})
	}

	return newsItems
}
