package users_postgres_repository

import "github.com/timururaimov/golang_todo_app/internal/core/domain"

type UserModel struct {
	Id          int
	Version     int
	FullName    string
	PhoneNumber *string
}

func userDomainsFromeModel(users []UserModel) []domain.User {
	userDomains := make([]domain.User, len(users))

	for i, user := range users {
		userDomains[i] = domain.NewUser(
			user.Id,
			user.Version,
			user.FullName,
			user.PhoneNumber,
		)
	}
	return userDomains
}
