package ziface

/*
	路由接口，这里面路由是 使用框架者给该链接自定的处理业务员的方法
	路由里的 IRequest 到包含用该链接的链接信息和该链接请求的数据信息
*/
type IRouter interface {
	PreHandle(request IRequest)  // 在处理conn业务之前的钩子方法
	Handle(request IRequest)     // 处理conn业务的方法
	PostHandle(request IRequest) // 处理conn业务之后的钩子方法
}
