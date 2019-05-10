package ziface

/*
	IRequest 接口：
	实际上是把客户端请求的链接信息 和 请求的数据 包装到了 Request里
*/
type IRequest interface {
	// 获取请求连接数据
	GetConnection() IConnection

	// 获取请求消息的数据
	GetData() []byte

	// 获取请求消息ID
	GetMsgID() uint32
}
