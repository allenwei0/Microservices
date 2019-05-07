package main

import (
	"fmt"

	pb "../go_micro_srv_consignment"

	// 使用 go-mircro
	micro "github.com/micro/go-micro"
	"golang.org/x/net/context"
)

type IRepository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
	GetAll() []*pb.Consignment
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
func (repo *Repository) GetAll() []*pb.Consignment {
	return repo.consignments
}

// service要实现在proto中定义的所有方法。当你不确定时
// 可以去对应的*.pb.go文件里查看需要实现的方法及其定义
type service struct {
	repo IRepository
}

// CreateConsignment - 在proto中，我们只给这个微服务定一个了一个方法
// 就是这个CreateConsignment方法，它接受一个context以及proto中定义的
// Consignment消息，这个Consignment是由gRPC的服务器处理后提供给你的
func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {
	consignment, err := s.repo.Create(req)
	if err != nil {
		return err
	}
	res.Created = true
	res.Consignment = consignment
	return nil
}
func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
	consignments := s.repo.GetAll()
	res.Consignments = consignments
	return nil
}
func main() {
	repo := &Repository{}
	// 注意，在这里我们使用go-micro的NewService方法来创建新的微服务服务器，
	// 而不是上一篇文章中所用的标准
	srv := micro.NewService(
		// This name must match the package name given in your protobuf definition
		// 注意，Name方法的必须是你在proto文件中定义的package名字
		micro.Name("go.micro.srv.consignment"),
		micro.Version("latest"),
	)
	// Init方法会解析命令行flags
	srv.Init()
	pb.RegisterShippingServiceHandler(srv.Server(), &service{repo})
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
