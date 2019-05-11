package utils

import (
	"ZINX_PRO/zinx/ziface"
	"encoding/json"
	"io/ioutil"
)

/*
	存储一切与Zinx框架的全局参数，供其它模块使用
	一些参数也可以通过 用户根据 zinx.json 来配置
*/
type GlobalObj struct {
	/*
		Server
	*/
	TcpServer ziface.IServer // 当前Zinx的全局Server对象
	Host      string         // 当前服务器主机 IP
	TcpPort   int            // 当前服务器主机监听端口号
	Name      string         // 当前服务器名称

	/*
	   Zinx
	*/
	Version          string // 当前Zinx版本号
	MaxPackageSize   uint32 // 传输数据包的最大Size
	MaxConn          int    // 当前服务器主机允许的最大链接个数
	WorkerPoolSize   uint32 // 当前业务工作Worker池的Goroutine数量
	MaxWorkerTaskLen uint32 // Zinx 框架允许用户最多开辟多少个Worker
}

/*
	定义一个全局 GlobalObject 对象
	目的就是让其他模块都能访问到 Zinx 里面的参数。
*/
var GlobalObject *GlobalObj

func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}
	// 将json数据解析到struct中
	// fmt.Printf("json	:%s\n",	data)
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

/*
	提供 init方法，默认加载
	初始化当前的GlobalObject
*/
func init() {
	GlobalObject = &GlobalObj{
		Name:             "ZinxServerApp",
		Version:          "V0.6",
		TcpPort:          8999,
		Host:             "0.0.0.0",
		MaxConn:          15,
		MaxPackageSize:   4096,
		WorkerPoolSize:   10,   // worker工作池的队列的个数
		MaxWorkerTaskLen: 1024, // 每个worker对应的消息队列的数量最大值
	}

	// 从配置文件中加载用户配置参数
	GlobalObject.Reload()
}
