package main

import (
	"fmt"
	"mydig/dns"
	"os"
	"strings"
)

func main() {
	var domain, server string
	var qtype uint16
	var trace bool
	for _, s := range os.Args[1:] {
		if strings.HasPrefix(s, "@") {
			server = strings.TrimPrefix(s, "@")
		}
		if strings.HasPrefix(s, "+") {
			switch s {
			case "+trace":
				trace = true
			}
		}
		if s == "A" || s == "a" {
			qtype = 1
		}
		if s == "NS" || s == "ns" {
			qtype = 2
		}
		if s == "CNAME" || s == "cname" {
			qtype = 5
		} else if !strings.HasPrefix(s, "@") && !strings.HasPrefix(s, "+") && s != "A" && s != "a" && s != "NS" && s != "ns" && s != "CNAME" && s != "cname" {
			domain = s
		}
	}
	if server == "" {
		server = "198.41.0.4"
	}
	if qtype == 0 {
		qtype = 1
	}
	fmt.Println("domain:", domain, "server:", server, "qtype:", qtype)
	ress, err := dns.Resolve(domain, server, qtype, trace)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, res := range ress {
		fmt.Println(res)
	}
}
