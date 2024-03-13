package testgrp

type Status struct {
	Status string
	Limit  string
	Offset string
}

type QueryParams struct {
	Limit  string `query:"limit"`
	Offset string `query:"offset"`
}

func (qp QueryParams) Params() map[string]string {
	return map[string]string{
		"limit":  qp.Limit,
		"offset": qp.Offset,
	}
}
