package core_http_request

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	core_errors "github.com/timururaimov/golang_todo_app/internal/core/errors"
)

func GetIntQueryParams(r *http.Request, key string) (*int, error) {
	params := r.URL.Query().Get(key)
	if params == "" {
		return nil, nil
	}

	//чтобы строку сделат int
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

func GetDateQueryParam(r *http.Request, key string) (*time.Time, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}

	//формат даты
	layot := "2006-01-02"

	//чтобы сделат из строки дату
	date, err := time.Parse(layot, param)
	if err != nil {
		return nil, fmt.Errorf(
			"param='%s' by key: '%s' not valid date: '%v': %w",
			param,
			key,
			err,
			core_errors.ErrInvalidArgument,
		)
	}
	return &date, nil
}
