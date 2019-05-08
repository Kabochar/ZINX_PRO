package main

import "ZINX_PRO/zinx/znet"

/*
	基于 Zinx框架来开发，服务端相应程序
*/
func main() {
	// 1 创建一个server句柄，使用Zinx的 api
	server := znet.NewServer("[Zinx V0.2]")
	// 2 启动server
	server.Serve()
}
