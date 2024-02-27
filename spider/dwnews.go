package spider

import (
	"encoding/json"
	"os"
	"regexp"
	"strconv"

	"github.com/echosoar/news/utils"
	"github.com/valyala/fasthttp"
)

func init() {
	spiderName := os.Getenv("SPIDER")
	if spiderName == "" || spiderName == "dwnews" {
		spiderManager.list = append(spiderManager.list, dwnewsSpider)
	}
}

type dwRankType []struct {
	ID     int    `json:"id"`
	Schema string `json:"schema"`
}

type dwItemType struct {
	PublishURL  string `json:"publishUrl"`
	Title       string `json:"title"`
	PublishTime int    `json:"publishTime"`
}

func dwnewsSpider() []NewsItem {
	newsItems := make([]NewsItem, 0)
	url := "https://prod-site-api.dwnews.com/v2/articles/"
	status, resp, err := fasthttp.Get(nil, "https://www.dwnews.com/")

	if err != nil || status != fasthttp.StatusOK {
		return newsItems
	}

	reg, _ := regexp.Compile("\"name\":\"新闻排行榜 24h\".*?\"items\":(\\[.*?])")

	rankListStr := reg.FindStringSubmatch(string(resp))

	if len(rankListStr) < 1 {
		return newsItems
	}

	dwRankList := dwRankType{}
	json.Unmarshal([]byte(rankListStr[1]), &dwRankList)

	for _, rank := range dwRankList {
		rankStatus, rankResp, rankErr := fasthttp.Get(nil, url+strconv.Itoa(rank.ID))
		if rankErr != nil || rankStatus != fasthttp.StatusOK {
			continue
		}

		dwItem := dwItemType{}
		json.Unmarshal(rankResp, &dwItem)
		newsItems = append(newsItems, NewsItem{
			Title:  utils.FormatTitle(dwItem.Title),
			Link:   dwItem.PublishURL,
			Filter: IsNeedFilter(dwItem.Title, []string{}),
			Origin: "多维新闻",
			Time:   int64(dwItem.PublishTime),
		})
	}
	return newsItems
}
