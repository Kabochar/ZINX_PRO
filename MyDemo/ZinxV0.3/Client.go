package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	fmt.Println("Client strat...")

	time.Sleep(1 * time.Second)

	// 1 直接连接远程服务器，得到一个conn链接
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("net dial err", err)
		return
	}

	for {
		// 2 链接调用Write 写数据
		_, err := conn.Write([]byte("Hello Zinx V0.1"))
		if err != nil {
			fmt.Println("write conn err", err)
			return
		}

		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf err", err)
			return
		}

		fmt.Printf("server call back: %s, cnt = %d\n", buf[:cnt], cnt)

		// CPU堵塞：让CPU释放资源，避免for{}耗光资源
		time.Sleep(1 * time.Second)
	}

}
