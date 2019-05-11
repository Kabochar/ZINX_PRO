package znet

import (
	"ZINX_PRO/zinx/utils"
	"ZINX_PRO/zinx/ziface"
	"bytes"
	"encoding/binary"
	"errors"
)

// 封包，拆包具体模块
type DataPack struct{}

// 拆包‘封包实例的初始化方法
func NewDataPack() *DataPack {
	return &DataPack{}
}

// 获取包 头的长度
func (dp *DataPack) GetHandLen() uint32 {
	// DataLen uint32(4字节) + ID uint32(4字节)
	return 8
}

// 封包方法
// 数据格式：|datalen|msgID|data
func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	// 创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	// 将 datalen 写进 databuff
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}
	// 将 MsgId 写进 databuff
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgID()); err != nil {
		return nil, err
	}
	// 将 data 写进 databuff
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

// 拆包方法 (将包的Hand信息都出来) 之后，再根据 head信息里的data的长度，再进行读取
func (dp *DataPack) Unpack(binaryData []byte) (ziface.IMessage, error) {
	// 创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)

	// 只解压head信息，得到 datalen 和 megID
	msg := &Message{}

	// 读取 datalen
	// &msg.DataLen: 使用 取地址 & 的原因：需要修改值
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	// 读取 msgID
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.ID); err != nil {
		return nil, err
	}

	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("too large msg data recv")
	}

	return msg, nil
}
