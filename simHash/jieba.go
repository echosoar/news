package simHash

import (
	"path"
	"runtime"

	"github.com/yanyiwu/gojieba"
)

// 词性对照 https://github.com/fxsjy/jieba#%E5%9F%BA%E4%BA%8E-textrank-%E7%AE%97%E6%B3%95%E7%9A%84%E5%85%B3%E9%94%AE%E8%AF%8D%E6%8A%BD%E5%8F%96
func GetJieba() *gojieba.Jieba {
	_, filePath, _, _ := runtime.Caller(0)
	newDict := path.Join(path.Dir(filePath), "dict.utf8")
	x := gojieba.NewJieba(
		gojieba.DICT_PATH,
		gojieba.HMM_PATH,
		newDict,
		gojieba.IDF_PATH,
		gojieba.STOP_WORDS_PATH,
	)
	return x
}
