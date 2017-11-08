package api

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/satori/go.uuid"
	"net/http"
	"time"
)

type (
	// URLExtractor interface for extracting uuid/country/date from the request
	URLExtractor interface {
		// ExtractUUIDFromRoute extracts the uuid from the url
		UUIDFromRoute(r *http.Request, routeParam string) (*uuid.UUID, error)
		// ExtractDateFromQuery extracts the date from the request
		DateFromQuery(r *http.Request) (*time.Time, error)
		// ExtractDateFromRoute extracts the date from the url
		DateFromRoute(r *http.Request) (*time.Time, error)
	}

	urlExtractor struct{}
)

// NewURLExtractor is the constructor for the UUID Extractor interface
func NewURLExtractor() URLExtractor {
	return &urlExtractor{}
}

func (e *urlExtractor) UUIDFromRoute(r *http.Request, routeParam string) (*uuid.UUID, error) {
	routeUUID, err := uuid.FromString(chi.URLParam(r, routeParam))
	if err != nil {
		return nil, err
	}

	return &routeUUID, nil
}

func (e *urlExtractor) DateFromQuery(r *http.Request) (*time.Time, error) {
	dateStr := r.URL.Query().Get("date")

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, fmt.Errorf("Invalid date %s", dateStr)
	}

	return &date, nil
}

func (e *urlExtractor) DateFromRoute(r *http.Request) (*time.Time, error) {
	dateStr := fmt.Sprintf(
		"%s-%s-%s",
		chi.URLParam(r, "year"),
		chi.URLParam(r, "month"),
		chi.URLParam(r, "day"),
	)

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, fmt.Errorf("Invalid date %s", dateStr)
	}

	return &date, nil
}
