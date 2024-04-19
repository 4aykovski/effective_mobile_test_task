package ownerservice

import (
	"context"
	"fmt"

	"github.com/4aykovski/effective_mobile_test_task/internal/model"
)

type ownerRepository interface {
	InsertOwner(ctx context.Context, owner model.Owner) error
}

type Service struct {
	ownerRepository ownerRepository
}

func New(ownerRepository ownerRepository) *Service {
	return &Service{
		ownerRepository: ownerRepository,
	}
}

type AddNewOwnerInput struct {
	Name       string
	Surname    string
	Patronymic string
}

func (s *Service) AddNewOwner(ctx context.Context, input AddNewOwnerInput) error {
	owner := model.Owner{
		Name:       input.Name,
		Surname:    input.Surname,
		Patronymic: input.Patronymic,
	}

	err := s.ownerRepository.InsertOwner(ctx, owner)
	if err != nil {
		return fmt.Errorf("can't add new owner: %w", err)
	}

	return nil
}
