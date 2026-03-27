package orders

import (
	"context"
	"github.com/ShivankSharma070/go-microservices-ecommerce/orders/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type Client struct {
	service pb.OrderServiceClient
	conn    *grpc.ClientConn
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:    conn,
		service: pb.NewOrderServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) PostOrder(ctx context.Context, accountID string, products []OrderedProduct) (*Order, error) {
	protoProducts := []*pb.PostOrderRequest_OrderedProduct{}
	for _, p := range products {
		protoProducts = append(protoProducts, &pb.PostOrderRequest_OrderedProduct{
			ProductId: p.ID,
			Quantity:  p.Quantity,
		})
	}

	r, err := c.service.PostOrder(ctx, &pb.PostOrderRequest{AccountId: accountID, Products: protoProducts})
	if err != nil {
		return nil, err
	}
	newOrderCreateAt := time.Time{}
	newOrderCreateAt.UnmarshalBinary(r.Order.CreateAt)

	return &Order{
		ID:         r.Order.Id,
		CreatedAt:  newOrderCreateAt,
		AccountId:  accountID,
		TotalPrice: r.Order.TotalPrice,
		Products:   products,
	}, nil
}

func (c *Client) GetOrderForAccount(ctx context.Context, accountID string) ([]Order, error) {
	r, err := c.service.GetOrderForAccount(ctx, &pb.GetOrderForAccountRequest{
		AccountId: accountID,
	})

	if err != nil {
		return nil, err
	}
	orders := []Order{}

	for _, o := range r.Orders {
		newOrder := Order{
			ID:         o.Id,
			AccountId:  o.AccountId,
			TotalPrice: o.TotalPrice,
		}

		newOrder.CreatedAt = time.Time{}
		newOrder.CreatedAt.UnmarshalBinary(o.CreateAt)

		products := []OrderedProduct{}
		for _, p := range o.Products {
			products = append(products, OrderedProduct{
				ID:          p.Id,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
				Quantity:    p.Quantity,
			})
		}
		newOrder.Products = products
		orders = append(orders, newOrder)
	}

	return orders, nil
}
