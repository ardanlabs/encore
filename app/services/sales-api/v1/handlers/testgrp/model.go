package testgrp

type Status struct {
	Status string
	Limit  int
	Offset int
}

type QueryParams struct {
	Limit  int `query:"limit"`
	Offset int `query:"offset"`
}
