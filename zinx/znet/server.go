package znet

import (
	"ZINX_PRO/zinx/utils"
	"ZINX_PRO/zinx/ziface"
	"fmt"
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
	// 给当前Server 由用户绑定的回调Router，也就是server 注册的链接对应的处理业务
	// 当前只能存在一个
	Router ziface.IRouter
}

// 启动服务器
func (s *Server) Start() {
	fmt.Printf("[Zinx] Server Name: %s, listenner at IP: %s, Port: %d is starting...\n",
		utils.GlobalObject.Name,
		utils.GlobalObject.Host,
		utils.GlobalObject.TcpPort)

	fmt.Printf("[Zinx] Version %s, MaxConn: %d, MaxPackage: %d\n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPackageSize)

	fmt.Printf("[Start] Server Listenner at IP :%s, Port %d, is starting\n",
		s.IP, s.Port)

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

			// 处理新连接业务方法 和 conn 进行绑定，得到我们的链接模块
			delConn := NewConnection(conn, cid, s.Router)
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

// 路由功能：给当前服务注册一个路由业务功能，供客户端链接处理使用
func (s *Server) AddRouter(router ziface.IRouter) {
	s.Router = router
	fmt.Println("Add Router Succ!")
}

/*
	初始化Server模块的方法
*/
func NewServer(name string) ziface.IServer {
	server := &Server{
		Name:      utils.GlobalObject.Name,
		IPVersion: "tcp4",
		IP:        utils.GlobalObject.Host,
		Port:      utils.GlobalObject.TcpPort,
		Router:    nil,
	}

	return server
}
