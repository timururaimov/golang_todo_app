package core_http_utils

import (
	"fmt"
	"net/http"
	"strconv"

	core_errors "github.com/timururaimov/golang_todo_app/internal/core/errors"
)

func GetIntQueryParams(r *http.Request, key string) (*int, error) {
	params := r.URL.Query().Get(key)
	if params == "" {
		return nil, nil
	}

	val, err := strconv.Atoi(params)
	if err != nil {
		return nil, fmt.Errorf(
			"param='%s' by key:'%s' not a validate integer: %v : %w",
			params,
			key,
			err,
			core_errors.ErrInvalidArgument,
		)

	}

	return &val, nil

}
