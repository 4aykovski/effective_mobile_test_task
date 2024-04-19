package carservice

import (
	"context"
	"fmt"

	"github.com/4aykovski/effective_mobile_test_task/internal/model"
)

type carRepository interface {
	CreateCar(ctx context.Context, car *model.Car) error
}

type Service struct {
	carRepository carRepository
}

func NewCarService(carRepository carRepository) *Service {
	return &Service{
		carRepository: carRepository,
	}
}

type AddNewCarInput struct {
	RegistrationNumber string
	Mark               string
	Model              string
	Year               int
	OwnerName          string
	OwnerSurname       string
}

func (service *Service) AddNewCar(ctx context.Context, car AddNewCarInput) error {

	carInfo := &model.Car{
		RegistrationNumber: car.RegistrationNumber,
		Mark:               car.Mark,
		Model:              car.Model,
		Year:               car.Year,
		OwnerName:          car.OwnerName,
		OwnerSurname:       car.OwnerSurname,
	}

	if err := service.carRepository.CreateCar(ctx, carInfo); err != nil {
		return fmt.Errorf("failed to create car: %w", err)
	}

	return nil
}