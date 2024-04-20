package handler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"reflect"
	"sync"

	"github.com/4aykovski/effective_mobile_test_task/internal/model"
	"github.com/4aykovski/effective_mobile_test_task/internal/repository"
	"github.com/4aykovski/effective_mobile_test_task/internal/service/carservice"
	"github.com/4aykovski/effective_mobile_test_task/internal/service/ownerservice"
	"github.com/4aykovski/effective_mobile_test_task/pkg/api/filter"
	"github.com/4aykovski/effective_mobile_test_task/pkg/client/carinfo"
	"github.com/4aykovski/effective_mobile_test_task/pkg/mapper"
	"github.com/4aykovski/effective_mobile_test_task/pkg/response"
	"github.com/4aykovski/effective_mobile_test_task/pkg/tag"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

const (
	requestWithWrongBody = "request with wrong body"
	invalidParameter     = "invalid parameter"
)

type carInfoService interface {
	GetCarInfoByRegNumber(ctx context.Context, regNumber []string) map[string]carinfo.CarInfo
}

type carService interface {
	AddNewCars(ctx context.Context, cars []carservice.AddNewCarInput, errs chan error) *sync.Map
	DeleteCar(ctx context.Context, regNumber string) error
	UpdateCar(ctx context.Context, car carservice.UpdateCarInput) error
	GetCars(ctx context.Context, limit, offset int, filterOptions filter.Options) ([]model.Car, error)
}

type ownerService interface {
	AddNewOwners(ctx context.Context, owners []ownerservice.AddNewOwnerInput, errs chan error)
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
	RegNumber []string `json:"regNumber"`
}

type AddNewCarResponse struct {
	response.Response
	ProcessedCars map[string]string `json:"processed_cars"`
}

func (h *CarHandler) AddNewCar(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		log = log.With(
			slog.String("handler", "AddNewCar"),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var input AddNewCarInput
		if err := render.DecodeJSON(r.Body, &input); err != nil {
			log.Info("request with wrong body")

			renderResponse(w, r, response.BadRequest(requestWithWrongBody), http.StatusBadRequest)
			return
		}
		log.Debug("input", slog.String("input", fmt.Sprint(input)))

		carInfos := h.carInfoService.GetCarInfoByRegNumber(r.Context(), input.RegNumber)
		if len(carInfos) == 0 {
			log.Info("can't find any car info")

			renderResponse(w, r, response.BadRequest(invalidParameter), http.StatusBadRequest)
			return
		}
		log.Debug("car infos", slog.Any("car_infos", carInfos))

		cars, owners := mapper.CarInfoIntoCarAndOwner(carInfos)
		log.Debug("cars", slog.Any("cars", cars))
		log.Debug("owners", slog.Any("owners", owners))

		errs := make(chan error, len(owners))
		h.ownerService.AddNewOwners(r.Context(), owners, errs)
		for err := range errs {
			log.Debug("failed to add new owner", slog.String("error", err.Error()))
			if err != nil && !errors.Is(err, repository.ErrOwnerExists) {
				log.Error("failed to add new owner", slog.String("error", err.Error()))

				renderResponse(w, r, response.InternalError(), http.StatusInternalServerError)
				return
			}
		}

		errs = make(chan error, len(cars))
		validNumbers := h.carService.AddNewCars(r.Context(), cars, errs)
		for err := range errs {
			log.Debug("failed to add new car", slog.String("error", err.Error()))
			if err != nil && !errors.Is(err, repository.ErrCarExists) {
				log.Error("failed to add new car", slog.String("error", err.Error()))

				renderResponse(w, r, response.InternalError(), http.StatusInternalServerError)
				return
			}
		}

		processedCars := make(map[string]string)
		validNumbers.Range(func(key, value any) bool {
			regNumber := key.(string)
			status := value.(string)
			processedCars[regNumber] = status
			return true
		})
		log.Debug("processed cars", slog.Any("processed_cars", processedCars))

		log.Info("cars processed", slog.Any("cars", processedCars))

		render.Status(r, http.StatusOK)
		render.JSON(w, r, AddNewCarResponse{Response: response.OK(), ProcessedCars: processedCars})
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
		log.Debug("reg number", slog.String("reg_number", regNumber))

		err := h.carService.DeleteCar(r.Context(), regNumber)
		if err != nil {
			if errors.Is(err, repository.ErrCarNotFound) {
				log.Info("can't find car with this registration number", slog.String("reg_number", regNumber))

				renderResponse(w, r, response.BadRequest(fmt.Sprintf("%s - regNumber", invalidParameter)), http.StatusBadRequest)
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
		log.Debug("reg number", slog.String("reg_number", regNumber))

		var input UpdateCarInput
		if err := render.DecodeJSON(r.Body, &input); err != nil {
			log.Info("request with wrong body")

			renderResponse(w, r, response.BadRequest(requestWithWrongBody), http.StatusBadRequest)
			return
		}
		log.Debug("input", slog.Any("input", input))

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

				renderResponse(w, r, response.BadRequest(fmt.Sprintf("%s - regNumber", invalidParameter)), http.StatusBadRequest)
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

		allowedFilters := make(map[string]string)
		carType := reflect.TypeOf(model.Car{})
		for i := 0; i < carType.NumField(); i++ {
			allowedFilter := tag.ParseJsonTag(carType.Field(i).Tag.Get("json"))
			type_ := carType.Field(i).Type.Name()
			allowedFilters[allowedFilter] = type_
		}
		log.Debug("allowed filters", slog.Any("filters", allowedFilters))

		filterOptions, err := getFiltersFromUrlQuery(r, allowedFilters)
		if err != nil {
			log.Info("invalid filter", slog.String("filter", err.Error()))

			renderResponse(w, r, response.BadRequest(fmt.Sprintf("%s - filter", invalidParameter)), http.StatusBadRequest)
			return
		}
		log.Debug("filter options", slog.Any("filters", filterOptions))

		limit, err := getLimitFromUrlQuery(r)
		if err != nil {
			log.Info("invalid limit", slog.String("limit", string(rune(limit))))

			renderResponse(w, r, response.BadRequest(fmt.Sprintf("%s - limit", invalidParameter)), http.StatusBadRequest)
			return
		}
		log.Debug("limit", slog.Int("limit", limit))

		offset, err := getOffsetFromUrlQuery(r)
		if err != nil {
			log.Info("invalid offset", slog.String("offset", string(rune(offset))))

			renderResponse(w, r, response.BadRequest(fmt.Sprintf("%s - offset", invalidParameter)), http.StatusBadRequest)
			return
		}
		log.Debug("offset", slog.Int("offset", offset))

		cars, err := h.carService.GetCars(r.Context(), limit, offset, filterOptions)
		if err != nil {
			if errors.Is(err, repository.ErrCarsNotFound) {
				log.Info("cars not found")

				render.Status(r, http.StatusOK)
				render.JSON(w, r, GetCarsResponse{
					Cars:     nil,
					Response: response.OK(),
				})
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
