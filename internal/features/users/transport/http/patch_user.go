package users_transport_http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/timururaimov/golang_todo_app/internal/core/domain"
	core_logger "github.com/timururaimov/golang_todo_app/internal/core/logger"
	core_http_request "github.com/timururaimov/golang_todo_app/internal/core/request"
	core_http_respone "github.com/timururaimov/golang_todo_app/internal/core/transport/http/respone"
	core_http_types "github.com/timururaimov/golang_todo_app/internal/core/transport/http/types"
	core_http_utils "github.com/timururaimov/golang_todo_app/internal/core/transport/http/utils"
)

type PatchUserRequest struct {
	FullName    core_http_types.Nullable[string] `json:"full_name"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number"`
}

func (r *PatchUserRequest) Validate() error {
	if r.FullName.Set {
		if r.FullName.Value == nil {
			return fmt.Errorf("FullName cant be null")
		}
		fulNameLen := len([]rune(*r.FullName.Value))
		if fulNameLen < 3 || fulNameLen > 100 {
			return fmt.Errorf("FulName must be betwen 3 and 100")
		}
	}
	if r.PhoneNumber.Set {
		if r.PhoneNumber.Value != nil {
			phoneNumverLen := len([]rune(*r.PhoneNumber.Value))
			if phoneNumverLen < 10 || phoneNumverLen > 15 {
				return fmt.Errorf("PhoneNumber must be betwen 10 and 15")
			}
			if !strings.HasPrefix(*r.PhoneNumber.Value, "+") {
				return fmt.Errorf("PhoneNumber must starts with +")
			}
		}
	}
	return nil

}

type PatchUserDTOResponse DTOUserResponse

func (h *UsersHTTPHandler) PatchUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromeContext(ctx)
	responseHandler := core_http_respone.NewHTTPResponseHandler(log, rw)
	userId, err := core_http_utils.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"faild to get userId path value",
		)
		return
	}

	var request PatchUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"faild to decode request",
		)
		return
	}

	userPatch := userPatchFromeRequest(request)

	userDomain, err := h.usersService.PatchUser(ctx, userId, userPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"faild to patch user",
		)
		return
	}
	response := PatchUserDTOResponse(userDTOFromDomain(userDomain))
	responseHandler.JSONRespone(response, http.StatusOK)

}

func userPatchFromeRequest(request PatchUserRequest) domain.UserPatch {
	return domain.UserPatch{
		FullName:    request.FullName.ToDomain(),
		PhoneNumber: request.PhoneNumber.ToDomain(),
	}
}
