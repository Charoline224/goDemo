package main

import (
	"fmt"
	"os"
	"time"
)

func countdown() {
	fmt.Println("Commencing countdown.  Press return to abort.")
	abort := make(chan int)
	go func() {
		os.Stdin.Read(make([]byte, 1))
		abort <- 1
	}()
	tick := time.Tick(1 * time.Second)
	for i := 10; i > 0; i-- {
		fmt.Println(i)
		select {
		case <-tick:
		case <-abort:
			fmt.Printf("stop")
			return
		}
	}
	// lauch()发射，不过没这个函数（
}
