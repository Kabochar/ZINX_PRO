package ziface

import "net"

// 定义接口
type IConnection interface {
	// 启动链接，让当前连接开始工作
	Start()
	// 停止连接，结束当前连接状态
	Stop()
	// 从当前连接获取原始的 socket TCPConn
	GetTCPConnection() *net.TCPConn
	// 获取当前连接 ID
	GetConnID() uint32
	// 获取远程客户端地址信息
	RemoteAddr() net.Addr
	// 发送数据，将数据发送给远程的客户端.先封包，后发送
	SendMsg(uint32, []byte) error
}

type HandFunc func(*net.TCPConn, []byte, int) error
