package api

import (
	"errors"
	"net/http"
	"strconv"
)

// FetchPagination returns pagination vars from request
func FetchPagination(r *http.Request) (*uint64, *uint64, error) {
	q := r.URL.Query()

	var skip *uint64
	skipStr := q.Get("$skip")
	if skipStr != "" {
		skipInt, err := strconv.ParseUint(skipStr, 10, 64)
		if err != nil {
			return nil, nil, errors.New("Invalid $skip query parameter")
		}
		skip = &skipInt
	}

	var top *uint64
	topStr := r.URL.Query().Get("$top")
	if topStr != "" {
		topInt, err := strconv.ParseUint(topStr, 10, 64)
		if err != nil {
			return nil, nil, errors.New("Invalid $top query parameter")
		}
		top = &topInt
	}

	return skip, top, nil
}
