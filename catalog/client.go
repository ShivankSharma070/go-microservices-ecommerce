package catalog

import (
	"context"

	"github.com/ShivankSharma070/go-microservices-ecommerce/catalog/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.CatalogServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:    conn,
		service: pb.NewCatalogServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) PostProduct(ctx context.Context, name, description string, price float64) (*Product, error) {
	resp, err := c.service.PostProduct(
		ctx,
		&pb.PostProductRequest{
			Name:        name,
			Description: description,
			Price:       price,
		},
	)

	if err != nil {
		return nil, err
	}

	return &Product{
		ID:          resp.Product.Id,
		Name:        resp.Product.Name,
		Description: resp.Product.Description,
		Price:       resp.Product.Price,
	}, nil
}

func (c *Client) GetProduct(ctx context.Context, id string) (*Product, error) {
	resp, err := c.service.GetProduct(
		ctx,
		&pb.GetProductRequest{
			Id: id,
		},
	)

	if err != nil {
		return nil, err
	}

	return &Product{
		ID:          resp.Product.Id,
		Name:        resp.Product.Name,
		Description: resp.Product.Description,
		Price:       resp.Product.Price,
	}, nil
}

func (c *Client) GetProducts(ctx context.Context, query string, ids []string, skip, take uint64) ([]Product, error) {
	resp, err := c.service.GetProducts(
		ctx,
		&pb.GetProductsRequest{
			Query: query,
			Ids:   ids,
			Skip:  skip,
			Take:  take,
		},
	)
	if err != nil {
		return nil, err
	}

	var products []Product

	for _, product := range resp.Products {
		products = append(products, Product{
			ID:          product.Id,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
		})
	}

	return products, nil
}
