package main

import (
	"os"
	"github.com/ruandao/goPractice/thesaurus"
	"bufio"
	"log"
	"fmt"
)

func main() {
	apiKey := os.Getenv("BHT_APIKEY")
	thesaurus := &thesaurus.BigHuge{APIKey:apiKey}
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		word := s.Text()
		syns, err := thesaurus.Synonyms(word)
		if err != nil {
			log.Fatal("Failed when looking for synonyms for " + word + err.Error())
		}
		if len(syns) == 0 {
			log.Fatal("Couldn't find any synonyms for " + word)
		}
		for _, syn := range syns {
			fmt.Println(syn)
		}
	}
}
