package spider

import (
	"encoding/json"
	"os"
	"strconv"

	"github.com/valyala/fasthttp"
)

func init() {
	spiderName := os.Getenv("SPIDER")
	if spiderName == "" || spiderName == "sinasports" {
		spiderManager.list = append(spiderManager.list, sinaSportsSpider)
	}
}

type sinaSports struct {
	Result struct {
		Data []struct {
			Base struct {
				Base struct {
					Url string `json:"url"`
				} `json:"base"`
			} `json:"base"`
			Info struct {
				Title    string `json:"title"`
				ShowTime string `json:"showTime"`
			} `json:"info"`
		} `json:"data"`
	} `json:"result"`
}

func sinaSportsSpider() []NewsItem {
	url := "https://feeds.sina.cn/api/v4/tianyi?action=0&up=0&down=0&length=15&cre=tianyi&mod=wspt"
	newsItems := make([]NewsItem, 0)
	status, resp, err := fasthttp.Get(nil, url)
	if err != nil || status != fasthttp.StatusOK {
		return newsItems
	}
	sinaSportsStruct := sinaSports{}
	json.Unmarshal(resp, &sinaSportsStruct)
	for _, matchedItem := range sinaSportsStruct.Result.Data {
		if IsNeedFilter(matchedItem.Info.Title) {
			continue
		}
		time, _ := strconv.Atoi(matchedItem.Info.ShowTime)
		newsItems = append(newsItems, NewsItem{
			Title:  matchedItem.Info.Title,
			Link:   matchedItem.Base.Base.Url,
			Origin: "新浪体育",
			Time:   int64(time),
		})
	}
	return newsItems
}
