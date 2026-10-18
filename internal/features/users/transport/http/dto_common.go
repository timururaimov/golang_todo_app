package users_transport_http

import "github.com/timururaimov/golang_todo_app/internal/core/domain"

type DTOUserResponse struct {
	Id          int     `json:"id" example:"5"`
	Version     int     `json:"version" example:"1"`
	FullName    string  `json:"full_name" example:"Иван Иванов"`
	PhoneNumber *string `json:"phone_number" example:"+9967478888"`
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
