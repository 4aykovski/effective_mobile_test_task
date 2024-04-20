package carinfoservice

import (
	"context"
	"encoding/json"
	"sync"

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

func (service *Service) GetCarInfoByRegNumber(ctx context.Context, regNumbers []string, errs chan error) map[string]carinfo.CarInfo {
	carInfos := make(map[string]carinfo.CarInfo)
	var wg sync.WaitGroup
	wg.Add(len(regNumbers))

	for _, regNumber := range regNumbers {
		go func(regNumber string) {
			defer wg.Done()
			res, err := service.client.GetCarInfoByRegNumber(ctx, regNumber)
			if err != nil {
				errs <- err
				carInfos[regNumber] = carinfo.CarInfo{}
				return
			}
			var carInfo carinfo.CarInfo
			if err = json.Unmarshal(res, &carInfo); err != nil {
				errs <- err
				carInfos[regNumber] = carinfo.CarInfo{}
				return
			}
			carInfos[regNumber] = carInfo
		}(regNumber)
	}

	wg.Wait()
	close(errs)

	return carInfos
}
