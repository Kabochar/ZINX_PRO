package znet

import (
	"ZINX_PRO/zinx/ziface"
	"fmt"
	"github.com/kataras/iris/core/errors"
	"sync"
)

/*
	连接管理模块
*/
type ConnManager struct {
	connections map[uint32]ziface.IConnection
	connLock    sync.RWMutex
}

/*
	创建一个连接管理
*/
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

// 添加链接
func (connMgr *ConnManager) Add(conn ziface.IConnection) {
	// 保护共享资源map 读写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	// 将conn连接添加到ConnManager中
	connMgr.connections[conn.GetConnID()] = conn

	fmt.Println("connection add to ConnManager successfully: conn num = ",
		connMgr.Len())
}

// 删除链接
func (connMgr *ConnManager) Remove(conn ziface.IConnection) {
	// 保护共享资源Map 加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	// 删除连接信息
	delete(connMgr.connections, conn.GetConnID())

	fmt.Println()
}

// 根据ConnID 获取链接
func (connMgr *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	// 保护共享资源Mao 加读锁
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not found")
	}
}

// 得到当前链接总数
func (connMgr *ConnManager) Len() int {
	return len(connMgr.connections)
}

// 清除并终止所有的链接
func (connMgr *ConnManager) ClearConn() {
	// 保护共享资源Map 加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	// 停止并删除全部的连接信息
	for connID, conn := range connMgr.connections {
		// 停止
		conn.Stop()
		// 删除
		delete(connMgr.connections, connID)
	}
}
