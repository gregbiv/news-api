package middleware

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/render"
	"github.com/gregbiv/news-api/pkg/assets/docs"
	"github.com/gregbiv/news-api/pkg/context"
	"github.com/xeipuuv/gojsonschema"
)

// JSONRequestSchema middleware that will validate the provided request based on a json schema
func JSONRequestSchema(schema string) func(next http.Handler) http.Handler {
	// Create the schema loader
	schemaLoader := gojsonschema.NewBytesLoader(
		docs.MustAsset(fmt.Sprintf("schema/request/%s", schema)),
	)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Load the request body
			requestBody, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, map[string]string{
					"code":    "InternalError",
					"message": "Failed to read the request body",
				})
				return
			}
			requestLoader := gojsonschema.NewBytesLoader(requestBody)

			// Validate the JSON schema
			result, err := gojsonschema.Validate(schemaLoader, requestLoader)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				render.JSON(w, r, map[string]string{
					"code":    "InternalError",
					"message": "Failed to validate against schema",
				})

				context.Logger(r.Context()).Error(err)
				return
			}

			// Handle invalid requests in a nice and pretty way
			if !result.Valid() {
				details := make(map[string]string, len(result.Errors()))
				for _, schemaErr := range result.Errors() {
					details[schemaErr.Field()] = schemaErr.String()
				}

				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, map[string]interface{}{
					"code":    "ResponseError",
					"message": fmt.Sprintf("Response schema validation failed for %s", schema),
					"details": details,
				})
				return
			}

			r.Body = ioutil.NopCloser(bytes.NewReader(requestBody))
			next.ServeHTTP(w, r)
		})
	}
}

// JSONDebugResponseSchema middleware that will validate the provided request based on a json schema
func JSONDebugResponseSchema(schemas map[int]string) func(next http.Handler) http.Handler {
	// Do not enable response validation when not in debug mode
	if !Debug {
		return func(next http.Handler) http.Handler {
			return next
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			buffer := bufferResponse{}

			defer func() {
				// no need to check empty body
				if buffer.resp == http.StatusNoContent {
					return
				}

				// Find the response schema for the status code
				schema, ok := schemas[buffer.resp]
				if !ok {
					buffer.headers.Get("Status")

					w.WriteHeader(http.StatusInternalServerError)
					render.JSON(w, r, map[string]string{
						"code":    "ResponseError",
						"message": fmt.Sprintf("No schema is configured for response code %d", buffer.resp),
					})
					return
				}

				// Create the schema loader
				schemaLoader := gojsonschema.NewBytesLoader(
					docs.MustAsset(fmt.Sprintf("schema/response/%s", schema)),
				)

				// Create the response loader
				responseLoader := gojsonschema.NewBytesLoader(buffer.Bytes())

				// Validate the JSON schema
				result, err := gojsonschema.Validate(schemaLoader, responseLoader)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					render.JSON(w, r, map[string]string{
						"code":    "ResponseError",
						"message": err.Error(),
					})
					return
				}

				if !result.Valid() {
					errors := make([]string, len(result.Errors()))
					for i, schemaErr := range result.Errors() {
						errors[i] = schemaErr.String()
					}

					w.WriteHeader(http.StatusInternalServerError)
					render.JSON(w, r, map[string]interface{}{
						"code":    "ResponseError",
						"message": fmt.Sprintf("Response schema validation failed for %s", schema),
						"errors":  errors,
					})
					return
				}

				context.Logger(r.Context()).Debugf("Valid response based on JSON schema %s", schema)
				buffer.Apply(w)
			}()

			next.ServeHTTP(&buffer, r)
		})
	}
}

// dateFormatChecker define the format date in json schema
type dateFormatChecker struct{}

// IsFormat ensure it meets the gojsonschema.FormatChecker interface
func (f dateFormatChecker) IsFormat(input interface{}) bool {
	asString, ok := input.(string)

	if ok == false {
		return false
	}

	_, err := time.Parse("2006-01-02", asString)
	return err == nil
}

func init() {
	// Add it to the library
	gojsonschema.FormatCheckers.Add("date", dateFormatChecker{})
}

// bufferResponse is a type that implements http.ResponseWriter but buffers all the data and headers.
type bufferResponse struct {
	bytes.Buffer
	resp    int
	headers http.Header
	once    sync.Once
}

// Header implements the header method of http.ResponseWriter
func (b *bufferResponse) Header() http.Header {
	b.once.Do(func() {
		b.headers = make(http.Header)
	})
	return b.headers
}

// WriteHeader implements the WriteHeader method of http.ResponseWriter
func (b *bufferResponse) WriteHeader(resp int) {
	b.resp = resp
}

// Apply takes an http.ResponseWriter and calls the required methods on it to
// output the buffered headers, response code, and data. It returns the number
// of bytes written and any errors flushing.
func (b *bufferResponse) Apply(w http.ResponseWriter) (n int, err error) {
	if len(b.headers) > 0 {
		h := w.Header()
		for key, val := range b.headers {
			h[key] = val
		}
	}
	if b.resp > 0 {
		w.WriteHeader(b.resp)
	}
	n, err = w.Write(b.Bytes())
	return
}
