package main

import (
	"bufio"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"os"
	"sync"
	"time"
)

// ALE Auto learn English
type ALE struct {
	wordInfos []WordInfo
	wordInfoLock sync.Mutex
	translator   *translator
}

func NewALE() *ALE{
	return &ALE{
		wordInfos:    make([]WordInfo, 0),
		wordInfoLock: sync.Mutex{},
		translator:   NewTranslator(),
	}
}

func (a *ALE)Start() {

	//读入所有的单词
	wordsTxtPath := conf.getValue("wordstxt")

	wordsTxtFd, err := os.Open(wordsTxtPath)
	if err != nil {
		panic(err.Error())
	}

	sc := bufio.NewScanner(wordsTxtFd)

	for sc.Scan() {
		oneWord := sc.Text()
		if oneWord != "" {
			a.translator.Translate(oneWord)
		}
	}

	//开启
	a.autoLearn()
}

func (a *ALE) autoLearn() {
	logger.Info("auto learn begin")
	//一直循环播放word info中的单词
	index := 0
	for {

		time.Sleep(1 * time.Second)
		select {
		case wi, ok := <- a.translator.wordInfoChan:
			if !ok {
				logger.Error("get word info chan info fail")
			} else {
				a.wordInfos = append(a.wordInfos, wi)
			}
			break
		case err := <- a.translator.errChan:
			if err != nil {
				panic(err.Error())
			}
			break
		default:
			//没有数据
			if len(a.wordInfos) == 0 {
				break
			}

			if index >= len(a.wordInfos) {
				index = 0
			}

			if err := a.playWordInfo(a.wordInfos[index]); err != nil {
				logger.Error(err.Error())
			}
			index++
		}
	}
}

func (a *ALE) playWordInfo(info WordInfo) error {
	time.Sleep(1 * time.Second)
	if err := play(info.srcMp3Path); err != nil {
		return err
	}
	time.Sleep(1 * time.Second)
	if err := play(info.dstMp3Path); err != nil {
		return err
	}

	return nil
}

func play(str string) error {
	f, err := os.Open(str)
	if err != nil {
		return err
	}
	defer f.Close()

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		return err
	}
	defer streamer.Close()

	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		return err
	}
	defer speaker.Close()

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	<-done
	logger.Info("play one mp3", str)

	return nil
}


