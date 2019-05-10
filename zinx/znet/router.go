package znet

import "ZINX_PRO/zinx/ziface"

// 实现router时，先嵌入该 基类，然后根据需要对这个基类的方法进行重写
type BaseRouter struct{}

// 为什么BaseRouter 方法都为空
// 因为，有的 Router 不希望有的 PreHandle or PostHandle
// 所以，Router 全部继承 BaseRouter 好处：
// 不需要实现 PreHandle or PostHandle 也可以实例化
// 在处理conn业务之前的钩子方法
func (br *BaseRouter) PreHandle(request ziface.IRequest) {}

// 处理conn业务的方法
func (br *BaseRouter) Handle(request ziface.IRequest) {}

// 处理conn业务之后的钩子方法
func (br *BaseRouter) PostHandle(request ziface.IRequest) {}
