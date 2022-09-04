package main

import (
	"bufio"
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"io"
	"os"
	"sync"
	"time"
)

// ALE Auto learn English
type ALE struct {
	wordInfos []WordInfo
	wordInfoLock sync.Mutex
	translator   *translator

	speakerSampleRate beep.SampleRate
}

func NewALE() *ALE{
	ale := &ALE{
		wordInfos:    make([]WordInfo, 0),
		wordInfoLock: sync.Mutex{},
		translator:   NewTranslator(),
		speakerSampleRate: beep.SampleRate(16000),
	}

	err := speaker.Init(ale.speakerSampleRate, ale.speakerSampleRate.N(time.Second/10))
	if err != nil {
		panic(err.Error())
	}

	return ale
}

func (a *ALE)Start() {

	//1.

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
	if err := play(info.srcMp3Path, a.speakerSampleRate); err != nil {
		return err
	}
	time.Sleep(1 * time.Second)
	if err := play(info.dstMp3Path, a.speakerSampleRate); err != nil {
		return err
	}

	return nil
}

func playWordByLetter(word string, sr beep.SampleRate) error {

	seqSlice := make([]beep.Streamer, 0)
	fdSlice := make([]io.ReadCloser, 0)
	for _, v := range word {
		f, err := os.Open(fmt.Sprintf("%c.ch.mp3", v))
		fdSlice = append(fdSlice, f)

		streamer, format, err := mp3.Decode(f)
		reStream := beep.Resample(4, format.SampleRate, sr, streamer)

		seqSlice = append(seqSlice, reStream)
	}

	speaker.Play(beep.Seq(seqSlice..., beep.Callback(func() {

	})))
}

func play(str string, sr beep.SampleRate) error {
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


	done := make(chan bool)

	reStream := beep.Resample(4, format.SampleRate, sr, streamer)

	speaker.Play(beep.Seq(reStream, beep.Callback(func() {
		done <- true
	})))


	<-done
	logger.Info("play one mp3", str)

	return nil
}


