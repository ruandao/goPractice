package main
import (
	"fmt"
	"time"
)
func saferoutine(c chan bool) {
	for i := 0; i < 10; i++ {
		fmt.Println("Count:", i)
		time.Sleep(1 * time.Second)
	}
	c <- true
}
func panicgoroutine(c chan bool) {
	time.Sleep(5 * time.Second)
	panic("Panic, omg ...")
	c <- true
}
func main() {
	c := make(chan bool, 2)
	go saferoutine(c)
	go panicgoroutine(c)
	for i := 0; i < 2; i++ {
		<-c
	}
}
