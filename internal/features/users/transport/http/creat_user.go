package users_transport_http

import (
	"net/http"

	"github.com/timururaimov/golang_todo_app/internal/core/domain"
	core_logger "github.com/timururaimov/golang_todo_app/internal/core/logger"
	core_http_request "github.com/timururaimov/golang_todo_app/internal/core/request"
	core_http_respone "github.com/timururaimov/golang_todo_app/internal/core/transport/http/respone"
)

// DTO для creatUser
type CreatUserRequest struct {
	FullName    string  `json:"full_name" validate:"required,min=3,max=100"`
	PhoneNumber *string `json:"phone_number" validate:"omitempty,min=3,max=100,startswith=+"`
}

type CreatUserResponse DTOUserResponse

func (h *UsersHTTPHandler) CreatUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromeContext(ctx)
	responseHandler := core_http_respone.NewHTTPResponseHandler(log, w)

	var request CreatUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "Failed to desode and validate")

		return

	}

	userDomain := domainFromDTO(request)
	userDomain, err := h.usersService.CreatUser(ctx, userDomain)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to cret user")
		return
	}

	response := CreatUserResponse(userDTOFromDomain(userDomain))

	responseHandler.JSONRespone(response, http.StatusCreated)
}

func domainFromDTO(dto CreatUserRequest) domain.User {
	return domain.NewUserUnInizilized(
		dto.FullName, dto.PhoneNumber,
	)
}
