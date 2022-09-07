package main

import "autoLearnEnglish/src/conf"

func main() {
	conf.Init()
	ale := NewALE()
	ale.Start()
}
