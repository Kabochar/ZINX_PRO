package znet

import (
	"ZINX_PRO/zinx/ziface"
	"errors"
	"fmt"
	"io"
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
	//消息的管理MsgID 和对应的处理业务API关系
	MsgHandler ziface.IMsgHandle
	// 告知当前连接一斤推出/停止 channel.(由Reader告知Writer退出)
	ExitChan chan bool
	// 无缓冲管道，用于 读、写两个goroutine之间的通信
	msgChan chan []byte
}

func NewConnection(conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandle) *Connection {
	c := &Connection{
		Conn:       conn,
		ConnID:     connID,
		MsgHandler: msgHandler,
		isClosed:   false,
		ExitChan:   make(chan bool, 1),
		msgChan:    make(chan []byte),
	}

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

		// 从路由中，找到注册绑定的Conn对应的router调用
		// 根据绑定好的MsgID 找到对应处理api业务 执行
		go c.MsgHandler.DoMsgHandler(&req)
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

	// 告知Writer关闭
	c.ExitChan <- true

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
