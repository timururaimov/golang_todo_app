package domain

import (
	"fmt"
	"regexp"

	core_errors "github.com/timururaimov/golang_todo_app/internal/core/errors"
)

type User struct {
	ID          int
	Version     int
	FulName     string
	PhoneNumber *string
}

func NewUser(
	id int,
	version int,
	fullName string,
	phoneNumber *string,
) User {
	return User{
		ID:          id,
		Version:     version,
		FulName:     fullName,
		PhoneNumber: phoneNumber,
	}
}
func NewUserUnInizilized(
	fullName string,
	phoneNumber *string,

) User {
	return NewUser(
		UnInitializedId,
		UnInitializedVersion,
		fullName,
		phoneNumber,
	)
}

func (u *User) Validate() error {
	// количество символов в строке
	fullNameLenth := len([]rune(u.FulName))

	if fullNameLenth < 3 || fullNameLenth > 100 {
		return fmt.Errorf(
			"invaild fullname lenth : %d: %w",
			fullNameLenth,
			core_errors.ErrInvalidArgument,
		)
	}

	if u.PhoneNumber != nil {
		phonNumberLenth := len([]rune(*u.PhoneNumber))
		if phonNumberLenth < 10 || phonNumberLenth > 15 {
			return fmt.Errorf(
				"invalid phonNumber lenth : %d: %w",
				phonNumberLenth,
				core_errors.ErrInvalidArgument,
			)
		}
		re := regexp.MustCompile(`^\+[0-9]+$`)
		if !re.MatchString(*u.PhoneNumber) {
			return fmt.Errorf(
				"invaild fullname format : %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}
	return nil
}

type UserPatch struct {
	FullName    Nullable[string]
	PhoneNumber Nullable[string]
}

func (p *UserPatch) Validate() error {
	if p.FullName.Set && p.FullName.Value == nil {
		return fmt.Errorf(
			"FulName can't be null: %w",
			core_errors.ErrInvalidArgument,
		)

	}
	return nil
}

func (u *User) ApplyPatch(patch UserPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("Validate user patch : %w", err)
	}

	tmp := *u
	if patch.FullName.Set {
		tmp.FulName = *patch.FullName.Value
	}

	if patch.PhoneNumber.Set {
		tmp.PhoneNumber = patch.PhoneNumber.Value
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate patch user : %w", err)
	}

	*u = tmp

	return nil

}
