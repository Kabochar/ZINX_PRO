package main

import (
	"ZINX_PRO/zinx/znet"
	"fmt"
	"io"
	"net"
	"time"
)

func main() {
	fmt.Println("Client strat...")

	time.Sleep(1 * time.Second)

	// 1 直接连接远程服务器，得到一个conn链接
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("net dial err", err)
		return
	}

	for {
		// 发送封包的message的消息
		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(znet.NewMsgPackage(0, []byte("ZinxV0.5 client Test Message")))
		if err != nil {
			fmt.Println("Pack error:", err)
			return
		}
		// 发送消息
		if _, err = conn.Write(binaryMsg); err != nil {
			fmt.Println("write error:", err)
			return
		}

		// 服务器应回复一个messsage数据，Msg: 1 ping

		// 1 先读取流中的head部分，得到id 和 datalen
		binaryHead := make([]byte, dp.GetHandLen())
		if _, err := io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("read head error:", err)
			break
		}

		// 将二进制的head拆包到msg结构体重
		msgHead, err := dp.Unpack(binaryHead)

		if msgHead.GetMsgLen() > 0 {
			// 2 再根据datalen 进行第二次读取，将data 读出来
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetMsgLen())

			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("read msg data err:", err)
				return
			}

			fmt.Printf("---> Recv Server MsgID: %d, len: %d. Data: %s\n",
				msg.ID, msg.DataLen, string(msg.Data))
		}

		// CPU堵塞：让CPU释放资源，避免for{}耗光资源
		time.Sleep(1 * time.Second)
	}

}
