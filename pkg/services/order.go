package service

import (
	"context"
	"net/http"

	"github.com/RohithER12/order-svc/pkg/client"
	"github.com/RohithER12/order-svc/pkg/db"
	"github.com/RohithER12/order-svc/pkg/models"
	"github.com/RohithER12/order-svc/pkg/pb"
	orderinterface "github.com/RohithER12/order-svc/pkg/repo/orderInterface"
)

type Server struct {
	H          db.Handler
	ProductSvc client.ProductServiceClient
	pb.UnimplementedOrderServiceServer
	Order orderinterface.Order
}

func (s *Server) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	product, err := s.ProductSvc.FindOne(req.ProductId)

	if err != nil {
		return &pb.CreateOrderResponse{Status: http.StatusBadRequest, Error: err.Error()}, nil
	} else if product.Status >= http.StatusNotFound {
		return &pb.CreateOrderResponse{Status: product.Status, Error: product.Error}, nil
	} else if product.Data.Stock < req.Quantity {
		return &pb.CreateOrderResponse{Status: http.StatusConflict, Error: "Out of stock"}, nil
	}

	order := models.Order{
		Price:     product.Data.Price,
		ProductId: product.Data.Id,
		UserId:    req.UserId,
		Quantity:  req.Quantity,
	}

	// s.H.DB.Create(&order)

	orderId, err := s.Order.CreateOrder(order)
	if err != nil {
		return &pb.CreateOrderResponse{
			Status: http.StatusBadRequest,
			Error:  "Oreder create failed",
		}, nil
	}

	res, err := s.ProductSvc.DecreaseStock(req.ProductId, orderId, order.Quantity)

	if err != nil {
		return &pb.CreateOrderResponse{Status: http.StatusBadRequest, Error: err.Error()}, nil
	} else if res.Status == http.StatusConflict {
		s.H.DB.Delete(&models.Order{}, order.Id)

		return &pb.CreateOrderResponse{Status: http.StatusConflict, Error: res.Error}, nil
	}

	return &pb.CreateOrderResponse{
		Status: http.StatusCreated,
		Id:     order.Id,
	}, nil
}

// Uncomment the mustEmbedUnimplementedOrderServiceServer method
// func (s *Server) mustEmbedUnimplementedOrderServiceServer() {
// 	// Empty implementation
// }
