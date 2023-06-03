package spider

import (
	"encoding/json"
	"os"

	"github.com/echosoar/news/utils"
	"github.com/valyala/fasthttp"
)

func init() {
	spiderName := os.Getenv("SPIDER")
	if spiderName == "" || spiderName == "wsc" {
		spiderManager.list = append(spiderManager.list, wallstreetcnSpider)
	}
}

type wallstreetcnJson struct {
	Ok   int `json:"ok"`
	Data struct {
		Items []struct {
			Resource struct {
				Title string `json:"title"`
				Url   string `json:"uri"`
				Time  int64  `json:"display_time"`
				IsVIP  bool  `json:"is_in_vip_privilege"`
			} `json:"resource"`
		} `json:"items"`
	} `json:"data"`
}

func wallstreetcnSpider() []NewsItem {
	url := "https://api-one-wscn.awtmt.com/apiv1/content/information-flow?channel=global&accept=article&limit=20&action=upglide"
	newsItems := make([]NewsItem, 0)
	status, resp, err := fasthttp.Get(nil, url)
	if err != nil || status != fasthttp.StatusOK {
		return newsItems
	}

	jsonStruct := wallstreetcnJson{}
	json.Unmarshal(resp, &jsonStruct)

	for _, item := range jsonStruct.Data.Items {
		if item.Resource.IsVIP {
			continue
		}
		if len(item.Resource.Title) < 5 {
			continue
		}
		if IsNeedFilter(item.Resource.Title, []string{}) {
			continue
		}
		newsItems = append(newsItems, NewsItem{
			Title:  utils.FormatTitle(item.Resource.Title),
			Link:   item.Resource.Url,
			Origin: "华尔街见闻",
			Time:   item.Resource.Time,
		})
	}
	return newsItems
}
