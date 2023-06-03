package utils

import (
	"strconv"
)

func CompareKeywords(keywords1, keywords2 []string) int32 {
	var matched float64 = 0.0
	for _, keyword := range keywords1 {
		// 并且在两个关键词列表中都存在
		if Contains(keywords2, keyword) {
			isIntNum, num := isNumeric(keyword)
			if isIntNum {
				switch num {
				case 315:
					fallthrough
				case 520:
					matched += 1.0
				default:
					matched += 0.6
				}
			} else {
				matched++
			}
		}
	}
	return int32(matched)
}
func isNumeric(str string) (bool, int) {
	num, err := strconv.Atoi(str)
	return err == nil, num
}

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
