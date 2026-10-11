package core_http_request

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
	core_errors "github.com/timururaimov/golang_todo_app/internal/core/errors"
)

var requestValidator = validator.New()

type validateable interface {
	Validate() error
}

func DecodeAndValidateRequest(r *http.Request, dest any) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf("Decode json err:%v , %w", err, core_errors.ErrInvalidArgument)
	}

	var (
		err error
	)

	v, ok := dest.(validateable)
	if ok {
		err = v.Validate()
	} else {
		err = requestValidator.Struct(dest)
	}

	if err != nil {
		return fmt.Errorf("Validate error : %v, %w", err, core_errors.ErrInvalidArgument)
	}
	return nil
}
