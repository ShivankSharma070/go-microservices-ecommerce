package main

import (
	"context"
	"log"
	"time"
)

type queryResolver struct {
	server *Server
}

func (r *queryResolver) Account(ctx context.Context, pagination *PaginationInput, id *string) ([]*Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if id != nil {
		resp, err := r.server.accountClient.GetAccount(ctx, *id)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		return []*Account{{
			Id:   resp.ID,
			Name: resp.Name,
		}}, nil
	}

	skip, take := uint64(0), uint64(0)
	if pagination != nil {
		skip, take = pagination.bounds()
	}

	accountList, err := r.server.accountClient.GetAccounts(ctx, skip, take)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	accounts := []*Account{}
	for _, ac := range accountList {
		accounts = append(accounts, &Account{
			Id:   ac.ID,
			Name: ac.Name,
		})
	}

	return accounts, nil
}

func (r *queryResolver) Product(ctx context.Context, pagination *PaginationInput, query *string, id *string) ([]*Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if id != nil {
		resp, err := r.server.catalogClient.GetProduct(ctx, *id)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		return []*Product{{
			ID:          resp.ID,
			Description: resp.Description,
			Name:        resp.Name,
			Price:       resp.Price,
		}}, nil
	}

	skip, take := uint64(0), uint64(0)
	if pagination != nil {
		skip, take = pagination.bounds()
	}

	q := ""
	if query != nil {
		q = *query
	}
	// NOTE: handle ids if bug found
	productList, err := r.server.catalogClient.GetProducts(ctx, q, nil, skip, take)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	products := []*Product{}
	for _, p := range productList {
		products = append(products, &Product{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		})
	}

	return products, nil
}

func (p *PaginationInput) bounds() (uint64, uint64) {
	skip := uint64(0)
	take := uint64(0)
	if p.Skip != nil {
		skip = uint64(*p.Skip)
	}
	if p.Take != nil {
		take = uint64(*p.Take)
	}
	return skip, take
}
