package main

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/ShivankSharma070/go-microservices-ecommerce/orders"
)

var (
	ErrInvalidParameter= errors.New("Invalid parameter")
)

type mutationResolver struct {
	server *Server
}

func (r *mutationResolver) CreateAccount(ctx context.Context, in *AccountInput) (*Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	acc, err := r.server.accountClient.PostAccount(ctx, in.Name)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &Account{
		Id: acc.ID,
		Name: acc.Name,
	}, nil
}

func (r *mutationResolver) CreateProduct(ctx context.Context, in *ProductInput) (*Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	product, err := r.server.catalogClient.PostProduct(ctx,in.Name, in.Description, in.Price)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &Product{
		ID: product.ID,
		Name: product.Name,
		Description: product.Description,
		Price: product.Price,
	}, nil
}

func (r *mutationResolver) CreateOrder(ctx context.Context, in *OrderInput) (*Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	var products []orders.OrderedProduct 
	for _, p := range in.Products {
		if p.Quantity <= 0 {
			return nil, ErrInvalidParameter
		}
		products = append(products, orders.OrderedProduct{
			ID : p.ID,
			Quantity: uint32(p.Quantity),
		})
	}

	o, err:= r.server.orderClient.PostOrder(ctx, in.AccountID, products)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &Order{
		ID: o.ID,
		CreatedAt: o.CreatedAt,
		TotalPrice: o.TotalPrice,
	}, nil
}
