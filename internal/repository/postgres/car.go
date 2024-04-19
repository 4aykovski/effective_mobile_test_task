package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/4aykovski/effective_mobile_test_task/internal/model"
	"github.com/4aykovski/effective_mobile_test_task/internal/repository"
	"github.com/4aykovski/effective_mobile_test_task/pkg/database/postgres"
	"github.com/lib/pq"
)

type CarRepository struct {
	postgres *postgres.Postgres
}

func NewCarRepository(postgres *postgres.Postgres) *CarRepository {
	return &CarRepository{
		postgres: postgres,
	}
}

func (r *CarRepository) InsertCar(ctx context.Context, car model.Car) error {
	stmt, err := r.postgres.Prepare(
		`INSERT INTO cars (registration_number, mark, model, year, owner_name, owner_surname) 
  			 	VALUES ($1, $2, $3, $4, $5, $6)`,
	)
	if err != nil {
		return fmt.Errorf("failed to prepare add new car statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, car.RegistrationNumber, car.Mark, car.Model, car.Year, car.OwnerName, car.OwnerSurname)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code.Name() == "unique_violation" {
				return repository.ErrCarExists
			}
		}

		return fmt.Errorf("failed to execute add new car statement: %w", err)
	}

	return nil
}

func (r *CarRepository) DeleteCar(ctx context.Context, regNumber string) error {
	stmt, err := r.postgres.Prepare("DELETE FROM cars WHERE registration_number = $1")
	if err != nil {
		return fmt.Errorf("failed to prepare delete car statement: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, regNumber)
	if err != nil {
		return fmt.Errorf("failed to execute delete car statement: %w", err)
	}

	deleted, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if deleted == 0 {
		return repository.ErrCarNotFound
	}

	return nil
}

func (r *CarRepository) UpdateCar(ctx context.Context, car model.Car) error {

	stmt, err := r.postgres.Prepare(
		`UPDATE cars
		SET mark = $1, model = $2, year = $3, owner_name = $4, owner_surname = $5
		WHERE registration_number = $6`,
	)
	if err != nil {
		return fmt.Errorf("failed to prepare update car statement: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, car.Mark, car.Model, car.Year, car.OwnerName, car.OwnerSurname, car.RegistrationNumber)
	if err != nil {
		return fmt.Errorf("failed to execute update car statement: %w", err)
	}

	updated, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if updated == 0 {
		return repository.ErrCarNotFound
	}

	return nil
}

func (r *CarRepository) GetCar(ctx context.Context, regNumber string) (model.Car, error) {

	stmt, err := r.postgres.Prepare("SELECT * FROM cars WHERE registration_number = $1")
	if err != nil {
		return model.Car{}, fmt.Errorf("failed to prepare get car statement: %w", err)
	}
	defer stmt.Close()

	var car model.Car
	err = stmt.QueryRowContext(ctx, regNumber).Scan(&car.RegistrationNumber, &car.Mark, &car.Model, &car.Year, &car.OwnerName, &car.OwnerSurname)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Car{}, repository.ErrCarNotFound
		}

		return model.Car{}, fmt.Errorf("failed to execute get car statement: %w", err)
	}

	return car, nil
}
