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
	if spiderName == "" || spiderName == "ithome" {
		spiderManager.list = append(spiderManager.list, ithomeSpider)
	}
}

type ithomeJson struct {
	DataList []struct {
		Url   string `json:"WapNewsUrl"`
		Title string `json:"title"`
		Time  string `json:"postdate"`
	} `json:"Result"`
}

func ithomeSpider() []NewsItem {
	url := "https://m.ithome.com/api/news/newslistpageget?Tag=&ot=" + strconv.Itoa(int(time.Now().Unix())) + "000&page=0&hitCountAuthority=false"
	newsItems := make([]NewsItem, 0)
	status, resp, err := fasthttp.Get(nil, url)
	if err != nil || status != fasthttp.StatusOK {
		return newsItems
	}

	jsonStruct := ithomeJson{}
	json.Unmarshal(resp, &jsonStruct)
	for _, item := range jsonStruct.DataList {
		if IsNeedFilter(item.Title, []string{}) {
			continue
		}
		newsItems = append(newsItems, NewsItem{
			Title:  utils.FormatTitle(item.Title),
			Link:   item.Url,
			Origin: "IT之家",
			Time:   utils.FormatTimeTLocation(item.Time),
		})
	}
	return newsItems
}
