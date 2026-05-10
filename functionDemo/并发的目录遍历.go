package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func walkDir(dir string, fileSize chan<- int64, n *sync.WaitGroup) {
	defer n.Done()
	for _, entry := range dirent(dir) {
		if entry.IsDir() {
			subdir := filepath.Join(dir, entry.Name())
			n.Add(1)
			go walkDir(subdir, fileSize, n)
		} else {
			info, err := entry.Info()
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v", err)
			}
			fileSize <- info.Size()
		}
	}
}

var sema = make(chan int, 20)

func dirent(dir string) []os.DirEntry {
	sema <- 1
	entries, err := os.ReadDir(dir)
	defer func() { <-sema }()
	if err != nil {
		fmt.Fprintf(os.Stderr, "du1: %v\n", err)
		return nil
	}
	return entries
}

var v = flag.Bool("v", false, "show verbose progress messages")

func main() {
	flag.Parse()
	roots := flag.Args()
	fileSizes := make(chan int64)
	//启用等待队列，用于计要完成的事务数量
	var n sync.WaitGroup
	for _, root := range roots {
		n.Add(1)
		go walkDir(root, fileSizes, &n)
	}
	go func() {
		n.Wait()
		close(fileSizes)
	}()

	//数filesize+控制打印过程
	var tick <-chan time.Time
	if *v {
		tick = time.Tick(500 * time.Millisecond)
	}

	var nfiles, nbytes int64
loop:
	for {
		select {
		case size, ok := <-fileSizes:
			if !ok {
				break loop
			}
			nfiles++
			nbytes += size
		case <-tick:
			fmt.Printf("%d files  %.1f GB\n", nfiles, float64(nbytes)/1e9)
		}
	}
}
