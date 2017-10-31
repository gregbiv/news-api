package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	log "github.com/sirupsen/logrus"
)

// Render is taking care of rendering the Err
func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

// ErrResponse struct
type ErrResponse struct {
	Errors         Error `json:"error"`
	HTTPStatusCode int   `json:"-"` // http response status code
}

// Error struct
type Error struct {
	Code    string `json:"code"`
	Target  string `json:"target,omitempty"`
	Message string `json:"message"`
}

// ErrRender is taking care of rendering the errors correctly
func ErrRender(code, target, message string, statusCode int) render.Renderer {
	return &ErrResponse{
		Errors: Error{
			Code:    code,
			Target:  target,
			Message: message,
		},
		HTTPStatusCode: statusCode,
	}
}

// NotFound is being called when a invalid route is requested.
func NotFound(w http.ResponseWriter, r *http.Request) {
	render.Render(w, r, ErrRender("InvalidUri", "", "The requested URI does not represent any resource on the server.", http.StatusNotFound))
}

// RenderErrMissingURIParam is being called when a query parameter is missing
func RenderErrMissingURIParam(w http.ResponseWriter, r *http.Request, param string) {
	render.Render(
		w,
		r,
		ErrRender("MissingUriParam", "", fmt.Sprintf("The '%s' query parameter is required.", param), http.StatusBadRequest),
	)
}

// RenderInternalServerError is being called when there is an internal server error
func RenderInternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Warn(err)

	render.Render(
		w,
		r,
		ErrRender("InternalError", "", err.Error(), http.StatusInternalServerError),
	)
}

// RenderInvalidInput is being called when the request input is invalid
func RenderInvalidInput(w http.ResponseWriter, r *http.Request, target, message string) {
	render.Render(
		w,
		r,
		ErrRender("InvalidInput", target, message, http.StatusBadRequest),
	)
}

// RenderBadGateway is being called when the server receives an invalid response from another server
func RenderBadGateway(w http.ResponseWriter, r *http.Request, err error) {
	render.Render(
		w,
		r,
		ErrRender("BadGateway", "", err.Error(), http.StatusBadGateway),
	)
}
