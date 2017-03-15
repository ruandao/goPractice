package main

import (
	"math/rand"
	"time"
	"bufio"
	"os"
	"fmt"
)

const (
	duplicateVowel bool = true
	removeVowel		bool = false
)

func randBool() bool {
	return rand.Intn(2) == 0
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		word := []byte(s.Text())
		if randBool() {
			var v1 int = -1
			for i, char := range word {
				switch char {
				case 'a', 'e', 'i', 'o', 'u', 'A', 'E', 'I', 'O', 'U':
					if randBool() {
						v1 = i
					}
				}
			}
			if v1 >= 0 {
				switch randBool() {
				case duplicateVowel:
					word = append(word[:v1+1], word[v1:]...)
				case removeVowel:
					word = append(word[:v1], word[v1+1:]...)
				}
			}
		}
		fmt.Println(string(word))
	}
}
