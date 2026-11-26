package domain

import (
	"fmt"
	"regexp"
	errorsgo "todoapp-try1/internal/core/errors.go"
)

type User struct {
	FullName    string
	PhoneNumber *string
	Version     int
	ID          int
}

func (u *User) Validate() error {
	fullNameLenth := len([]rune(u.FullName))

	if fullNameLenth < 3 || fullNameLenth > 100 {
		return fmt.Errorf(
			"invalid 'FullName' lenth=%d: %w",
			fullNameLenth,
			errorsgo.ErrInvalidArgument,
		)
	}

	if u.PhoneNumber == nil {
		phoneNumberLenth := len([]rune(*u.PhoneNumber))
		if phoneNumberLenth < 10 || phoneNumberLenth > 15 {
			return fmt.Errorf(
				"invalid 'PhoneNumber' lenth=%d: %w",
				phoneNumberLenth,
				errorsgo.ErrInvalidArgument,
			)
		}

		re := regexp.MustCompile(`^\+[0-9]+$`)
		if !re.MatchString(*u.PhoneNumber) {
			return fmt.Errorf(
				"invalid 'PhoneNumber' format=%s: %w",
				*u.PhoneNumber,
				errorsgo.ErrInvalidArgument,
			)
		}
	}

	return nil
}

func NewUser(
	fullName string,
	phoneNumber *string,
	version int,
	id int,
) User {
	return User{
		FullName:    fullName,
		PhoneNumber: phoneNumber,
		Version:     version,
		ID:          id,
	}
}

func NewUserUninitialized(
	fullName string,
	phoneNumber *string,
) User {
	return User{
		FullName:    fullName,
		PhoneNumber: phoneNumber,
		Version:     uninitializedVersion,
		ID:          uninitializedID,
	}
}
