package main

import (
	"autoLearnEnglish/src/cache"
	"autoLearnEnglish/src/player"
	"autoLearnEnglish/src/youdao"
	"sync"
	"time"
)

// ALE Auto learn English
type ALE struct {
	wordInfoLock sync.Mutex
	wPlay        *player.Player
	wCache       *cache.Cache
	wYouDao      *youdao.WordBook
}

func NewALE() *ALE {
	ale := &ALE{
		wordInfoLock: sync.Mutex{},
		wPlay:        player.New(),
		wYouDao:      youdao.NewWordBook(),
		wCache:       cache.New(),
	}

	return ale
}

func (a *ALE) Start() {

	//1.获取cache
	wl, err := a.wCache.GetWordList()
	if err != nil {
		panic(err)
	}

	var subTime int64 = 0
	if wl != nil {
		subTime = time.Now().Unix() - wl.AddTime
	}
	//1. cache is nil
	//2. overtime
	if wl == nil || subTime < 0 || subTime > 60*60*60 {
		rwl, err := a.wYouDao.Get()
		if err != nil {
			logger.Warn("get you dao word book fail, err = ", err.Error())
			panic(err)
		}

		//update cache
		wordStrList := youDaoWL2WordList(rwl)
		if err = a.wCache.UpdateWordList(wordStrList); err != nil {
			logger.Error("update word list to cache fail err = ", err.Error())
		}

		if wl, err = a.wCache.GetWordList(); err != nil {
			logger.Error("get cache word list fail err = ", err.Error())
			panic(err)
		}
	}

	wordDetlCacheList := make([]*cache.WordDetailCache, 0)
	for _, wordStr := range wl.Words {

		wordCache, err := a.wCache.GetOneWordDetail(wordStr)
		if err != nil {
			logger.Error("get one word cache err = ", err.Error())
			continue
		}

		if wordCache == nil {
			//cache is nil
			if rwi, err := a.wYouDao.GetOneWord(wordStr); err != nil {
				//get info fail, skip
				continue
			} else {
				//update cache
				if err = a.wCache.UpdateWordInfo(&cache.InputWordCache{
					EngMp3Url: rwi.SpeakUrl,
					ChMp3Url:  rwi.TSpeakUrl,
					WordEng:   wordStr,
					WordCh:    "",
				}); err != nil {
					logger.Error("update word cache fail, word = ", wordStr, " err info = ", err.Error())
					continue
				}

				if wordCache, err = a.wCache.GetOneWordDetail(wordStr); err != nil {
					logger.Error("update word cache fail, word = ", wordStr)
					continue
				}
			}
		}

		//wordCache must have value
		wordDetlCacheList = append(wordDetlCacheList, wordCache)
	}

	//wordDetailCache to wordPlayInfo
	var wordPlayInfoList []*player.WordPlayInfo
	for _, v := range wordDetlCacheList {
		wordPlayInfoList = append(wordPlayInfoList, &player.WordPlayInfo{
			IdleTime:   time.Second * 1,
			EngMp3Path: v.EngMp3,
			ChMp3Path:  v.ChMp3,
			Word:       v.English,
		})
	}

	//begin play
	a.autoLearn(wordPlayInfoList)
}

func (a *ALE) autoLearn(wordPlayInfoList []*player.WordPlayInfo) {
	logger.Info("auto learn begin")
	//一直循环播放word info中的单词
	for {
		time.Sleep(1 * time.Second)
		for _, v := range wordPlayInfoList {
			if err := a.wPlay.PlayWord(v); err != nil {
				logger.Error(err.Error())
			}
		}
	}
}

func youDaoWL2WordList(rwb *youdao.RespWordBookData) []string {
	var wl []string
	for _, v := range rwb.ItemList {
		wl = append(wl, v.Word)
	}
	return wl
}

func youDaoWl2PlayWl(rwb *youdao.RespWordBookData) *player.WordPlayInfo {
	return nil
}

func cacheWl2PlayWl(list *cache.WordListCache) *player.WordPlayInfo {
	return nil
}
