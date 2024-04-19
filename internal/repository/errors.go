package repository

import "errors"

var (
	ErrOwnerExists  = errors.New("owner with this name already exists")
	ErrCarExists    = errors.New("car with this registration number already exists")
	ErrCarNotFound  = errors.New("car with this registration number not found")
	ErrCarsNotFound = errors.New("cars not found")
)
