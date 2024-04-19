package handler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/4aykovski/effective_mobile_test_task/internal/model"
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
	UpdateCar(ctx context.Context, car carservice.UpdateCarInput) error
	GetCars(ctx context.Context, limit, offset int) ([]model.Car, error)
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

type UpdateCarInput struct {
	Mark         string `json:"mark,omitempty"`
	Model        string `json:"model,omitempty"`
	Year         int    `json:"year,omitempty"`
	OwnerName    string `json:"ownerName,omitempty"`
	OwnerSurname string `json:"ownerSurname,omitempty"`
}

func (h *CarHandler) UpdateCar(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		log = log.With(
			slog.String("handler", "UpdateCar"),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		regNumber := chi.URLParam(r, "reg_number")

		var input UpdateCarInput
		if err := render.DecodeJSON(r.Body, &input); err != nil {
			log.Info("request with empty body")

			renderResponse(w, r, response.BadRequest(), http.StatusBadRequest)
			return
		}

		err := h.carService.UpdateCar(r.Context(), carservice.UpdateCarInput{
			RegistrationNumber: regNumber,
			Mark:               input.Mark,
			Model:              input.Model,
			Year:               input.Year,
			OwnerName:          input.OwnerName,
			OwnerSurname:       input.OwnerSurname,
		})
		if err != nil {
			if errors.Is(err, repository.ErrCarNotFound) {
				log.Info("can't find car with this registration number", slog.String("reg_number", regNumber))

				renderResponse(w, r, response.BadRequest(), http.StatusBadRequest)
				return
			}

			log.Error("failed to update car", slog.String("error", err.Error()))

			renderResponse(w, r, response.InternalError(), http.StatusInternalServerError)
			return
		}

		log.Info("car updated", slog.String("reg_number", regNumber))

		renderResponse(w, r, response.OK(), http.StatusOK)
		return
	}
}

type GetCarsResponse struct {
	Cars []model.Car `json:"cars"`
	response.Response
}

func (h *CarHandler) GetCars(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		log = log.With(
			slog.String("handler", "GetCars"),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		limit, err := getLimitFromUrlQuery(r)
		if err != nil {
			log.Info("invalid limit", slog.String("limit", string(rune(limit))))

			renderResponse(w, r, response.BadRequest(), http.StatusBadRequest)
			return
		}

		offset, err := getOffsetFromUrlQuery(r)
		if err != nil {
			log.Info("invalid offset", slog.String("offset", string(rune(offset))))

			renderResponse(w, r, response.BadRequest(), http.StatusBadRequest)
			return
		}

		cars, err := h.carService.GetCars(r.Context(), limit, offset)
		if err != nil {
			if errors.Is(err, repository.ErrCarsNotFound) {
				log.Info("cars not found")

				renderResponse(w, r, response.BadRequest(), http.StatusBadRequest)
				return
			}

			log.Error("failed to get cars", slog.String("error", err.Error()))

			renderResponse(w, r, response.InternalError(), http.StatusInternalServerError)
			return
		}

		log.Info("cars found", slog.Int("cars_count", len(cars)))

		render.Status(r, http.StatusOK)
		render.JSON(w, r, GetCarsResponse{
			Cars:     cars,
			Response: response.OK(),
		})
		return
	}
}
