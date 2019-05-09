package znet

import (
	"ZINX_PRO/zinx/utils"
	"ZINX_PRO/zinx/ziface"
	"fmt"
	"net"
)

/**
连接模块
*/
type Connection struct {
	// 从当前连接的 socket TCP 套接字
	Conn *net.TCPConn
	// 连接 ID
	ConnID uint32
	// 当前的连接状态
	isClosed bool
	// 该链接的处理方法router
	Router ziface.IRouter
	// 告知当前连接一斤推出/停止 channel
	ExitChan chan bool
}

func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connID,
		Router:   router,
		isClosed: false,
		ExitChan: make(chan bool, 1),
	}

	return c
}

// 连接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running..")
	defer fmt.Println("ConnID = ", c.ConnID, "Reader is exit, Remote Addr is ", c.RemoteAddr().String())
	defer c.Stop()

	for {
		// 读取客户端的数据到buf中，最大 512字节
		buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf err", err)
			continue
		}

		// 得到当前客户端请求的 Request 数据
		req := Request{
			conn: c,
			data: buf,
		}

		// 从路由 Routers 中找到注册绑定Conn的对应 Handle
		go func(request ziface.IRequest) {
			// 执行注册的路由方法
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)
	}
}

// 启动链接，让当前连接开始工作
func (c *Connection) Start() {
	fmt.Println("Conn Start()... ConnID = ", c.ConnID)

	// 启动从当前连接的读数据业务
	go c.StartReader()

	// TODO 启动从当前连接写数的业务
}

// 停止连接，结束当前连接状态
func (c *Connection) Stop() {
	fmt.Println("Conn Stop().. ConnID = ", c.ConnID)

	// 如果连接已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	// 关闭socket连接
	c.Conn.Close()
}

// 从当前连接获取原始的 socket TCPConn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// 获取当前连接 ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// 获取远程客户端地址信息
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// 发送数据，将数据发送给远程的客户端
func (c *Connection) Send(data []byte) error {
	return nil
}
