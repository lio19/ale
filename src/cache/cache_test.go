package cache

import (
	"autoLearnEnglish/src/conf"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestCache_GetWordList(t *testing.T) {
	a := assert.New(t)

	testDir := "test_cache_dir"
	err := os.Mkdir(testDir, os.ModePerm)
	a.Nil(err, "create test dir fail")

	defer os.RemoveAll(testDir)

	conf.MockConf(map[string]string{"cache": "test_cache_dir"})

	c := New()
	a.NotNil(c)

	c.wordCacheFile("one")
	wl, err := c.GetWordList()
	a.Nil(err)
	a.Nil(wl)

	err = c.UpdateWordList([]string{"test"})
	a.Nil(err)
	wl, err = c.GetWordList()
	a.Nil(err)
	a.Equal(len(wl.Words), 1)
	a.Equal(wl.Words[0], "test")
	err = c.UpdateWordList(nil)
	wl, err = c.GetWordList()
	a.Nil(err)
	a.Nil(wl)
}

func TestCache_GetOneWordInfo(t *testing.T) {
	a := assert.New(t)

	testDir := "test_cache_dir"
	err := os.Mkdir(testDir, os.ModePerm)
	a.Nil(err, "create test dir fail")

	conf.MockConf(map[string]string{"cache": "test_cache_dir"})

	defer os.RemoveAll(testDir)

	c := New()
	a.NotNil(c)
	wi, err := c.GetOneWordDetail("test")
	a.Nil(err)
	a.Nil(wi)

	wc := InputWordCache{
		EngMp3Url: "http://downsc.chinaz.net/Files/DownLoad/sound1/201906/11582.mp3",
		ChMp3Url:  "http://downsc.chinaz.net/Files/DownLoad/sound1/201906/11582.mp3",
		WordEng:   "test",
		WordCh:    "测试",
	}
	err = c.UpdateWordInfo(&wc)
	a.Nil(err)
	wi, err = c.GetOneWordDetail("one")
	a.Nil(wi)
	a.Nil(err)
	wi, err = c.GetOneWordDetail("test")
	a.NotNil(wi)
	a.Nil(err)
	wcf := c.wordCacheFile("test")
	a.NotNil(wcf)
	a.Equal(wcf.EngMp3Path, wi.EngMp3)
	a.Equal(wcf.ChMp3Path, wi.ChMp3)
	a.Equal("测试", wi.Chinese)
	a.Equal("test", wi.English)

	bs1, err := os.ReadFile(wi.EngMp3)
	a.Nil(err)
	bs2, err := os.ReadFile(wi.ChMp3)
	a.Nil(err)
	a.Equal(bs1, bs2)
}
