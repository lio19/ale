package player

import (
	"autoLearnEnglish/src/conf"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_playWordByLetter(t *testing.T) {

	a := assert.New(t)

	conf.MockConf(map[string]string{
		"lettersDir": "../../letters",
	})

	p := New()
	err := p.playWordByLetter("abandon")
	a.Nil(err)

	err = p.playPath("../../letters/a.ch.mp3")
	a.Nil(err)
}
