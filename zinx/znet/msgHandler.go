package znet

import (
	"ZINX_PRO/zinx/ziface"
	"fmt"
	"strconv"
)

/*
	消息处理模块的实现
*/
type MsgHandle struct {
	//存放每个MsgID 所对应的处理方法
	APIs map[uint32]ziface.IRouter
}

// 初始化/创建MsgHandle方法
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		APIs: make(map[uint32]ziface.IRouter),
	}
}

// 调度/执行对应的Router消息处理方法
func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	// 1 从 Request 中找到 msgID
	handler, ok := mh.APIs[request.GetMsgID()]
	if !ok {
		fmt.Printf("API msgID = %v, is NOT FOUND! Need Register!\n", request.GetMsgID())
	}
	// 2 根据MsgID 调度对应router业务即可
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// 为消息添加具体的处理逻辑
func (mh *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	// 1 判断 当前msg绑定的API处理方法是否已经存在
	if _, ok := mh.APIs[msgID]; ok {
		// id已经注册了
		panic("repeat api , msgID = " + strconv.Itoa(int(msgID)))
	}
	// 2 添加msg与API的绑定关系
	mh.APIs[msgID] = router
	fmt.Println("Add api MsgID = ", msgID, " succ!")
}
