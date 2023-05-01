package spider

import (
	"encoding/json"
	"os"

	"github.com/echosoar/news/utils"
	"github.com/valyala/fasthttp"
)

func init() {
	spiderName := os.Getenv("SPIDER")
	if spiderName == "" || spiderName == "abc" {
		spiderManager.list = append(spiderManager.list, abcSpider)
	}
}

type abcNewsStruct struct {
	Collection []struct {
		Link struct {
			To string `json:"to"`
		} `json:"link"`
		Title struct {
			Children string `json:"children"`
		} `json:"title"`
		Timestamp struct {
			Dates struct {
				FirstPublished struct {
					LabelDate string `json:"labelDate"`
				} `json:"firstPublished"`
			} `json:"dates"`
		} `json:"timestamp"`
	} `json:"collection"`
}

func abcSpider() []NewsItem {
	newsItems := make([]NewsItem, 0)
	url := "https://www.abc.net.au/news-web/api/loader/channelrefetch?name=PaginationArticles&documentId=13544780&prepareParams=%7B%22imagePosition%22:%7B%22mobile%22:%22right%22,%22tablet%22:%22right%22,%22desktop%22:%22right%22%7D%7D&loaderParams=%7B%22pagination%22:%7B%22size%22:10%7D%7D&offset=0&size=10&total=250"
	status, resp, err := fasthttp.Get(nil, url)
	if err != nil || status != fasthttp.StatusOK {
		return newsItems
	}

	abcJSON := abcNewsStruct{}
	json.Unmarshal(resp, &abcJSON)

	for _, item := range abcJSON.Collection {
		if IsNeedFilter(item.Title.Children) {
			continue
		}
		newsItems = append(newsItems, NewsItem{
			Title:  item.Title.Children,
			Link:   "https://www.abc.net.au" + item.Link.To,
			Time:   utils.FormatTimeT(item.Timestamp.Dates.FirstPublished.LabelDate),
			Origin: "ABC",
		})
	}

	return newsItems
}
