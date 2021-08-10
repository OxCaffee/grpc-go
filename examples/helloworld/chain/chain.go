package chain

import (
	"context"
	"google.golang.org/grpc"
)

// UnaryServerInterceptorChain unary server interceptor chain
func UnaryServerInterceptorChain(interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	// 需要添加的拦截器个数，我们需要将这些拦截器组织成一个拦截链的形式
	length := len(interceptors)

	if length == 0 {
		// 说明此时没有拦截器，默认返回一个空的拦截器
		return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (
			resp interface{}, err error) {
			return handler(ctx, req)
		}
	}

	if length == 1 {
		// 拦截器个数为1
		return interceptors[0]
	}

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (
		resp interface{}, err error) {
		// 构造一个拦截器链
		chain := func(currentInterceptor grpc.UnaryServerInterceptor, currenHandler grpc.UnaryHandler) grpc.UnaryHandler {
			return func(ctx context.Context, req interface{}) (interface{}, error) {
				return currentInterceptor(ctx, req, info, handler)
			}
		}

		// 从当前的Handler开始调用
		chainHandler := handler
		for i := length - 1; i >= 0; i-- {
			chainHandler = chain(interceptors[i], handler)
		}
		return chainHandler(ctx, req)
	}
}
