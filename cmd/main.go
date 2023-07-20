package main

import (
	"fmt"
	"log"
	"net"

	"github.com/RohithER12/order-svc/pkg/client"
	"github.com/RohithER12/order-svc/pkg/config"
	"github.com/RohithER12/order-svc/pkg/db"
	"github.com/RohithER12/order-svc/pkg/pb"
	"github.com/RohithER12/order-svc/pkg/repo"
	orderinterface "github.com/RohithER12/order-svc/pkg/repo/orderInterface"
	service "github.com/RohithER12/order-svc/pkg/services"
	"github.com/google/wire"
	"google.golang.org/grpc"
)

func main() {
	c, err := config.LoadConfig()

	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	h := db.Init(c.DBUrl)

	lis, err := net.Listen("tcp", c.Port)

	if err != nil {
		log.Fatalln("Failed to listing:", err)
	}

	productSvc := client.InitProductServiceClient(c.ProductSvcUrl)

	if err != nil {
		log.Fatalln("Failed to listing:", err)
	}

	fmt.Println("Order Svc on", c.Port)

	order := InitializeOrderImpl(&h)
	s := service.Server{
		H:          h,
		ProductSvc: productSvc,
		Order:      order,
	}

	grpcServer := grpc.NewServer()

	pb.RegisterOrderServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
	}
}

func InitializeOrderImpl(h *db.Handler) orderinterface.Order {
	wire.Build(orderinterface.NewOrderImpl)
	return &repo.OrderRepo{
		H: *h,
	}
}
