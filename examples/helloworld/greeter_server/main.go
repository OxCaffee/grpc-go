/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a server for Greeter service.
package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

const (
	port = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

// InterceptChain 创建服务器端的拦截器链
func InterceptChain(interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	// 获取拦截器链的长度
	l := len(interceptors)
	// 我们返回一个拦截器
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {
		// 构造一个拦截器链
		chain := func(currentInterceptor grpc.UnaryServerInterceptor,
			currentHandler grpc.UnaryHandler) grpc.UnaryHandler {
			return func(ctx context.Context, req interface{}) (interface{}, error) {
				return currentInterceptor(ctx, req, info, currentHandler)
			}
		}

		// 声明一个Handler
		chainHandler := handler
		for i := l - 1; i >= 0; i-- {
			// 递归调用
			chainHandler = chain(interceptors[i], chainHandler)
		}
		return chainHandler(ctx, req)
	}
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// !!! 在这里开始模拟拦截器interceptor
	var myInterceptor grpc.UnaryServerInterceptor
	myInterceptor = func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (
		resp interface{}, err error) {
		// 拦截之前我们需要做什么
		fmt.Println("这里是server端的服务器配置")
		// 拦截后置处理
		return handler(ctx, req)
	}

	// <1> 创建server
	s := grpc.NewServer(
		[]grpc.ServerOption{grpc.UnaryInterceptor(myInterceptor)}...)
	// <2> 将服务注册到这个grpc server当中
	// s 代表的是 grpc server，后面的要注册的server是服务的实例
	pb.RegisterGreeterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
