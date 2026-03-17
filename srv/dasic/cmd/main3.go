package main

import (
	"flag"
	"fmt"
	"lianxi/srv/dasic/config"
	"lianxi/srv/dasic/inits"
	"lianxi/srv/handler/service/goods"
	"lianxi/srv/handler/service/order"
	"log"
	"net"

	__ "lianxi/srv/proto/goods"
	__2 "lianxi/srv/proto/order"

	"google.golang.org/grpc"
)

//var (
//	port = flag.Int("port", 50051, "The server port")
//)

func main() {

	if err := inits.ConsulInit(); err != nil {
		log.Fatalf("Consul初始化失败: %v", err)
	}
	log.Println("Consul初始化成功")
	services, err := inits.GetServiceWithLoadBalancer(config.Gen.Consul.ServiceName)
	if err != nil {
		log.Printf("获取用户服务失败: %v", err)
	} else {
		log.Printf("获取到用户服务: %s, 地址: %s:%d", services.Service, services.Address, services.Port)
	}
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 50052))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	__.RegisterGoodsServer(s, &goods.Server{})
	__2.RegisterOrderServer(s, &order.Server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	err = inits.ConsulShutdown()
	if err != nil {
		return
	}
	fmt.Println("服务已退出")
}
