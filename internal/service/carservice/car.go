package carservice

import (
	"context"
	"fmt"

	"github.com/4aykovski/effective_mobile_test_task/internal/model"
)

type carRepository interface {
	InsertCar(ctx context.Context, car model.Car) error
	DeleteCar(ctx context.Context, regNumber string) error
	UpdateCar(ctx context.Context, car model.Car) error
	GetCar(ctx context.Context, regNumber string) (model.Car, error)
	GetCars(ctx context.Context, limit, offset int) ([]model.Car, error)
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

func (s *Service) AddNewCar(ctx context.Context, car AddNewCarInput) error {

	carInfo := model.Car{
		RegistrationNumber: car.RegistrationNumber,
		Mark:               car.Mark,
		Model:              car.Model,
		Year:               car.Year,
		OwnerName:          car.OwnerName,
		OwnerSurname:       car.OwnerSurname,
	}

	if err := s.carRepository.InsertCar(ctx, carInfo); err != nil {
		return fmt.Errorf("failed to create car: %w", err)
	}

	return nil
}

func (s *Service) DeleteCar(ctx context.Context, regNumber string) error {
	err := s.carRepository.DeleteCar(ctx, regNumber)
	if err != nil {
		return fmt.Errorf("failed to delete car: %w", err)
	}

	return nil
}

type UpdateCarInput struct {
	RegistrationNumber string
	Mark               string
	Model              string
	Year               int
	OwnerName          string
	OwnerSurname       string
}

func (s *Service) UpdateCar(ctx context.Context, car UpdateCarInput) error {

	oldCar, err := s.carRepository.GetCar(ctx, car.RegistrationNumber)
	if err != nil {
		return fmt.Errorf("failed to update car: %w", err)
	}

	carInfo := model.Car{
		RegistrationNumber: car.RegistrationNumber,
		Mark:               car.Mark,
		Model:              car.Model,
		Year:               car.Year,
		OwnerName:          car.OwnerName,
		OwnerSurname:       car.OwnerSurname,
	}

	if car.Year == 0 {
		carInfo.Year = oldCar.Year
	}

	if car.Mark == "" {
		carInfo.Mark = oldCar.Mark
	}

	if car.Model == "" {
		carInfo.Model = oldCar.Model
	}

	if car.OwnerName == "" {
		carInfo.OwnerName = oldCar.OwnerName
	}

	if car.OwnerSurname == "" {
		carInfo.OwnerSurname = oldCar.OwnerSurname
	}

	if err = s.carRepository.UpdateCar(ctx, carInfo); err != nil {
		return fmt.Errorf("failed to update car: %w", err)
	}

	return nil
}

func (s *Service) GetCars(ctx context.Context, limit, offset int) ([]model.Car, error) {
	cars, err := s.carRepository.GetCars(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get cars: %w", err)
	}

	return cars, nil
}
