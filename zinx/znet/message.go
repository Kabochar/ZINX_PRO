package znet

type Message struct {
	ID      uint32 // 消息的ID
	DataLen uint32 // 消息长度
	Data    []byte // 消息内容
}

// 创建一个Message消息包
func NewMsgPackage(id uint32, data []byte) *Message {
	return &Message{
		ID:      id,
		DataLen: uint32(len(data)),
		Data:    data,
	}
}

// 获取消息ID
func (m *Message) GetMsgID() uint32 {
	return m.ID
}

// 获取消息长度
func (m *Message) GetMsgLen() uint32 {
	return m.DataLen
}

// 获取消息内容
func (m *Message) GetData() []byte {
	return m.Data
}

// 设置消息ID
func (m *Message) SetMsgID(id uint32) {
	m.ID = id
}

// 设置消息内容
func (m *Message) SetData(data []byte) {
	m.Data = data
}

// 设置消息长度
func (m *Message) SetDataLen(len uint32) {
	m.DataLen = len
}
