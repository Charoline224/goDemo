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
	// TODO : A +trace根服务器查询直接给我答案了，但是NS +trace可以返回trace路径
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
		server = "199.9.14.201"
	}
	if qtype == 0 {
		qtype = 1
	}
	ress, err := dns.Resolve(domain, server, qtype, trace)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, res := range ress {
		fmt.Println(res)
	}
}
