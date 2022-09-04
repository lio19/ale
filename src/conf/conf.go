package conf

import (
	"encoding/json"
	"os"
)

type Conf struct {
	C map[string]string
}

var conf Conf

func Init() {
	bs, err := os.ReadFile("conf.json")
	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal(bs, &conf.C); err != nil {
		panic(err)
	}
}

func GetValue(key string) string {
	if conf.C == nil {
		panic("conf do not init")
	}
	return conf.C[key]
}

func MockConf(c map[string]string) {
	conf.C = c
}
