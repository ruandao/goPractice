package main

import (
	"math/rand"
	"time"
	"bufio"
	"os"
	"strings"
	"unicode"
	"fmt"
)

var tlds = []string{"com", "net"}

const allowedChars = "qwertyuiopasdfghjklzxcvbnm1234567890_-"

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		text := strings.ToLower(strings.TrimSpace(s.Text()))
		var newText []rune
		for _, r := range text {
			if unicode.IsSpace(r) {
				r = '-'
			}
			if !strings.ContainsRune(allowedChars, r) {
				continue
			}
			newText = append(newText, r)
		}
		fmt.Println(string(newText) + "." + tlds[rand.Intn(len(tlds))])
	}
}
