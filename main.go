package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/echosoar/news/simHash"
	"github.com/echosoar/news/spider"
	"github.com/echosoar/news/utils"
)

type Result struct {
	Distances []*DitanceItem
	Items     []*ResultItem
}

type ResultItem struct {
	Title    string            `json:"title"`
	Links    []spider.NewsItem `json:"links"`
	Time     int64             `json:"time"`
	Keywords []string          `json:"keywords"`
}

type DitanceItem struct {
	Distance uint64      `json:"distance"`
	Item     *ResultItem `json:"item"`
}

type CacheResult struct {
	Items []*ResultItem `json:"items"`
	Date  string        `json:"time"`
}

func main() {
	isFilterCache := os.Getenv("FILTER_CACHE") == "true"
	list := make([]spider.NewsItem, 0)
	nowDay := utils.FormatNow("2006-01-02")
	nowDayTime := utils.FormatTimeYMDToUnix(nowDay)
	nowTime := time.Now().Unix()
	if nowTime-nowDayTime < 6*3600 {
		nowDayTime = nowDayTime - 60*3600
	}
	result := Result{
		Distances: make([]*DitanceItem, 0),
		Items:     make([]*ResultItem, 0),
	}
	noCache := os.Getenv("NO_CACHE")
	cacheFile := "./result/cache.json"
	if noCache != "true" {
		cacheFileHandler, err := os.Open(cacheFile)
		if err == nil {
			defer cacheFileHandler.Close()
			byteValue, _ := ioutil.ReadAll(cacheFileHandler)
			var cacheStruct CacheResult
			json.Unmarshal([]byte(byteValue), &cacheStruct)
			if cacheStruct.Date == nowDay {
				fmt.Println("load cache", len(cacheStruct.Items))
				for _, items := range cacheStruct.Items {
					list = append(list, items.Links...)
				}
			}
		}
	}

	list = append(list, spider.Get()...)

	x := simHash.GetJieba()
	defer x.Free()

	for _, item := range list {
		if item.Time < nowDayTime {
			continue
		}
		titleLen := float64(len(item.Title))
		if titleLen == 0.0 {
			continue
		}
		if isFilterCache && spider.IsNeedFilter(item.Title) {
			continue
		}
		hash, keywords := simHash.Calc(x, item.Title)
		distance := simHash.Distance(hash)

		isEqual := false
		for _, distanceItem := range result.Distances {
			lenCheck := titleLen / float64(len(distanceItem.Item.Title))
			// title difference too large
			if lenCheck < 0.5 || lenCheck > 2.0 {
				continue
			}
			if simHash.IsEqual(distanceItem.Distance, distance, 6) && utils.CompareKeywords(distanceItem.Item.Keywords, keywords) {
				isEqual = true
				isExists := false
				for _, link := range distanceItem.Item.Links {
					// same source, only one is kept
					if link.Origin == item.Origin {
						isExists = true
						break
					}
				}
				if isExists {
					break
				}
				if len(item.Title) > len(distanceItem.Item.Title) {
					distanceItem.Item.Title = item.Title
					distanceItem.Item.Keywords = keywords
				}
				if item.Time > distanceItem.Item.Time {
					distanceItem.Item.Time = item.Time
				}
				distanceItem.Item.Links = append(distanceItem.Item.Links, item)
				break
			}
		}
		if !isEqual {
			resultItem := ResultItem{
				Title:    item.Title,
				Time:     item.Time,
				Links:    []spider.NewsItem{item},
				Keywords: keywords,
			}
			distanceItem := DitanceItem{
				Distance: distance,
				Item:     &resultItem,
			}
			result.Items = append(result.Items, &resultItem)
			result.Distances = append(result.Distances, &distanceItem)
		}
	}
	sort.Slice(result.Items, func(i, j int) bool {
		jLinksLen := len(result.Items[j].Links)
		iLinksLen := len(result.Items[i].Links)
		if jLinksLen != iLinksLen {
			return jLinksLen < iLinksLen
		}
		return result.Items[j].Time < result.Items[i].Time
	})

	now := utils.FormatNow("2006-01-02 15:04:05")

	if noCache != "true" {
		cacheJson, _ := os.Create(cacheFile)
		defer cacheJson.Close()
		cacheResult := CacheResult{
			Items: result.Items,
			Date:  nowDay,
		}
		cacheJsonStr, _ := json.Marshal(cacheResult)
		cacheJson.Write(cacheJsonStr)
	}

	size := len(result.Items)
	if size > 100 {
		size = 100
	}
	result.Items = result.Items[0:size]

	jsonStr, _ := json.Marshal(result.Items)

	json, _ := os.Create("./result/news.json")
	defer json.Close()
	json.Write(jsonStr)

	jsonp, _ := os.Create("./result/news.jsonp")
	defer jsonp.Close()
	jsonp.Write([]byte("/* */window.newsJsonp && window.newsJsonp(\"" + now + "\", " + string(jsonStr) + ");"))

	mdStr := "## News Update\n---\n" + now + "\n---\n"

	for index, item := range result.Items {
		if len(item.Links) > 1 {
			mdStr += strconv.Itoa(index+1) + ". " + item.Title + " (" + strconv.Itoa(len(item.Links)) + ")\n"
			for _, link := range item.Links {
				mdStr += "    +  " + spider.ItemToHtml(&link) + "\n"
			}
			mdStr += "\n"
		} else {
			mdStr += strconv.Itoa(index+1) + ". " + spider.ItemToHtml(&(item.Links[0])) + "\n"
		}
	}

	// md, _ := os.Create("readme.md")
	// defer md.Close()
	// md.Write([]byte(mdStr))
}
