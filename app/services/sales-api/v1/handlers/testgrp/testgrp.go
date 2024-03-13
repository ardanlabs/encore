package testgrp

import (
	"context"
	"errors"
	"math/rand"
	"net/http"

	v1 "github.com/ardanlabs/encore/business/web/v1"
)

type Handlers struct{}

func New() *Handlers {
	return &Handlers{}
}

func (h *Handlers) Test(ctx context.Context, qp *QueryParams) (*Status, error) {
	if n := rand.Intn(100); n%2 == 0 {
		return nil, v1.NewTrustedError(errors.New("trusted error"), http.StatusBadRequest)
	}

	status := Status{
		Status: "OK",
		Limit:  qp.Limit,
		Offset: qp.Offset,
	}

	return &status, nil
}

func (h *Handlers) TestAuth(ctx context.Context, qp *QueryParams) (*Status, error) {
	status := Status{
		Status: "OK",
		Limit:  qp.Limit,
		Offset: qp.Offset,
	}

	return &status, nil
}
