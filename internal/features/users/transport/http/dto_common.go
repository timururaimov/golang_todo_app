package users_transport_http

import "github.com/timururaimov/golang_todo_app/internal/core/domain"

type DTOUserResponse struct {
	Id          int     `json:"id"`
	Version     int     `json:"version"`
	FullName    string  `json:"full_name"`
	PhoneNumber *string `json:"phone_number"`
}

func userDTOFromDomain(user domain.User) DTOUserResponse {
	return DTOUserResponse{
		Id:          user.ID,
		Version:     user.Version,
		FullName:    user.FulName,
		PhoneNumber: user.PhoneNumber,
	}

}

func usersDTOFromeDomains(users []domain.User) []DTOUserResponse {
	usersDTO := make([]DTOUserResponse, len(users))
	for i, user := range users {
		usersDTO[i] = userDTOFromDomain(user)
	}
	return usersDTO
}
