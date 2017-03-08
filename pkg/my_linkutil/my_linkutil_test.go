package my_linkutil

import (
	"testing"
	"fmt"
	"strings"
)

func TestLinksFromURL(t *testing.T) {
	urls, err := LinksFromURL("http://www.qq.com")
	if err != nil {
		t.Errorf("err %s\n", err)
	}
	fmt.Printf("total link: %d\n", len(urls))
	for _, url := range urls {
		fmt.Println(url)
	}
}

func TestLinksFromReader(t *testing.T) {
	content := `<a href="http://www.baidu.com">百度</a>`
	r := strings.NewReader(content)
	urls, err := LinksFromReader(r)
	if err != nil {
		t.Errorf("err %s\n", err)
	}
	fmt.Printf("total link: %d\n", len(urls))
	for _, url := range urls {
		fmt.Printf("%s\n", url)
	}
}