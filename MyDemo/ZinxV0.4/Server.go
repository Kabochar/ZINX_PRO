package main

import (
	"ZINX_PRO/zinx/ziface"
	"ZINX_PRO/zinx/znet"
	"fmt"
)

/*
	基于 Zinx框架来开发，服务端相应程序
*/

// ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

// test PreHandle
func (pr *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("Call Router PreHandle...")
	_, err := request.GetConnection().GetTCPConnection().
		Write([]byte("before ping...\n"))
	if err != nil {
		fmt.Println("call back: before ping error")
	}
}

// test Handle
func (pr *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle...")
	_, err := request.GetConnection().GetTCPConnection().
		Write([]byte("current ping...\n"))
	if err != nil {
		fmt.Println("call back: current ping error")
	}
}

// test PostHandle
func (pr *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("Call Router PostHandle...")
	_, err := request.GetConnection().GetTCPConnection().
		Write([]byte("after ping...\n"))
	if err != nil {
		fmt.Println("call back: after ping error")
	}
}
func main() {
	// 1 创建一个server句柄，使用Zinx的 api
	s := znet.NewServer("[Zinx V0.4]")
	// 2 给当前的zinx框架添加一个自定义的router
	s.AddRouter(&PingRouter{})
	// 3 启动server
	s.Serve()
}
