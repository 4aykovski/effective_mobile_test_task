package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/4aykovski/effective_mobile_test_task/pkg/api/filter"
	"github.com/4aykovski/effective_mobile_test_task/pkg/response"
	"github.com/go-chi/render"
)

func renderResponse(w http.ResponseWriter, r *http.Request, resp response.Response, statusCode int) {
	render.Status(r, statusCode)
	render.JSON(w, r, resp)
}

func getLimitFromUrlQuery(r *http.Request) (int, error) {
	strLimit := r.URL.Query().Get("limit")
	limit := -1
	var err error
	if strLimit != "" {
		limit, err = strconv.Atoi(strLimit)
		if err != nil || limit < -1 {
			return -1, fmt.Errorf("failed to parse limit: %w", err)
		}
	}
	return limit, nil
}

func getOffsetFromUrlQuery(r *http.Request) (int, error) {
	strOffset := r.URL.Query().Get("offset")
	offset := 0
	var err error
	if strOffset != "" {
		offset, err = strconv.Atoi(strOffset)
		if err != nil || offset < -1 {
			return -1, fmt.Errorf("failed to parse offset: %w", err)
		}
	}
	return offset, nil
}

func getFiltersFromUrlQuery(r *http.Request, allowedFilters map[string]string) (filter.Options, error) {
	filterOptions := filter.NewOptions()
	for filterName, filterType := range allowedFilters {
		strValue := r.URL.Query().Get(filterName)
		if strValue != "" {
			var type_ string
			var operator string

			switch filterType {
			case filter.DataTypeInt:
				type_ = filter.DataTypeInt
				operator = filter.OperatorEq

				if strings.Index(strValue, ":") != -1 {
					split := strings.Split(strValue, ":")
					operator = split[0]
					strValue = split[1]

					if _, err := strconv.Atoi(strValue); err != nil {
						return nil, fmt.Errorf("failed to parse filter: %w", err)
					}
					type_ = filter.DataTypeInt
				} else {
					if _, err := strconv.Atoi(strValue); err != nil {
						return nil, fmt.Errorf("failed to parse filter: %w", err)
					}
				}
			case filter.DataTypeStr:
				type_ = filter.DataTypeInt
				operator = filter.OperatorEq
			}

			if err := filterOptions.AddField(filterName, operator, strValue, type_); err != nil {
				return nil, fmt.Errorf("failed to parse filter: %w", err)
			}
		}
	}
	return filterOptions, nil
}
