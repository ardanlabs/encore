package testgrp

import "context"

type Handlers struct{}

func New() *Handlers {
	return &Handlers{}
}

func (h *Handlers) Test(ctx context.Context, qp *QueryParams) (*Status, error) {
	status := Status{
		Status: "OK",
		Limit:  qp.Limit,
		Offset: qp.Offset,
	}

	return &status, nil
}
