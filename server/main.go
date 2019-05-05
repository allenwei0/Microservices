package main

import (
	"log"
	"net"

	// 导入生成的consignment.pb.go文件
	pb "../go_micro_srv_consignment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"golang.org/x/net/context"
)

const (
	port = ":50051"
)

type IRepository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
}

// Repository - 模拟一个数据库，我们会在此后使用真正的数据库替代他
type Repository struct {
	consignments []*pb.Consignment
}

func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	updated := append(repo.consignments, consignment)
	repo.consignments = updated
	return consignment, nil
}

// service要实现在proto中定义的所有方法。当你不确定时
// 可以去对应的*.pb.go文件里查看需要实现的方法及其定义
type service struct {
	repo IRepository
}

// CreateConsignment - 在proto中，我们只给这个微服务定一个了一个方法
// 就是这个CreateConsignment方法，它接受一个context以及proto中定义的
// Consignment消息，这个Consignment是由gRPC的服务器处理后提供给你的
func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment) (*pb.Response, error) {
	// 保存我们的consignment
	consignment, err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}
	// 返回的数据也要符合proto中定义的数据结构
	return &pb.Response{Created: true, Consignment: consignment}, nil
}
func main() {
	repo := &Repository{}
	// 设置gRPC服务器
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	// 在我们的gRPC服务器上注册微服务，这会将我们的代码和*.pb.go中
	// 的各种interface对应起来
	pb.RegisterShippingServiceServer(s, &service{repo})
	// 在gRPC服务器上注册reflection
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
