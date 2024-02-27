package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/echosoar/news/simHash"
	"github.com/echosoar/news/spider"
	"github.com/echosoar/news/utils"
)

type Result struct {
	Items []*ResultItem
}

type ResultItem struct {
	Title  string            `json:"title"`
	Time   int64             `json:"time"`
	Filter []string          `json:"filter"`
	Links  []spider.NewsItem `json:"links"`
}

type JSONItem struct {
	Title string         `json:"title"`
	Time  int64          `json:"time"`
	Links []JSONItemLink `json:"links"`
}
type JSONItemLink struct {
	Origin string `json:"origin"`
	Time   int64  `json:"time"`
}

type CacheResult struct {
	Items []*ResultItem `json:"items"`
	Date  string        `json:"time"`
}

type LinkResult struct {
	Item  spider.NewsItem
	Index int
}

func checkItemIsEqual(item spider.NewsItem, rItem *ResultItem) bool {
	if rItem.Title == item.Title {
		return true
	}
	// checkSize :=
	for start := 0; start < len(rItem.Links); start++ {
		// index := start * 2
		// if index >= len(rItem.Links) {
		// 	break
		// }
		checkItem := rItem.Links[start]
		sameKeywords := utils.CompareKeywords(checkItem.Keywords, item.Keywords)
		if sameKeywords <= 1 {
			return false
		}
		// >= 3
		matchDistance := 12
		switch sameKeywords {
		case 2:
			matchDistance = 6
		}
		if simHash.IsEqual(checkItem.Distance, item.Distance, matchDistance) {
			return true
		}
	}
	return false
}

func compareItemIsEqual(item spider.NewsItem, checkItem spider.NewsItem) bool {
	if checkItem.Title == item.Title {
		return true
	}
	sameKeywords := utils.CompareKeywords(checkItem.Keywords, item.Keywords)
	if sameKeywords <= 1 {
		return false
	}
	// >= 3
	matchDistance := 12
	switch sameKeywords {
	case 2:
		matchDistance = 6
	}
	if simHash.IsEqual(checkItem.Distance, item.Distance, matchDistance) {
		return true
	}
	return false
}

func main() {
	isFilterCache := os.Getenv("FILTER_CACHE") == "true"
	list := make([]spider.NewsItem, 0)
	nowDay, nowDayTime := utils.GetTodayStrAndTime()
	nowTime := time.Now().Unix()
	if nowTime-nowDayTime < 6*3600 {
		nowDayTime = nowDayTime - 60*3600
	}
	result := Result{
		Items: make([]*ResultItem, 0),
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

	x := simHash.GetJieba()
	defer x.Free()

	testTitleList := []string{}

	if len(testTitleList) > 1 {
		hash1, keywords1 := simHash.Calc(x, testTitleList[0])
		distance1 := simHash.Distance(hash1)
		hash2, keywords2 := simHash.Calc(x, testTitleList[1])
		distance2 := simHash.Distance(hash2)
		matched := utils.CompareKeywords(keywords1, keywords2)
		isEqual := simHash.IsEqual(distance1, distance2, int(matched))
		fmt.Println("test1:", testTitleList[0], keywords1, distance1)
		fmt.Println("test2:", testTitleList[1], keywords2, distance2)
		fmt.Println("test-:", matched, isEqual)
		return
	}

	list = append(list, spider.Get()...)

	/*
		全新的比对方法
	*/
	index := 1
	indexLinkMap := make(map[int]int, 0)
	linkResultList := make([]LinkResult, 0)
	for _, item := range list {
		if item.Time < nowDayTime {
			continue
		}
		titleLen := float64(utf8.RuneCountInString(item.Title))
		if titleLen <= 6.0 {
			continue
		}
		// 强制更新 filter
		if isFilterCache {
			item.Title = utils.FormatTitle(item.Title)
			item.Filter = spider.IsNeedFilter(item.Title, []string{})
		}
		// 如果没有关键词，则计算一次
		if len(item.Keywords) == 0 {
			hash, keywords := simHash.Calc(x, item.Title)
			if len(keywords) == 0 {
				continue
			}
			item.Keywords = keywords
			item.Distance = simHash.Distance(hash)
		}

		newLinkItem := LinkResult{
			Item:  item,
			Index: 0,
		}
		needSkip := false
		for _, linkResult := range linkResultList {
			if linkResult.Item.Link == newLinkItem.Item.Link {
				needSkip = true
				break
			}
			if compareItemIsEqual(newLinkItem.Item, linkResult.Item) {
				if newLinkItem.Index == 0 {
					newLinkItem.Index = linkResult.Index
				} else {
					// 找到指向小的，循环合并多个  index，把大的都指向小的
					from := newLinkItem.Index
					to := linkResult.Index
					for {
						if from == to {
							break
						}
						if from < to {
							mid := from
							from = to
							to = mid
						}
						value, exists := indexLinkMap[from]
						indexLinkMap[from] = to
						if !exists {
							break
						}
						from = value
					}
				}
			}
		}

		if needSkip {
			continue
		}
		if newLinkItem.Index == 0 {
			newLinkItem.Index = index
			index += 1
		}
		linkResultList = append(linkResultList, newLinkItem)
	}

	aggregationIndexMap := make(map[int]int, 0)
	speedIndexLinkMap := make(map[int]int, 0)
	for _, linkItem := range linkResultList {
		finalIndex := linkItem.Index
		speedFinalIndex, speedExists := speedIndexLinkMap[linkItem.Index]
		if speedExists {
			finalIndex = speedFinalIndex
		} else {
			for {
				value, exists := indexLinkMap[finalIndex]
				if !exists {
					break
				}
				finalIndex = value
			}
			speedIndexLinkMap[linkItem.Index] = finalIndex
		}
		itemIndex, exists := aggregationIndexMap[finalIndex]
		// 不存在关联，新的新闻
		if !exists {
			aggregationIndexMap[finalIndex] = len(result.Items)
			resultItem := ResultItem{
				Title:  linkItem.Item.Title,
				Time:   linkItem.Item.Time,
				Links:  []spider.NewsItem{linkItem.Item},
				Filter: linkItem.Item.Filter,
			}
			result.Items = append(result.Items, &resultItem)
		} else {
			// 存在关联
			rItem := result.Items[itemIndex]
			item := linkItem.Item
			if item.Time > rItem.Time {
				rItem.Time = item.Time
			}
			// 添加到结果的链接列表中
			rItem.Links = append(rItem.Links, item)
			sort.Slice(rItem.Links, func(i, j int) bool {
				iTitleLen := len(rItem.Links[i].Title)
				jTitleLen := len(rItem.Links[j].Title)
				return jTitleLen < iTitleLen
			})
			center := len(rItem.Links) / 2
			centerItem := rItem.Links[center]
			rItem.Title = centerItem.Title
			rItem.Filter = centerItem.Filter
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
	if size > 150 {
		size = 150
	}

	filtedItems := make([]*ResultItem, 0)
	jsonItems := make([]JSONItem, 0)
	for _, item := range result.Items {
		if len(item.Filter) == 0 {
			filtedItems = append(filtedItems, item)
			jsonItemsLinks := make([]JSONItemLink, len(item.Links))
			for index, link := range item.Links {
				jsonItemsLinks[index] = JSONItemLink{
					Origin: link.Origin,
					Time:   link.Time,
				}
			}
			jsonItems = append(jsonItems, JSONItem{
				Title: item.Title,
				Time:  item.Time,
				Links: jsonItemsLinks,
			})
			if len(filtedItems) >= size {
				break
			}
		}
	}
	if len(result.Items) > size*4 {
		result.Items = result.Items[0 : size*4]
	}

	jsonStr, _ := json.Marshal(jsonItems)

	json, _ := os.Create("./result/news.json")
	defer json.Close()
	json.Write(jsonStr)

	jsonp, _ := os.Create("./result/news.jsonp")
	defer jsonp.Close()
	jsonp.Write([]byte("/* */window.newsJsonp && window.newsJsonp(\"" + now + "\", " + string(jsonStr) + ");"))

	mdStr := "## News Update\n---\n" + now + "\n---\n"

	for index, item := range filtedItems {
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

	mdStr += "\n---\n\n## No Filter News Update\n---\n" + now + "\n---\n"

	for index, item := range result.Items {
		addon := ""
		if len(item.Filter) > 0 {
			addon = "【Filter by '" + strings.Join(item.Filter, "', '") + "'】"
		}
		if len(item.Links) > 1 {
			mdStr += strconv.Itoa(index+1) + ". " + item.Title + addon + " (" + strconv.Itoa(len(item.Links)) + ")\n"
			for _, link := range item.Links {
				mdStr += "    +  " + spider.ItemToHtml(&link) + "\n"
			}
			mdStr += "\n"
		} else {
			mdStr += strconv.Itoa(index+1) + ". " + spider.ItemToHtml(&(item.Links[0])) + addon + "\n"
		}
	}

	md, _ := os.Create("readme.md")
	defer md.Close()
	md.Write([]byte(mdStr))
}
