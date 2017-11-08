package api

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var (
	extractor = NewURLExtractor()
)

func executeRequest(t *testing.T, pattern string, path string, h http.HandlerFunc) {
	req, err := http.NewRequest(http.MethodGet, path, nil)
	assert.NoError(t, err)

	router := chi.NewMux()
	router.Get(pattern, h)

	router.ServeHTTP(httptest.NewRecorder(), req)
}

func TestUrlExtractor_UUIDFromRoute(t *testing.T) {
	t.Parallel()

	pattern := "/{category_id}"
	handler := func(actualUUID *uuid.UUID, err *error) http.HandlerFunc {
		return func(_ http.ResponseWriter, r *http.Request) {
			obtainedUUID, obtainedErr := extractor.UUIDFromRoute(r, "category_id")
			if obtainedUUID != nil {
				*actualUUID = *obtainedUUID
			}
			*err = obtainedErr
		}
	}

	t.Run("It successfully extracts the UUID", func(t *testing.T) {
		expectedUUID := uuid.NewV4().String()
		var actualUUID uuid.UUID
		var err error

		executeRequest(t, pattern, "/"+expectedUUID, handler(&actualUUID, &err))

		assert.NoError(t, err)
		assert.Equal(t, actualUUID.String(), expectedUUID)
	})

	t.Run("It fails on invalid UUID", func(t *testing.T) {
		expectedUUID := "invalid id"
		var actualUUID uuid.UUID
		var err error

		executeRequest(t, pattern, "/"+expectedUUID, handler(&actualUUID, &err))

		assert.NotNil(t, err)
		assert.EqualError(t, err, "uuid: UUID string too short: invalid id")
		assert.Equal(t, uuid.UUID{}, actualUUID)
	})
}

func TestUrlExtractor_DateFromQuery(t *testing.T) {
	t.Parallel()

	pattern := "/"
	handler := func(actualDate *time.Time, err *error) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			obtainedDate, obtainedErr := extractor.DateFromQuery(r)
			if obtainedDate != nil {
				*actualDate = *obtainedDate
			}
			*err = obtainedErr
		}
	}

	t.Run("It successfully extracts the date", func(t *testing.T) {
		var actualDate time.Time
		var err error

		expectedDateStr := "2017-09-07"
		expectedDate, err := time.Parse("2006-01-02", expectedDateStr)
		assert.NoError(t, err)

		executeRequest(t, pattern, "/?date="+expectedDateStr, handler(&actualDate, &err))

		assert.NoError(t, err)
		assert.EqualValues(t, actualDate, expectedDate)
	})

	t.Run("It fails on invalid date", func(t *testing.T) {
		invalidDate := "13th of Yesterday"
		var actualDate time.Time
		var err error

		executeRequest(t, pattern, "/?date="+invalidDate, handler(&actualDate, &err))

		assert.Error(t, err)
		assert.EqualError(t, err, fmt.Sprintf("Invalid date %s", invalidDate))
		assert.Equal(t, time.Time{}, actualDate)
	})

}

func TestUrlExtractor_DateFromRoute(t *testing.T) {
	t.Parallel()

	pattern := "/{year}-{month}-{day}"
	handler := func(actualDate *time.Time, err *error) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			obtainedDate, obtainedErr := extractor.DateFromRoute(r)
			if obtainedDate != nil {
				*actualDate = *obtainedDate
			}
			*err = obtainedErr
		}
	}

	t.Run("It successfully extracts the date", func(t *testing.T) {
		var actualDate time.Time
		var err error

		expectedDateStr := "2017-09-07"
		expectedDate, err := time.Parse("2006-01-02", expectedDateStr)
		assert.NoError(t, err)

		executeRequest(t, pattern, "/"+expectedDateStr, handler(&actualDate, &err))

		assert.Equal(t, err, nil)
		assert.EqualValues(t, actualDate, expectedDate)
	})

	t.Run("It fails on invalid date", func(t *testing.T) {
		invalidDate := "Two-Days-Later"
		var actualDate time.Time
		var err error

		executeRequest(t, pattern, "/"+invalidDate, handler(&actualDate, &err))

		assert.Error(t, err)
		assert.EqualError(t, err, fmt.Sprintf("Invalid date %s", invalidDate))
		assert.Equal(t, time.Time{}, actualDate)
	})
}
