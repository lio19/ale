package main

import (
	"encoding/json"
	"io/ioutil"
)

type Conf struct {
	C map[string]string
}

var conf Conf

func init() {
	bs, err := ioutil.ReadFile("conf.json")
	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal(bs, &conf.C); err != nil {
		panic(err)
	}
}

func (c *Conf) getValue(key string) string {
	return c.C[key]
}


