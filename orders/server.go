package orders

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/ShivankSharma070/go-microservices-ecommerce/account"
	"github.com/ShivankSharma070/go-microservices-ecommerce/catalog"
	"github.com/ShivankSharma070/go-microservices-ecommerce/orders/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	service       Service
	accountClient *account.Client
	catalogClient *catalog.Client
	pb.UnimplementedOrderServiceServer
}

func ListenGRPC(s Service, accountURL, catalogURL string, port int) (err error) {
	accountClient, err := account.NewClient(accountURL)
	if err != nil {
		return err
	}

	catalogClient, err := catalog.NewClient(catalogURL)
	if err != nil {
		accountClient.Close()
		return
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		accountClient.Close()
		catalogClient.Close()
		return
	}

	serv := grpc.NewServer()
	pb.RegisterOrderServiceServer(serv, &grpcServer{
		service:       s,
		accountClient: accountClient,
		catalogClient: catalogClient,
	})

	reflection.Register(serv)
	return serv.Serve(lis)
}

func (s *grpcServer) PostOrder(ctx context.Context, r *pb.PostOrderRequest) (*pb.PostOrderResponse, error) {
	_, err := s.accountClient.GetAccount(ctx, r.AccountId)
	if err != nil {
		log.Println("Error getting account: ", err)
		return nil, fmt.Errorf("Account not found")
	}

	productIds := []string{}
	for _, p := range r.Products {
		productIds = append(productIds, p.ProductId)
	}

	orderedProducts, err := s.catalogClient.GetProducts(ctx, "", productIds, 0, 0)
	if err != nil {
		log.Println("Error getting products")
		return nil, fmt.Errorf("Products not found")
	}

	products := []OrderedProduct{}
	for _, p := range orderedProducts {
		product := OrderedProduct{
			ID:          p.ID,
			Description: p.Description,
			Name:        p.Name,
			Price:       p.Price,
			Quantity:    0,
		}

		for _, rp := range r.Products {
			if rp.ProductId == p.ID {
				product.Quantity = rp.Quantity
				break
			}
		}

		if product.Quantity > 0 {
			products = append(products, product)
		}
	}

	order, err := s.service.PostOrder(ctx, r.AccountId, products)
	if err != nil {
		log.Println("Error posting error:", err)
		return nil, errors.New("could not post order")
	}

	orderProto := &pb.Order{
		Id:         order.ID,
		AccountId:  order.AccountId,
		TotalPrice: order.TotalPrice,
		Products:   []*pb.Order_OrderedProduct{},
	}

	orderProto.CreateAt, _ = order.CreatedAt.MarshalBinary()
	for _, p := range order.Products {
		orderProto.Products = append(orderProto.Products, &pb.Order_OrderedProduct{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Quantity:    p.Quantity,
		})
	}

	return &pb.PostOrderResponse{
		Order: orderProto,
	}, nil
}

func (s *grpcServer) GetOrderForAccount(ctx context.Context, r *pb.GetOrderForAccountRequest) (*pb.GetOrderForAccountResponse, error) {
	accountOrder, err := s.service.GetOrderForAccount(ctx, r.AccountId)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	productIDMap := map[string]bool{}
	for _, o := range accountOrder {
		for _, p := range o.Products {
			productIDMap[p.ID] = true
		}
	}
	productIDs := []string{}
	for id := range productIDMap {
		productIDs = append(productIDs, id)
	}

	products, err := s.catalogClient.GetProducts(ctx, "", productIDs, 0, 0)
	if err != nil {
		log.Println("Error getting account products: ", err)
		return nil, err
	}

	orders := []*pb.Order{}
	for _, o := range accountOrder {
		op := &pb.Order{
			Id:         o.ID,
			AccountId:  o.AccountId,
			TotalPrice: o.TotalPrice,
			Products:   []*pb.Order_OrderedProduct{},
		}

		op.CreateAt, _ = o.CreatedAt.MarshalBinary()
		for _, product := range o.Products {
			for _, p := range products {
				if product.ID == p.ID {
					product.Name = p.Name
					product.Description = p.Description
					product.Price = p.Price
					break
				}
			}

			op.Products = append(op.Products, &pb.Order_OrderedProduct{
				Id:          product.ID,
				Name:        product.Name,
				Description: product.Description,
				Quantity:    product.Quantity,
				Price:       product.Price,
			})
		}

		orders = append(orders, op)
	}

	return &pb.GetOrderForAccountResponse{
		Orders: orders,
	}, nil
}
