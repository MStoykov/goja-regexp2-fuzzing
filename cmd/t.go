package main

import (
	"io/ioutil"
	"os"

	regexp "github.com/mstoykov/goja-regexp2-fuzzing"
)

func main() {
	b, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	regexp.Fuzz(b)
}
