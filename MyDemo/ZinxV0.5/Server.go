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

// test Handle
func (pr *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle...")

	// 先读取客户端的数据，再回写ping...ping...ping
	fmt.Println("recv from client: msgID = ", request.GetMsgID(),
		", data = ", string(request.GetData()))

	// 发送消息给客户端
	err := request.GetConnection().SendMsg(1, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
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
