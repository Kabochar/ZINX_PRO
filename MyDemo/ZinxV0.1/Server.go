package main

import "ZINX_PRO/zinx/znet"

func main() {
	server := znet.NewServer("[Zinx V0.1]")
	server.Serve()
}
