package postgres

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/4aykovski/effective_mobile_test_task/internal/model"
	"github.com/4aykovski/effective_mobile_test_task/internal/repository"
	"github.com/4aykovski/effective_mobile_test_task/pkg/database/postgres"
	"github.com/lib/pq"
)

type OwnerRepository struct {
	postgres *postgres.Postgres
}

func NewOwnerRepository(postgres *postgres.Postgres) *OwnerRepository {
	return &OwnerRepository{
		postgres: postgres,
	}
}

func (r *OwnerRepository) InsertOwner(ctx context.Context, owner model.Owner) error {
	stmt, err := r.postgres.Prepare(
		`INSERT INTO owners (name, surname, patronymic)
  			 	VALUES ($1, $2, $3)`,
	)
	if err != nil {
		return fmt.Errorf("failed to prepare add new owner statement: %w", err)
	}
	defer stmt.Close()

	var mu sync.Mutex
	mu.Lock()
	_, err = stmt.ExecContext(ctx, owner.Name, owner.Surname, owner.Patronymic)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code.Name() == "unique_violation" {
				return repository.ErrOwnerExists
			}
		}

		return fmt.Errorf("failed to execute add new owner statement: %w", err)
	}
	mu.Unlock()

	return nil
}
