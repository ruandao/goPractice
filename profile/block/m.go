package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"time"
)

func main() {
	runtime.SetBlockProfileRate(1)
	for i:=0; i< 1000; i++ {
		go SleepTest()
		go TickerTest()
	}
	log.Fatalln(http.ListenAndServe(":8080", nil))
}

func SleepTest() {
	for {
		time.Sleep(10 * time.Second)
		log.Printf("sleep\n")
	}
}

func TickerTest() {
	t := time.NewTicker(10 * time.Second)
	defer t.Stop()

	for range t.C {
		log.Printf("tick\n")
	}
}