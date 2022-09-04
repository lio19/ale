package youdao

import (
	"autoLearnEnglish/src/conf"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type WordBook struct {
	cookie string
	appKey string
	appSec string
}

func NewWordBook() *WordBook {
	return &WordBook{
		cookie: conf.GetValue("cookie"),
		appSec: conf.GetValue("appSec"),
		appKey: conf.GetValue("appKey"),
	}
}

func (wb *WordBook) Get() (*RespWordBookData, error) {
	//
	wbUrl := "https://dict.youdao.com/wordbook/webapi/v2/word/list?limit=48&offset=0&sort=time&lanTo=&lanFrom="

	req, err := http.NewRequest(http.MethodGet, wbUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", wb.cookie)
	req.Header.Set("Host", "dict.youdao.com")
	req.Header.Set("Origin", "https://www.youdao.com")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Referer", "https://www.youdao.com/")
	req.Header.Set("sec-ch-ua", `"Chromium";v="104", " Not A;Brand";v="99", "Google Chrome";v="104"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var respWB RespWordBook
	err = json.Unmarshal(respBytes, &respWB)
	if err != nil {
		return nil, err
	}

	if respWB.Code != 0 {
		return nil, errors.New(respWB.Msg)
	}

	return &respWB.Data, nil
}
