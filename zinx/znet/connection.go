package znet

import (
	"ZINX_PRO/zinx/utils"
	"ZINX_PRO/zinx/ziface"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
)

/**
连接模块
*/
type Connection struct {
	// 当前Conn属于哪个Server
	// //当前conn属于哪个server，在conn初始化的时候添加即可
	TcpServer ziface.IServer
	// 从当前连接的 socket TCP 套接字
	Conn *net.TCPConn
	// 连接 ID
	ConnID uint32
	// 当前的连接状态
	isClosed bool
	//消息的管理MsgID 和对应的处理业务API关系
	MsgHandler ziface.IMsgHandle
	// 告知当前连接一斤推出/停止 channel.(由Reader告知Writer退出)
	ExitChan chan bool
	// 无缓冲管道，用于 读、写两个goroutine之间的通信
	msgChan chan []byte
	// 链接属性集合
	property map[string]interface{}
	// 保护链接属性的锁
	propertyLock sync.RWMutex
}

func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandle) *Connection {
	c := &Connection{
		TcpServer:  server, //将隶属的server传递进来
		Conn:       conn,
		ConnID:     connID,
		MsgHandler: msgHandler,
		isClosed:   false,
		ExitChan:   make(chan bool, 1),
		msgChan:    make(chan []byte),
		property:   make(map[string]interface{}),
	}

	// 将新创建的Conn添加到链接管理中
	c.TcpServer.GetConnMgr().Add(c)

	return c
}

// 连接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("[Reader Goroutine is running]")
	defer fmt.Printf("ConnID: %v, Reader is exit, Remote Addr: %v \n",
		c.ConnID, c.RemoteAddr().String())
	defer c.Stop()

	for {

		dp := NewDataPack()
		// 读取客户端msg head 二进制流前 8字节
		headData := make([]byte, dp.GetHandLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head err: ", err)
			break
		}

		// 拆包，得到 msgID 和 msgDataLen 放在 msg消息中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error: ", err)
			break
		}

		// 根据datalen 再次读取Data，放在 msg.Data中
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data err:", err)
				break
			}
		}
		msg.SetData(data)

		// 得到当前客户端请求的 Request 数据
		req := Request{
			conn: c,
			msg:  msg,
		}
		if utils.GlobalObject.WorkerPoolSize > 0 {
			// 已经开启了工作池机制，将消息发送给Worker工作池处理即可
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			// 从路由中，找到注册绑定的Conn对应的router调用
			// 根据绑定好的MsgID 找到对应处理api业务 执行
			go c.MsgHandler.DoMsgHandler(&req)
		}

	}
}

// 写消息goroutine，用户将数据发送给客户端
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println("[Writer exit!]", c.RemoteAddr().String())

	// 不断堵塞等待chanel的消息，进行写给客户端
	for {
		select {
		case data := <-c.msgChan:
			// 有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data errror, ", err)
				return
			}
		case <-c.ExitChan:
			// 代表reader已经退出，此时Writer也要退出
			fmt.Println("[Reader exit!]", c.RemoteAddr().String())
			return
		}
	}
}

// 启动链接，让当前连接开始工作
func (c *Connection) Start() {
	fmt.Println("Conn Start()... ConnID = ", c.ConnID)

	// 启动从当前连接的读数据业务
	go c.StartReader()

	// 启动从当前连接写数据的业务
	go c.StartWriter()

	// 按照用户传递进来的创建连接时需要处理的业务，执行钩子方法
	c.TcpServer.CallOnConnStart(c)
}

// 停止连接，结束当前连接状态
func (c *Connection) Stop() {
	fmt.Println("Conn Stop().. ConnID = ", c.ConnID)

	// 如果连接已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	//如果用户注册了该链接的关闭回调业务，那么在此刻应该显示调用
	c.TcpServer.CallOnConnStop(c)

	// 关闭socket连接
	c.Conn.Close()

	// 告知Writer关闭
	c.ExitChan <- true

	// 将当前conn从ConnMgr中摘除掉
	c.TcpServer.GetConnMgr().Remove(c)

	// 回收资源
	close(c.ExitChan)
	close(c.msgChan)
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

// 发送数据，将数据发送给远程的客户端.先封包，后发送
func (c *Connection) SendMsg(msgID uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closed when send msg")
	}

	// 将data进行封包
	// MsgDataLen|MsgID|Data
	dp := NewDataPack()

	// MsgID | MsgDatalen | Data 封包
	binaryMsg, err := dp.Pack(NewMsgPackage(msgID, data))
	if err != nil {
		fmt.Println("Pack err msg", err)
		return errors.New("Pack err Msg")
	}

	// 封包数据写进 msgChan 中
	c.msgChan <- binaryMsg

	return nil
}

// 设置链接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = value
}

// 获取链接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	// 读取属性
	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}

// 移除链接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	// 删除属性
	delete(c.property, key)
}
