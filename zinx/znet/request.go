package znet

import "ZINX_PRO/zinx/ziface"

type Request struct {
	// 已经和客户端建立好的链接
	conn ziface.IConnection
	// 客户端请求的数据
	msg ziface.IMessage
}

// 获取请求连接数据
func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

// 获取请求消息的数据
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

// 获取请求消息ID
func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgID()
}
