package ziface

/*
	将请求的消息封装到一个 Message 中
*/
type IMessage interface {
	// 获取消息ID
	GetMsgID() uint32
	// 获取消息长度
	GetDataLen() uint32
	// 获取消息内容
	GetData() []byte

	// 设置消息ID
	SetMsgID(uint32)
	// 设置消息内容
	SetData([]byte)
	// 设置消息长度
	SetDataLen(uint32)
}
