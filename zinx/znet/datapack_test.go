package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

// 只是负责测试datapack拆包，封包的单元测试
func TestNewDataPack(t *testing.T) {
	/*
		模拟 服务器
	*/
	// 1 创建 socketTCP
	listenner, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("server listen err:", err)
		return
	}
	// 创建一个go 承载 负责从客户端处理业务
	go func() {
		// 2 从客户端读取数据，拆包处理
		for {
			conn, err := listenner.Accept()
			if err != nil {
				fmt.Println("server accept error: ", err)
			}

			go func(conn net.Conn) {
				// 处理客户端的请求
				// ----- 拆包过程 -----
				// 定义一个拆包对象
				dp := NewDataPack()

				for {
					// 1 第一次从 conn 读，把包的 head 读出来
					headData := make([]byte, dp.GetHandLen())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read head err", err)
						break
					}

					msgHead, err := dp.Unpack(headData)
					if err != nil {
						fmt.Println("server unpack err: ", err)
						return
					}
					// 判断包中是否有数据 data
					if msgHead.GetDataLen() > 0 {
						// msg 有数据，进行第二次读取
						// 2 第二次从 conn 读，根据 head中的dataLen，再读取data 内容
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetDataLen())

						// 根据 datalen 的长度再次从 io 流中读取
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("server unpack data err:", err)
							return
						}

						// 完整的消息已经读取完毕
						fmt.Printf("---> Recv MsgID: %v, DataLen: %v, Data: %v \n",
							msg.ID, msg.DataLen, string(msg.Data))
					}

				}
			}(conn)
		}
	}()

	/*
		模拟客户端
	*/
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client dial err: ", err)
		return
	}

	// 创建一个封包对象 dp
	dp := NewDataPack()

	// 模拟粘包过程，封装两个 msg 一同发送
	// 封装第一个msg1 包
	msg1 := &Message{
		ID:      1,
		DataLen: 5,
		Data:    []byte{'z', 'i', 'n', 'x', '!'},
	}

	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack msg1 error:", err)
		return
	}

	// 封装第二个msg2 包
	msg2 := &Message{
		ID:      1,
		DataLen: 7,
		Data:    []byte{'n', 'i', ',', 'h', 'a', 'o', '!'},
	}

	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client pack msg1 error:", err)
		return
	}

	// 将两个包粘在一起
	sendDataFull := append(sendData1, sendData2...)

	// 一次性发送给服务端
	conn.Write(sendDataFull)

	// 客户端堵塞
	select {}
}
