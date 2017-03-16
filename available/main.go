package main

import (
	"net"
	"bufio"
	"strings"
	"os"
	"fmt"
	"time"
)

func exists(domain string) (bool, error) {
	const whoisServer string = "com.whois-servers.net"
	conn, err := net.Dial("tcp", whoisServer + ":43")
	if err != nil {
		return false, err
	}
	defer conn.Close()
	conn.Write([]byte(domain + "rn"))
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		if strings.Contains(strings.ToLower(scanner.Text()), "no match") {
			return false, nil
		}
	}
	return true, nil
}

var marks = map[bool]string{true: "ok", false : "fail"}
func main() {
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		domain := s.Text()
		exist, err := exists(domain)
		if err != nil {
			fmt.Printf("domain: %s err %s\n", domain, err)
			continue
		}
		fmt.Printf("domain: %s available: %s\n", domain, marks[!exist])
		time.Sleep(1 * time.Second)
	}
}
