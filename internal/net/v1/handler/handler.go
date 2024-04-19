package handler

import (
	"net/http"

	"github.com/4aykovski/effective_mobile_test_task/pkg/response"
	"github.com/go-chi/render"
)

func renderResponse(w http.ResponseWriter, r *http.Request, resp response.Response, statusCode int) {
	render.Status(r, statusCode)
	render.JSON(w, r, resp)
}
