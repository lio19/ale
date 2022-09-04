package main

import "fmt"

type log struct {
}

var logger log

func (l *log) Info(i ...interface{}) {
	fmt.Println(i...)
}

func (l *log) Error(i ...interface{}) {
	fmt.Println(i...)
}

func (l *log) Warn(i ...interface{}) {
	fmt.Println(i...)
}
