package users_service

import (
	"context"
	"fmt"

	"github.com/timururaimov/golang_todo_app/internal/core/domain"
)

func (s *UsersService) CreatUser(
	ctx context.Context,
	user domain.User,
) (domain.User, error) {
	if err := user.Validate(); err != nil {
		return domain.User{}, fmt.Errorf("validate user domain: %w", err)
	}

	user, err := s.usersRepository.CreatUser(ctx, user)
	if err != nil {
		return domain.User{}, fmt.Errorf("creat user: %w", err)
	}

	return user, nil
}
