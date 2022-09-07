package cache

import (
	"autoLearnEnglish/src/conf"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

//提供两个功能
//1. word list
//2. one word cache (mp3 + english to chinese)
//|
//|
//|----WordListCache
//|
//|---Words----word1-----ch.mp3
//          |         |
//          |         ---en.mp3
//          |         |
//          |         ---info
//          |
//          ---word2
//          |
//          ---word3
//
//
//WordList是一个文本，以json内容，存储着所有的单词和添加时间
//{
//	"AddTime": 1,
//	"Words": [
//		"1",
//		"2",
//		"3"
//	]
//}
//
//Word是一个文件夹，word1是每一个具体的单词文件夹
//
//word1中:
//ch.mp3中文读音
//en.mp3英文读音
//info是json文本，里面保存着中文释义和添加释义，和最后一次访问时间
//{
//	"WordEng": "",
//	"WordCh": "",
//	"AddTime": 1
//}

type Cache struct {
	local string
}

func New() *Cache {
	//todo: 是否需要检查是否存在?
	local := conf.GetValue("cache")
	return &Cache{
		local: local,
	}
}

func (c *Cache) wordListCacheFile() string {
	return fmt.Sprintf("%s%cWordList", c.local, filepath.Separator)
}

func (c *Cache) wordCacheFile(wordStr string) *innerWordFilePath {
	var wfp innerWordFilePath
	wfp.DirPath = fmt.Sprintf("%s%cWords%c%s", c.local, filepath.Separator, filepath.Separator, wordStr)
	wfp.EngMp3Path = fmt.Sprintf("%s%cen.mp3", wfp.DirPath, filepath.Separator)
	wfp.ChMp3Path = fmt.Sprintf("%s%cch.mp3", wfp.DirPath, filepath.Separator)
	wfp.InfoPath = fmt.Sprintf("%s%cinfo", wfp.DirPath, filepath.Separator)
	return &wfp
}

// GetWordList cache do not exit return nil, nil
func (c *Cache) GetWordList() (*WordListCache, error) {

	bs, err := os.ReadFile(c.wordListCacheFile())

	if os.IsNotExist(err) {
		return nil, nil
	}

	if err != nil {
		return nil, errors.New("read word list cache fail, msg = " + err.Error())
	}

	var wl WordListCache
	err = json.Unmarshal(bs, &wl)
	if err != nil {
		return nil, errors.New("cache content error, msg = " + err.Error())
	}

	return &wl, nil
}

// UpdateWordList wl is nil, delete WordListCache
func (c *Cache) UpdateWordList(words []string) error {

	wlp := c.wordListCacheFile()

	if words == nil {
		if _, err := os.Stat(wlp); err != nil && os.IsNotExist(err) {
			return nil
		}
		//删掉文件
		if err := os.Remove(wlp); err != nil {
			return err
		}
		return nil
	}

	if !fileIsExit(c.local) {
		if err := os.MkdirAll(c.local, os.ModePerm); err != nil {
			panic("create cache root dir fail, path =  " + c.local + "   err = " + err.Error())
		}
	}

	wl := WordListCache{
		AddTime: time.Now().Unix(),
		Words:   words,
	}
	bs, err := json.Marshal(wl)
	if err != nil {
		return errors.New("word list to cache fail, msg = " + err.Error())
	}

	fd, err := os.Create(wlp)
	if err != nil {
		return errors.New("word list to cache fail, create cache file fail, msg = " + err.Error())
	}

	defer fd.Close()

	if _, err = fd.Write(bs); err != nil {
		return errors.New("write word list to cache file fail, msg = " + err.Error())
	}

	return nil
}

func (c *Cache) UpdateWordInfo(wc *InputWordCache) error {
	if wc == nil {
		return nil
	}

	if wc.WordEng == "" || wc.EngMp3Url == "" || wc.ChMp3Url == "" {
		return fmt.Errorf("word info is error, %+v", *wc)
	}

	wcf := c.wordCacheFile(wc.WordEng)

	//is dir exit? create it
	if !fileIsExit(wcf.DirPath) {
		if err := os.MkdirAll(wcf.DirPath, os.ModePerm); err != nil {
			return errors.New("mkdir word cache dir fail, msg = " + err.Error())
		}
	}

	//download
	if err := cacheHttpUrlToPath(wc.EngMp3Url, wcf.EngMp3Path); err != nil {
		return errors.New("cache http to path fail, msg = " + err.Error())
	}

	if err := cacheHttpUrlToPath(wc.ChMp3Url, wcf.ChMp3Path); err != nil {
		return errors.New("cache http to path fail, msg = " + err.Error())
	}

	wic := innerWordCache{
		WordEng: wc.WordEng,
		WordCh:  wc.WordCh,
		AddTime: time.Now().Unix(),
	}

	bs, err := json.Marshal(wic)
	if err != nil {
		return errors.New("marshal json fail, msg = " + err.Error())
	}

	infoFd, err := os.Create(wcf.InfoPath)
	if err != nil {
		return errors.New("create file info fail msg = " + err.Error())
	}

	_, err = infoFd.Write(bs)
	if err != nil {
		return errors.New("write file info fail msg = " + err.Error())
	}

	return nil
}

func (c *Cache) GetOneWordDetail(word string) (*WordDetailCache, error) {

	wp := c.wordCacheFile(word)
	//check is all file exit

	if !fileIsExit(wp.DirPath) {
		return nil, nil
	}

	if !fileIsExit(wp.EngMp3Path) || !fileIsExit(wp.ChMp3Path) {
		_ = os.Remove(wp.DirPath)
		return nil, errors.New("cache is dirty english mp3 or chinese mp3 dono")
	}

	bs, err := os.ReadFile(wp.InfoPath)
	if err != nil {
		_ = os.Remove(wp.DirPath)
		return nil, errors.New("get cache info fail, msg = " + err.Error())
	}

	var wi innerWordCache
	if err = json.Unmarshal(bs, &wi); err != nil {
		_ = os.Remove(wp.DirPath)
		return nil, errors.New("get cache info fail, msg = " + err.Error())
	}

	var w WordDetailCache
	w.AddTime = wi.AddTime
	w.English = wi.WordEng
	w.Chinese = wi.WordCh
	w.ChMp3 = wp.ChMp3Path
	w.EngMp3 = wp.EngMp3Path

	return &w, nil
}

func cacheHttpUrlToPath(httpUrl, localPath string) error {
	resp, err := http.Get(httpUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	//don't success
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("http request fail, status code = %d", resp.StatusCode)
	}

	fd, err := os.Create(localPath)
	if err != nil {
		return errors.New("create chinese voice mp3 file fail, msg = %s" + err.Error())
	}
	_, err = io.Copy(fd, resp.Body)
	if err != nil {
		return errors.New("cache chinese voice mp3 file fail, msg = %s" + err.Error())
	}

	return nil
}

func fileIsExit(p string) bool {
	//todo: 有没有可能别的error导致获取stat失败？ 现在整个程序都没有处理这种问题
	_, err := os.Stat(p)
	return err == nil
}
