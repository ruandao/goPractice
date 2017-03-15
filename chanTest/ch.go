package main

func main() {
	ch := make(chan string)
	close(ch)
	ch <- "hello"
}
