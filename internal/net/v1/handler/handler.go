package handler

import (
	"fmt"
	"net/http"
	"strconv"

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
