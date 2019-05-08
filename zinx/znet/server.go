package znet

import (
	"ZINX_PRO/zinx/ziface"
	"fmt"
	"github.com/kataras/iris/core/errors"
	"net"
)

// IServer接口实现，定义一个Server的服务器模块
type Server struct {
	// server name
	Name string
	// IP version
	IPVersion string
	// listen IP
	IP string
	// listen PORT
	Port int
}

func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
	// 回显的业务
	fmt.Println("[Conn Handle] CallbackToClient...")
	if _, err := conn.Write(data[:cnt]); err != nil {
		fmt.Println("write back buf err", err)
		return errors.New("CallBackToClient error")
	}

	return nil
}

// 启动服务器
func (s *Server) Start() {
	fmt.Printf("[Start] Server Listenner at IP :%s, Port %d, is starting\n", s.IP, s.Port)

	go func() {
		// 1 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr err: ", err)
			return
		}

		// 2 监听服务器的地址
		listenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen ", s.IPVersion, " err ", err)
			return
		}
		var cid uint32
		cid = 0

		fmt.Println("start Zinx server succ, ", s.Name, " succ , Listenning...")

		// 3 堵塞等待客户端连接，处理客户端链接业务（读写）
		for {
			// 如果有客户端链接，阻塞会返回
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err: ", err)
				continue
			}

			// 使用 connections 模块
			// connection 模块绑定 server
			// 将处理新连接的业务方法和conn进行绑定，得到我们的链接模块
			delConn := NewConnection(conn, cid, CallBackToClient)
			cid++

			// 启动当前的连接业务处理
			go delConn.Start()
		}
	}()

}

// 停止服务器
func (s *Server) Stop() {
	// TODO 将一些服务器的资源、状态 or 一些已经开辟的链接信息，进行停止or回收
}

// 运行服务器
func (s *Server) Serve() {
	// 启动server的服务功能
	s.Start()
	// TODO 做一些启动服务器之后的额外业务

	// 堵塞状态
	select {}
}

/*
	初始化Server模块的方法
*/
func NewServer(name string) ziface.IServer {
	server := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
	}

	return server
}
