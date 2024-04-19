package carinfoservice

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/4aykovski/effective_mobile_test_task/pkg/client/carinfo"
)

type сarInfoClient interface {
	GetCarInfoByRegNumber(ctx context.Context, number string) ([]byte, error)
}

type Service struct {
	client сarInfoClient
}

func New(client сarInfoClient) *Service {
	return &Service{
		client: client,
	}
}

func (service *Service) GetCarInfoByRegNumber(ctx context.Context, regNumber string) (*carinfo.CarInfo, error) {
	res, err := service.client.GetCarInfoByRegNumber(ctx, regNumber)
	if err != nil {
		return nil, fmt.Errorf("can't get car info: %w", err)
	}

	var carInfo carinfo.CarInfo
	if err = json.Unmarshal(res, &carInfo); err != nil {
		return nil, fmt.Errorf("can't unmarshal car info: %w", err)
	}

	return &carInfo, nil
}
