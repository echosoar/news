package utils

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

var stripPrefixList []string = []string{
	"|",
	"丨",
	"｜",
	"】",
}

var emojiRegex = regexp.MustCompile(`[\x{1F300}-\x{1F6FF}|[\x{2600}-\x{26FF}]|[\x{2700}-\x{27BF}]|\x{24C2}|[\x{1F900}-\x{1F9FF}]|\x{1F1E6}-\x{1F1FF}|\x{1F191}-\x{1F251}|\x{1F600}-\x{1F64F}]`)
var stripCharRegex = regexp.MustCompile(`(^[\s\.，，]+)|([\s\.，，…！!]+$)|[”“「」]`)
var spaceRegex = regexp.MustCompile(`[!,！，]+`)

func FormatTitle(title string) string {
	// 去除前缀
	for _, preifx := range stripPrefixList {
		if strings.Contains(title, preifx) {
			list := strings.Split(title, preifx)
			firstLen := utf8.RuneCountInString(list[0])
			secondLen := utf8.RuneCountInString(list[1])
			if firstLen < secondLen {
				title = list[1]
			}
		}
		if len(title) <= 1 {
			break
		}
	}
	// 去除 emoji
	title = emojiRegex.ReplaceAllString(title, "")
	// 去除前后无效字符
	title = stripCharRegex.ReplaceAllString(title, "")
	// 替换字符
	title = spaceRegex.ReplaceAllString(title, "，")
	return title
}
