package utils

import "strconv"

func CompareKeywords(keywords1, keywords2 []string) int32 {
	var matched int32 = 0
	for _, keyword := range keywords1 {
		// 关键词不是数字
		// 并且在两个关键词列表中都存在
		if !isNumeric(keyword) && Contains(keywords2, keyword) {
			matched++
		}
	}
	return matched
}
func isNumeric(str string) bool {
	_, err := strconv.Atoi(str)
	return err == nil
}

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
