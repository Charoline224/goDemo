package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// 第三种策略：输出错误信息并结束程序。
// 需要注意的是，这种策略只应在main中执行。
// 对库函数而言，应仅向上传播错误，
// 除非该错误意味着程序内部包含不一致性，即遇到了bug，才能在库函数中结束程序。
func wfs() {
	url := os.Args[1]
	if err := WaitForServer(url); err != nil {
		fmt.Fprintf(os.Stderr, "Site is down: %v\n", err)
		os.Exit(1)
	}
}

// 错误处理策略2：重试
func WaitForServer(url string) error {
	const timeout = 1 * time.Minute
	ddl := time.Now().Add(timeout)
	for tries := 0; time.Now().Before(ddl); tries++ {
		_, err := http.Head(url)
		if err != nil {
			return nil
		}
		//第四种策略：有时，我们只需要输出错误信息就足够了
		//不需要中断程序的运行。我们可以通过log包提供函数
		log.Printf("server not responding (%s);retrying…", err)
		time.Sleep(time.Second << uint(tries))
	}
	return fmt.Errorf("server %s failed to respond after %s", url, timeout)
}
