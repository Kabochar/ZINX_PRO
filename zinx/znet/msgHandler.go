package znet

import (
	"ZINX_PRO/zinx/ziface"
	"fmt"
	"strconv"
)

/*
	消息处理模块的实现
*/

type MesHandle struct {
	// 存放每个MsgID所对应的方法
	Apis map[uint32]ziface.IRouter
}

// 初始化/创建MsgHandle方法
func NewMsgHandle() *MesHandle {
	return &MesHandle{
		Apis: make(map[uint32]ziface.IRouter),
	}
}

// 调度/执行对应的Router消息处理方法
func (mh *MesHandle) DoMsgHandler(request ziface.IRequest) {
	// 1 从 Request 中找到 msgID
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Printf("api msgID = %s is NOT FOUND! Need register!\n",
			request.GetMsgID())
	}
	// 2 根据MsgID 调度对应router业务即可
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// 为消息添加具体的处理逻辑
func (mh *MesHandle) Addrouter(msgID uint32, router ziface.IRouter) {
	// 1 判断 当前 msg 绑定的API处理防范是否已经注册
	if _, ok := mh.Apis[msgID]; ok {
		panic("repeat api, msgID = " + strconv.Itoa(int(msgID)))
	}

	// 2 添加msg与API绑定关系
	mh.Apis[msgID] = router
	fmt.Printf("Add api MsgID = %v\n", msgID)
}
