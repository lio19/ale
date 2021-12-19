package main

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

type translator struct {
	wordInfoChan chan WordInfo
	errChan chan error

	salt uint64

	lastSendTime int64
	lastSendTimeLock sync.Mutex
}

type WordInfo struct {
	srcMp3Path string
	dstMp3Path string
}

/*
{
  "from": "en",
  "to": "zh",
  "trans_result": [
    {
      "src": "apple",
      "dst": "苹果",
      "src_tts": "https:\/\/fanyiapp.cdn.bcebos.com\/api\/tts\/95e906875b87d342d7325a36a4e1ab42.mp3",
      "dst_tts": "https:\/\/fanyiapp.cdn.bcebos.com\/api\/tts\/62f4ff87617655bc1f65e24cf4ed4963.mp3",
      "dict": "{\"lang\":\"1\",\"word_result\":{\"simple_means\":{\"word_name\":\"apple\",\"from\":\"original\",\"word_means\":[\"苹果\"],\"exchange\":{\"word_pl\":[\"apples\"]},\"tags\":{\"core\":[\"高考\",\"考研\"],\"other\":[\"\"]},\"symbols\":[{\"ph_en\":\"ˈæpl\",\"ph_am\":\"ˈæpl\",\"parts\":[{\"part\":\"n.\",\"means\":[\"苹果\"]}],\"ph_other\":\"\"}]}}}"
    }
  ]
}
*/

type result struct {
	ErrorCode   string   `json:"errorCode"`
	Query       string   `json:"query"`
	Translation []string `json:"translation"`
	TSpeakUrl string `json:"tSpeakUrl"`
	SpeakUrl  string `json:"speakUrl"`
}

func NewTranslator() *translator {
	return &translator{
		wordInfoChan: make(chan WordInfo, 10),
		errChan: make(chan error, 2),
	}
}

func (t *translator) Translate(word string) {
	go func() {

		//控制在qps = 1
		for {
			t.lastSendTimeLock.Lock()
			if time.Now().UnixNano() - t.lastSendTime < 1000000 {
				t.lastSendTimeLock.Unlock()
				time.Sleep(1 * time.Second)
			} else {
				t.lastSendTime = time.Now().UnixNano()
				t.lastSendTimeLock.Unlock()
				break
			}
		}

		//请求
		t.salt++
		curTime := strconv.FormatInt(time.Now().Unix(), 10)
		saltStr :=strconv.FormatUint(t.salt, 10)
		sig := createSignal(word,saltStr , curTime)

		url := fmt.Sprintf("https://openapi.youdao.com/api?q=%s&from=zh-CHS&to=EN&appKey=%s&salt=%s&sign=%s&signType=v3&curtime=%s&ext=mp3&voice=0&strict=true", word, conf.getValue("appid"), saltStr, sig, curTime)

		resp, err := http.Get(url)
		defer resp.Body.Close()

		if err != nil {
			t.errChan <- err
			return
		}

		bs, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.errChan <- err
			return
		}

		var r result
		err = json.Unmarshal(bs, &r)
		if err != nil {
			t.errChan <- err
			return
		}

		if r.ErrorCode != "0" {
			t.errChan <- errors.New("get result less 0")
			return
		}

		var wi WordInfo
		wi.srcMp3Path, err = downloadMp3(r.SpeakUrl, conf.getValue("mp3TempDir"), word + ".en.mp3")
		if err != nil {
			t.errChan <- err
			return
		}

		wi.dstMp3Path, err = downloadMp3(r.TSpeakUrl, conf.getValue("mp3TempDir"), word + ".ch.mp3")
		if err != nil {
			t.errChan <- err
			return
		}

		t.wordInfoChan <- wi
	}()
}


func createSignal(word, salt, curTime string) string{
	var input string
	if len(word) <= 20 {
		input = word
	} else {
		input = fmt.Sprintf("%s%d%s", word[:10], len(word), word[10:])
	}


	return fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprint(conf.getValue("appid"), input, salt, curTime,conf.getValue("miyao")))))
}

func downloadMp3(uri, dirPath, filename string) (string,error) {

	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err = os.MkdirAll(dirPath, 0777); err != nil {
			return "", err
		}
	}

	//判断是非存在相同的文件
	fileTempPath := fmt.Sprintf("%s%c%s",dirPath,filepath.Separator,filename)
	_, err := os.Stat(fileTempPath)
	if os.IsExist(err) {
		return fileTempPath, nil
	}

	resp, err := http.Get(uri)
	defer resp.Body.Close()

	if err != nil {
		return "", err
	}

	fileTempFd, err := os.Create(fileTempPath)
	defer fileTempFd.Close()
	if err != nil {
		return "",err
	}

	c, err := io.Copy(fileTempFd, resp.Body)
	if err != nil {
		return "", err
	} else {
		logger.Info("copy count = ", c)
	}

	return fileTempPath, nil
}
