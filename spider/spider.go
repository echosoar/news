package spider

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/echosoar/news/utils"
)

var filterList [][]string = [][]string{
	{"学习"},
	{"思想"},
	{"牢牢"},
	{"扎根"},
	{"基层"},
	{"奋斗"},
	{"明天"},
	{"绿水青山"},
	{"伟大"},
	{"为人民"},
	{"和人民"},
	{"讲述"},
	{"微记录"},
	{"前程"},
	{"成长"},
	{"足迹"},
	{"见证"},
	{"谱写"},
	{"致辞"},
	{"红火"},
	{"给你我"},
	{"奔跑"},
	{"聆听"},
	{"追梦"},
	{"更美好"},
	{"关心"},
	{"主题"},
	{"面对面"},
	{"摇篮"},
	{"期待"},
	{"温度"},
	{"网评"},
	{"会议"},
	{"惦念"},
	{"记在"},
	{"委员"},
	{"正能量"},
	{"幸福"},
	{"乡愁"},
	{"讲话"},
	{"开门红"},
	{"走进"},
	{"开启"},
	{"奋进"},
	{"奋斗"},
	{"治国"},
	{"理政"},
	{"青春"},
	{"当选", "省"},
	{"当选", "市"},
	{"来了"},
	{"引领"},
	{"新时代"},
	{"闭幕"},
	{"访谈"},
	{"观察"},
	{"故事"},
	{"精神"},
	{"不负"},
	{"嘱托"},
	{"回来"},
	{"跑出"},
	{"加速度"},
	{"记录"},
	{"亮出"},
	{"走进"},
	{"看中国"},
	{"复苏"},
	{"领会"},
	{"信心"},
	{"出版"},
	{"征程"},
	{"出发"},
	{"团结"},
	{"会见"},
	{"欢迎"},
	{"魅力"},
	{"夯实"},
	{"韧性"},
	{"发展"},
	{"飘扬"},
	{"总书记"},
	{"惹人"},
	{"岁月"},
	{"正当时"},
	{"愿景"},
	{"机遇"},
	{"书写"},
	{"改写"},
	{"舒适"},
	{"沃土"},
	{"他们"},
	{"见闻"},
	{"就业"},
	{"样子"},
	{"敬献"},
	{"热闹"},
	{"收藏"},
	{"会谈"},
	{"春天"},
	{"春色"},
	{"春风"},
	{"胜利"},
	{"活力"},
	{"民心"},
	{"赚钱"},
	{"致富"},
	{"脱贫"},
	{"难忘"},
	{"非凡"},
	{"山水"},
	{"开学"},
	{"建议"},
	{"民生"},
	{"消费"},
	{"为啥"},
	{"成效"},
	{"扎实"},
	{"美好"},
	{"两会"},
	{"中国梦"},
	{"合力"},
	{"高质量"},
	{"成熟"},
	{"推进"},
	{"开局"},
	{"提效"},
	{"风光"},
	{"涌现"},
	{"明确"},
	{"为党"},
	{"嘛"},
	{"攻关"},
	{"制度"},
	{"逆袭"},
	{"坚守"},
	{"认识"},
	{"回升"},
	{"急了"},
	{"震撼"},
	{"攻坚"},
	{"综述"},
	{"美丽"},
	{"完成"},
	{"百姓"},
	{"自信"},
	{"民主"},
	{"加快"},
	{"中国特色"},
	{"好帮手"},
	{"尊重"},
	{"近平"},
	{"现代化"},
	{"外交"},
	{"大片"},
	{"航拍"},
	{"自己"},
	{"中国式"},
	{"群众"},
	{"热情"},
	{"成绩"},
	{"速览"},
	{"加强"},
	{"效率"},
	{"中", "关系"},
	{"造福"},
	{"第", "集"},
	{"协商"},
	{"通道"},
	{"亮点"},
	{"旅行"},
	{"文旅"},
	{"的美"},
	{"格局"},
	{"塑造"},
	{"平凡"},
	{"传统"},
	{"专访"},
	{"美食"},
	{"知识点"},
	{"焦虑"},
	{"读懂"},
	{"振兴"},
	{"多美"},
	{"迈进"},
	{"视频"},
	{"中年"},
	{"评论"},
	{"梦想"},
	{"诠释"},
	{"赋能"},
	{"青睐"},
	{"精心"},
	{"这个"},
	{"在家里"},
	{"在民间"},
	{"科普"},
}

type NewsItem struct {
	Title  string
	Link   string
	Origin string
	Time   int64
}

type SpiderManager struct {
	list []Spider
}

var spiderManager = &SpiderManager{
	list: make([]Spider, 0),
}

type Spider = func() []NewsItem

func Get() []NewsItem {
	spiderCount := len(spiderManager.list)
	var wg sync.WaitGroup
	wg.Add(spiderCount)
	channel := make(chan *[]NewsItem, spiderCount)
	for i := 0; i < spiderCount; i++ {
		go func(ch chan *[]NewsItem, index int, wg *sync.WaitGroup) {
			list := spiderManager.list[index]()
			ch <- &list
			wg.Done()
		}(channel, i, &wg)
	}
	go func() {
		wg.Wait()
		close(channel)
	}()
	newsItems := make([]NewsItem, 0)
	for ch := range channel {
		newsItems = append(newsItems, *ch...)
	}
	spiderName := os.Getenv("SPIDER")
	if spiderName != "" {
		fmt.Println("newsItems", newsItems)
	}
	return newsItems
}

func IsNeedFilter(title string) bool {
	for _, filterWords := range filterList {
		isNeedFilter := true
		for _, word := range filterWords {
			if !strings.Contains(title, word) {
				isNeedFilter = false
				break
			}
		}
		if isNeedFilter {
			// if strings.Contains(title, "---") {
			// 	fmt.Println("filter", title, filterWords)
			// }
			return true
		}
	}
	return false
}

func ItemToHtml(item *NewsItem) string {
	pubTime := ""
	if item.Time > 0 {
		pubTime = " - " + utils.FormatTime(item.Time, "01/02 15:04")
	}
	return "<a target=\"_blank\" href=\"" + item.Link + "\">" + item.Title + "</a> [" + item.Origin + pubTime + "]"
}
