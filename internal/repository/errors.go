package repository

import "errors"

var (
	ErrOwnerExists = errors.New("owner with this name already exists")
	ErrCarExists   = errors.New("car with this registration number already exists")
)
