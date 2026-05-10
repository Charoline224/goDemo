package main

import (
	"functionDemo/links"
	"os"
)

var tokens = make(chan int, 20)

func crawl(url string) []string {
	tokens <- 1
	list, err := links.Extract(url)
	<-tokens
	if err != nil {
		return nil
	}
	return list
}

func work() {
	// n:表示还有多少任务，每向worklist里面输入一个，任务就加一
	var n = 0
	worklist := make(chan []string)

	//任务的增加和处理是并发的，任务放在worklist里面
	n++
	go func() { worklist <- os.Args[1:] }()
	//标记
	seen := make(map[string]bool)
	// for _, link := range <-worklist{
	// 	if seen[link]{
	// 		continue
	// 	}
	// 	seen[link] = true
	// 	n++
	// 	list = append(list, crawl(link))
	// }
	// 用n控制处理任务，直到n=0就没有任务了
	for ; n > 0; n-- {
		//收到一批链接，处理
		list := <-worklist
		for _, link := range list {
			link := link
			if seen[link] {
				continue
			}
			seen[link] = true
			n++
			// link 是循环变量，如果不作为参数传入，
			// 所有 goroutine 会共享同一个 link，
			// 等 goroutine 真正执行时，
			// link 已经变成循环最后一个值了
			go func(link string) {
				worklist <- crawl(link)
			}(link)
		}
	}
}
