package main

import (
	"context"
)

type accountResolver struct {
	server *Server
}

func (r *accountResolver) Order(ctx context.Context, obj *Account) ([]*Order, error) {
}
