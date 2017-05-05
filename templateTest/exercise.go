package main

import (
	"os"
	"github.com/ruandao/goPractice/templateTest/tmpl"
)

type Inventory struct {
	Material string
	Count  uint
}
func main() {
	tmpl.AddDir("./templates")
	err := tmpl.Execute(os.Stdout, "hello/content", nil)
	if err != nil {
		panic(err)
	}
}
