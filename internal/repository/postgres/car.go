package postgres

import (
	"context"
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
