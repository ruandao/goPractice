package my_linkutil

import "io"
import (
	"net/http"
	"io/ioutil"
	"regexp"
)

func LinksFromURL(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return LinksFromReader(resp.Body)
}

func LinksFromReader(r io.Reader) ([]string, error) {
	content, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	sContent := string(content)
	rep, err := regexp.Compile(`href="([^"]*)"`)
	if err != nil {
		return nil, err
	}
	urls := rep.FindAllStringSubmatch(sContent, -1)
	urls2 := make(map[string]bool)
	for _, url := range urls {
		urls2[url[1]] = true
	}
	allurls := make([]string, 0)
	for k,_ := range urls2 {
		allurls = append(allurls, k)
	}
	return allurls, nil
}