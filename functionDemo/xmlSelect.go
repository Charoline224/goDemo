package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"
)

func xmlSelect() {
	dec := xml.NewDecoder(os.Stdin)
	var stack []string
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "input something invalid:%v", err)
		}
		switch tok := tok.(type) {
		case xml.StartElement:
			stack = append(stack, tok.Name.Local)
		case xml.EndElement:
			stack = stack[:len(stack)-1]
		case xml.CharData:
			if contains(stack, os.Args[1:]) {
				fmt.Printf("%s has : %s", strings.Join(stack, " "), tok)
			}
		}
	}
}

func contains(s, t []string) bool {
	for len(s) >= len(t) {
		if len(t) == 0 {
			return true
		}
		if s[0] == t[0] {
			t = t[1:]
		}
		s = s[1:]
	}
	return false
}
