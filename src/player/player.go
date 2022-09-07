package player

import (
	"autoLearnEnglish/src/conf"
	"errors"
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Player struct {
	lettersDir string
}

func New() *Player {

	err := speaker.Init(beep.SampleRate(16000), beep.SampleRate(16000).N(time.Second/10))
	if err != nil {
		panic(err.Error())
	}
	return &Player{lettersDir: conf.GetValue("lettersDir")}
}

func (p *Player) PlayWord(info *WordPlayInfo) error {

	if info == nil {
		return errors.New("play word info is nil")
	}

	//1.enMp3
	//2.sleep
	//3.spell
	//4.sleep
	//5.enMp3
	//6.sleep
	//7.chMp3
	if err := p.playPath(info.EngMp3Path); err != nil {
		return err
	}
	time.Sleep(500 * time.Millisecond)
	if err := p.playPath(info.ChMp3Path); err != nil {
		return err
	}
	time.Sleep(info.IdleTime)
	if err := p.playPath(info.EngMp3Path); err != nil {
		return err
	}
	time.Sleep(300 * time.Millisecond)
	if err := p.playWordByLetter(info.Word); err != nil {
		return err
	}
	time.Sleep(info.IdleTime)
	if err := p.playPath(info.EngMp3Path); err != nil {
		return err
	}
	time.Sleep(500 * time.Millisecond)
	if err := p.playPath(info.ChMp3Path); err != nil {
		return err
	}

	return nil
}

func (p *Player) playPath(str string) error {
	f, err := os.Open(str)
	if err != nil {
		return err
	}
	defer f.Close()

	streamer, _, err := mp3.Decode(f)
	if err != nil {
		return err
	}
	defer streamer.Close()

	done := make(chan bool)

	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))

	<-done
	return nil
}

func (p *Player) playWordByLetter(word string) error {

	seqSlice := make([]beep.Streamer, 0)
	closeStream := make([]beep.StreamSeekCloser, 0)
	fdSlice := make([]io.ReadCloser, 0)

	for _, v := range strings.ToLower(word) {
		f, err := os.Open(fmt.Sprintf("%s%c%c.ch.mp3", p.lettersDir, filepath.Separator, v))
		if err != nil {
			return err
		} else {
			fdSlice = append(fdSlice, f)
		}

		streamer, _, err := mp3.Decode(f)
		seqSlice = append(seqSlice, streamer)
		closeStream = append(closeStream, streamer)
	}

	defer func() {
		for _, v := range fdSlice {
			_ = v.Close()
		}

		for _, v := range closeStream {
			_ = v.Close()
		}
	}()

	done := make(chan bool)
	seqSlice = append(seqSlice, beep.Callback(func() {
		done <- true
	}))

	speaker.Play(beep.Seq(seqSlice...))

	<-done
	return nil
}
