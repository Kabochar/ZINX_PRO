package znet

import (
	"ZINX_PRO/zinx/utils"
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
	// 负责Worker取任务的消息队列
	TaskQueue []chan ziface.IRequest
	// 业务工作Worker池的worker数量
	WorkerPoolSize uint32
}

// 初始化/创建MsgHandle方法
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		APIs:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize, // 设置自定义。方便拓展
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
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

// 启动一个Worker工作池（开启工作池的动作只能发生一次，一个zinx框架只能有一个worker工作池）
func (mh *MsgHandle) StartWorkerPool() {
	// 根据workerPoolSize 分别开启Worker，每个Worker用一个go来承载
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		// 一个worker被启动
		// 1 当前的worker对应的channel消息队列。开辟空间 第0个workerchannerl
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		// 2 启动当前的worker，堵塞等待消息从channel传递过来
		go mh.StartOneWorker(uint32(i), mh.TaskQueue[i])
	}
}

// 启动一个Worker工作流程
func (mh *MsgHandle) StartOneWorker(workerID uint32, taskQueue chan ziface.IRequest) {
	fmt.Println("WorkerID=", workerID, " is start...")

	// 不断堵塞等待对应消息队列的消息
	for {
		select {
		// 如果有消息过来，出列的就是一个客户端的Request，执行当Request所绑定业务
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

// 将消息交给TaskQueue，由worker进行处理
func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	// 1 将消息平均分配给不通过的worker
	// 根据客户端建立的ConnID来进行分配
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Printf("Add ConnID = %v, RequestID = %v, workID = %v\n",
		request.GetConnection().GetConnID(),
		request.GetMsgID(),
		workerID,
	)
	// 2 将消息发送给对应的worker的TaskQueue即可
	mh.TaskQueue[workerID] <- request
}
