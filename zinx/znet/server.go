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
	//当前server的消息管理模块， 用来绑定MsgID和对应的处理业务API关系
	MsgHandler ziface.IMsgHandle
	// 当前Server的连接管理器
	ConnMgr ziface.IConnManager
	// 该Server的连接创建时Hook函数
	OnConnStart func(conn ziface.IConnection)
	// 该Server的连接断开时的Hook函数
	OnConnStop func(conn ziface.IConnection)
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
		// 0 开启消息队列及worker工作池
		s.MsgHandler.StartWorkerPool()

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

			// 设置最大连接个数的判断，如果超过最大连接，那么则关闭此新的连接
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				// TODO 给客户端相应一个超出最大连接的错课包
				fmt.Println("[STOP] Zinx server , name ", s.Name)
				conn.Close()
				continue
			}

			// 处理新连接业务方法 和 conn 进行绑定，得到我们的链接模块
			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++

			// 启动当前的连接业务处理
			go dealConn.Start()
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
func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("Add Router Succ!")
}

/*
	初始化Server模块的方法
*/
func NewServer(name string) ziface.IServer {
	server := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(),
	}

	return server
}

// 得到链接管理
func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

/*
	HOOK 钩子函数
*/
// 设置该Server的连接创建时Hook函数
func (s *Server) SetOnConnStart(hookFunc func(ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

// 设置该Server的连接断开时的Hook函数
func (s *Server) SetOnConnStop(hookFunc func(ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

// 调用连接OnConnStart Hook函数
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("---> CallOnConnStart....")
		s.OnConnStart(conn)
	}
}

// 调用连接OnConnStop Hook函数
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("---> CallOnConnStop....")
		s.OnConnStop(conn)
	}
}
