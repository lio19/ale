package youdao

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func (wb *WordBook) signal(q, curTime, salt string) string {

	var input string
	if len(q) <= 10 {
		input = fmt.Sprintf("%s%d", q, len(q))
	} else if len(q) <= 20 {
		input = fmt.Sprintf("%s%d%s", q[:10], len(q), q[10:])
	} else {
		input = q
	}

	signStr := wb.appKey + input + salt + curTime + wb.appSec

	m := sha256.New()
	m.Write([]byte(signStr))
	return hex.EncodeToString(m.Sum(nil))
}

func (wb *WordBook) GetOneWord(q string) (*RespWord, error) {
	now := time.Now()
	curTime := strconv.FormatInt(now.Unix(), 10)
	salt := strconv.FormatInt(now.UnixMilli(), 10)
	sign := wb.signal(q, curTime, salt)

	resp, err := http.PostForm("", url.Values{
		"q":        {q},
		"from":     {"EN"},
		"to":       {"zh-CHS"},
		"appKey":   {wb.appKey},
		"salt":     {salt},
		"sign":     {sign},
		"signType": {"v3"},
		"curtime":  {curTime},
		"ext":      {"mp3"},
		"voice":    {"0"},
		"strict":   {"false"},
	})

	if err != nil {
		return nil, err
	}

	var rw RespWord
	bs, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bs, &rw)
	if err != nil {
		return nil, err
	}

	return &rw, nil
}
