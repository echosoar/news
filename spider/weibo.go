package spider

import (
	"encoding/json"
	"os"

	"github.com/valyala/fasthttp"
)

func init() {
	spiderName := os.Getenv("SPIDER")
	if spiderName == "" || spiderName == "weibo" {
		spiderManager.list = append(spiderManager.list, weiboSpider)
	}
}

type weiboJson struct {
	Ok   int `json:"ok"`
	Data struct {
		Cards []struct {
			CardGroup []struct {
				Scheme string `json:"scheme"`
				Desc   string `json:"desc"`
			} `json:"card_group"`
		} `json:"cards"`
	} `json:"data"`
}

func weiboSpider() []NewsItem {
	url := "https://m.weibo.cn/api/container/getIndex?containerid=106003type%3D25%26t%3D3%26disable_hot%3D1%26filter_type%3Drealtimehot&title=%E5%BE%AE%E5%8D%9A%E7%83%AD%E6%90%9C&extparam=filter_type%3Drealtimehot%26mi_cid%3D100103%26pos%3D0_0%26c_type%3D30%26display_time%3D1540538388&luicode=10000011&lfid=231583"
	newsItems := make([]NewsItem, 0)
	status, resp, err := fasthttp.Get(nil, url)
	if err != nil || status != fasthttp.StatusOK {
		return newsItems
	}

	weiboJsonStruct := weiboJson{}
	json.Unmarshal(resp, &weiboJsonStruct)

	for _, item := range weiboJsonStruct.Data.Cards[0].CardGroup {
		if len(item.Desc) < 5 {
			continue
		}
		if IsNeedFilter(item.Desc) {
			continue
		}
		newsItems = append(newsItems, NewsItem{
			Title:  item.Desc,
			Link:   item.Scheme,
			Origin: "微博",
		})
	}
	return newsItems
}
