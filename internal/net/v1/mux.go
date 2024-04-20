package v1

import (
	"context"
	"log/slog"
	"sync"

	"github.com/4aykovski/effective_mobile_test_task/internal/model"
	"github.com/4aykovski/effective_mobile_test_task/internal/net/v1/handler"
	"github.com/4aykovski/effective_mobile_test_task/internal/net/v1/middleware"
	"github.com/4aykovski/effective_mobile_test_task/internal/service/carservice"
	"github.com/4aykovski/effective_mobile_test_task/internal/service/ownerservice"
	"github.com/4aykovski/effective_mobile_test_task/pkg/api/filter"
	"github.com/4aykovski/effective_mobile_test_task/pkg/client/carinfo"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

type carService interface {
	AddNewCars(ctx context.Context, cars []carservice.AddNewCarInput, errs chan error) *sync.Map
	DeleteCar(ctx context.Context, regNumber string) error
	UpdateCar(ctx context.Context, car carservice.UpdateCarInput) error
	GetCars(ctx context.Context, limit, offset int, filterOptions filter.Options) ([]model.Car, error)
}

type ownerService interface {
	AddNewOwners(ctx context.Context, owners []ownerservice.AddNewOwnerInput, errs chan error)
}

type carInfoService interface {
	GetCarInfoByRegNumber(ctx context.Context, regNumbers []string, errs chan error) map[string]carinfo.CarInfo
}

func NewMux(
	log *slog.Logger,
	carService carService,
	ownerService ownerService,
	carInfoService carInfoService,
) *chi.Mux {
	var (
		carHandler = handler.NewCarHandler(carInfoService, carService, ownerService)
		mux        = chi.NewMux()
	)

	mux.Use(chiMiddleware.RequestID)
	mux.Use(middleware.Logger(log))

	mux.Route("/api/v1", func(r chi.Router) {
		r.Route("/cars", func(r chi.Router) {
			r.Post("/", carHandler.AddNewCar(log))
			r.Delete("/{reg_number}", carHandler.DeleteCar(log))
			r.Put("/{reg_number}", carHandler.UpdateCar(log))
			r.Get("/", carHandler.GetCars(log))
		})
	})

	return mux
}
