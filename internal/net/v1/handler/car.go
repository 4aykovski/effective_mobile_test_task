package handler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/4aykovski/effective_mobile_test_task/internal/repository"
	"github.com/4aykovski/effective_mobile_test_task/internal/service/carservice"
	"github.com/4aykovski/effective_mobile_test_task/internal/service/ownerservice"
	"github.com/4aykovski/effective_mobile_test_task/pkg/client"
	"github.com/4aykovski/effective_mobile_test_task/pkg/client/carinfo"
	"github.com/4aykovski/effective_mobile_test_task/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type carInfoService interface {
	GetCarInfoByRegNumber(ctx context.Context, regNumber string) (*carinfo.CarInfo, error)
}

type carService interface {
	AddNewCar(ctx context.Context, car carservice.AddNewCarInput) error
	DeleteCar(ctx context.Context, regNumber string) error
}

type ownerService interface {
	AddNewOwner(ctx context.Context, owner ownerservice.AddNewOwnerInput) error
}

type CarHandler struct {
	carInfoService carInfoService
	carService     carService
	ownerService   ownerService
}

func NewCarHandler(carInfoService carInfoService, carService carService, ownerService ownerService) *CarHandler {
	return &CarHandler{
		carInfoService: carInfoService,
		carService:     carService,
		ownerService:   ownerService,
	}
}

type AddNewCarInput struct {
	RegNumber string `json:"regNumber"`
}

func (h *CarHandler) AddNewCar(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		log = log.With(
			slog.String("handler", "AddNewCar"),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var input AddNewCarInput
		if err := render.DecodeJSON(r.Body, &input); err != nil {
			log.Info("request with empty body")

			renderResponse(w, r, response.BadRequest(), http.StatusBadRequest)
			return
		}

		carInfo, err := h.carInfoService.GetCarInfoByRegNumber(r.Context(), input.RegNumber)
		if err != nil {
			if errors.Is(err, client.Err400StatusCode) {
				log.Info("can't find car with this registration number", slog.String("reg_number", input.RegNumber))

				renderResponse(w, r, response.BadRequest(), http.StatusBadRequest)
				return
			}

			log.Error("Failed to get car info", slog.String("error", err.Error()))

			renderResponse(w, r, response.InternalError(), http.StatusInternalServerError)
			return
		}

		car := carservice.AddNewCarInput{
			RegistrationNumber: carInfo.RegNumber,
			Mark:               carInfo.Mark,
			Model:              carInfo.Model,
			Year:               carInfo.Year,
			OwnerName:          carInfo.Owner.Name,
			OwnerSurname:       carInfo.Owner.Surname,
		}

		owner := ownerservice.AddNewOwnerInput{
			Name:       carInfo.Owner.Name,
			Surname:    carInfo.Owner.Surname,
			Patronymic: carInfo.Owner.Patronymic,
		}

		err = h.ownerService.AddNewOwner(r.Context(), owner)
		if err != nil && !errors.Is(err, repository.ErrOwnerExists) {
			log.Error("failed to add new owner", slog.String("error", err.Error()))

			renderResponse(w, r, response.InternalError(), http.StatusInternalServerError)
			return
		}

		err = h.carService.AddNewCar(r.Context(), car)
		if err != nil {
			if errors.Is(err, repository.ErrCarExists) {
				log.Info("car with this registration number already exists", slog.String("reg_number", car.RegistrationNumber))

				renderResponse(w, r, response.BadRequest(), http.StatusBadRequest)
				return
			}

			log.Error("failed to add new car", slog.String("error", err.Error()))

			renderResponse(w, r, response.InternalError(), http.StatusInternalServerError)
			return
		}

		renderResponse(w, r, response.OK(), http.StatusOK)
		return
	}
}

func (h *CarHandler) DeleteCar(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		log = log.With(
			slog.String("handler", "DeleteCar"),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		regNumber := chi.URLParam(r, "reg_number")

		err := h.carService.DeleteCar(r.Context(), regNumber)
		if err != nil {
			if errors.Is(err, repository.ErrCarNotFound) {
				log.Info("can't find car with this registration number", slog.String("reg_number", regNumber))

				renderResponse(w, r, response.BadRequest(), http.StatusBadRequest)
				return
			}

			log.Error("failed to delete car", slog.String("error", err.Error()))

			renderResponse(w, r, response.InternalError(), http.StatusInternalServerError)
			return
		}

		log.Info("car deleted", slog.String("reg_number", regNumber))

		renderResponse(w, r, response.OK(), http.StatusOK)
		return
	}
}
