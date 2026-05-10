package clock

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

func clock1() {

	port := flag.String("port", "8080", "port to listen on")
	flag.Parse()
	listener, err := net.Listen("TCP", "localhost:"+*port)
	if err != nil {
		fmt.Fprintf(os.Stderr, "端口无法监听:%s", err)
		return
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "连接失败:%s", err)
			continue
		}
		go handle(conn)
	}

}

func handle(c net.Conn) {
	defer c.Close()
	for {
		//由于net.Conn实现了io.Writer接口，我们可以直接向其写入内容
		_, err := io.WriteString(c, os.Stdin.Name())
		if err != nil {
			return
		}
		time.Sleep(1 * time.Second)
	}

}
