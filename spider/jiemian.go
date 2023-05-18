package spider

import (
	"encoding/json"
	"os"
	"regexp"

	"github.com/echosoar/news/utils"
	"github.com/valyala/fasthttp"
)

func init() {
	spiderName := os.Getenv("SPIDER")
	if spiderName == "" || spiderName == "jiemian" {
		spiderManager.list = append(spiderManager.list, jiemianSpider)
	}
}

type jiemianJson struct {
	Rst string `json:"rst"`
}

func jiemianSpider() []NewsItem {
	url := "https://a.jiemian.com/index.php?m=lists&a=ajaxNews&cid=4&page=1"
	newsItems := make([]NewsItem, 0)
	status, resp, err := fasthttp.Get(nil, url)
	if err != nil || status != fasthttp.StatusOK {
		return newsItems
	}
	jiemianJsonStruct := jiemianJson{}
	json.Unmarshal(resp[1:len(resp)-1], &jiemianJsonStruct)

	reg, _ := regexp.Compile("\"item-date\">(.*?)</div><div class=\"item-main\"><p>\n\\s*.*?href=\"(.*?)\".*?\"_blank\">(.*?)</a>.*?\n\\s*(.*?)</p>")
	res := reg.FindAllStringSubmatch(jiemianJsonStruct.Rst, -1)
	newsItems = make([]NewsItem, 0, len(res))
	for _, matchedItem := range res {
		if IsNeedFilter(matchedItem[3], []string{"界面"}) {
			continue
		}
		newsItems = append(newsItems, NewsItem{
			Title:  utils.FormatTitle(matchedItem[3]),
			Link:   matchedItem[2],
			Origin: "界面新闻",
			Time:   utils.FormatTimeHMToUnix(matchedItem[1]),
		})
	}
	return newsItems
}
