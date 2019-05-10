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

// PingHandle test Handle
func (pr *PingRouter) PingHandle(request ziface.IRequest) {
	fmt.Println("Call Router PingHandle...")

	// 先读取客户端的数据，再回写ping...ping...ping
	fmt.Println("recv from client: msgID = ", request.GetMsgID(),
		", data = ", string(request.GetData()))

	// 发送消息给客户端
	err := request.GetConnection().SendMsg(200, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
	}
}

// HelloRouter test 自定义路由
type HelloRouter struct {
	znet.BaseRouter
}

// test Handle
func (pr *HelloRouter) HelloHandle(request ziface.IRequest) {
	fmt.Println("Call Router HelloHandle...")

	// 先读取客户端的数据，再回写ping...ping...ping
	fmt.Println("recv from client: msgID = ", request.GetMsgID(),
		", data = ", string(request.GetData()))

	// 发送消息给客户端
	err := request.GetConnection().SendMsg(201, []byte("bee...bee...bee"))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	// 1 创建一个server句柄，使用Zinx的 api
	s := znet.NewServer("[Zinx V0.6]")

	// 2 给当前的zinx框架添加一个自定义的router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &PingRouter{})

	// 3 启动server
	s.Serve()
}
