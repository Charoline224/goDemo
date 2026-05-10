package clock

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

var mu sync.Mutex
var times map[string]string

func main() {
	//读入参数，不断更新map
	for _, arg := range os.Args[1:] {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) != 2 {
			log.Fatalf("invalid argument: %s", arg)
		}
		name, address := parts[0], parts[1]
		go watch(name, address)
	}
	//输出
	for {
		time.Sleep(1 * time.Second)
		mu.Lock()
		for name, time := range times {
			fmt.Printf("%-10s %s\n", name, time)
		}
		mu.Unlock()
	}
}

func watch(name, address string) {
	conn, err := net.Dial("TCP", address)
	if err != nil {
		log.Fatal("无法连接")
		return
	}
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		mu.Lock()
		times[name] = scanner.Text()
		mu.Unlock()
	}
}
