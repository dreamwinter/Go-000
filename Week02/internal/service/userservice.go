package service

import (
	"database/sql"
	"errors"
	"fmt"

	"xll.com/go-000/Week02/internal/dao"
)

// UserServiceType is the enum for types of UserService
type UserServiceType string

// UserServiceInMemoryType is InMemory UserServiceType
const UserServiceInMemoryType = "InMemory"

var errorUserNotFound = errors.New("user not found")

// UserService is the interface for all user related business logic
type UserService interface {
	GetUser(id int64) (dao.User, error)
}

// InMemoryUserService is the concrete UserService with in memory store
type InMemoryUserService struct {
}

// GetUser returns the user with the give id
func (s *InMemoryUserService) GetUser(id int64) (dao.User, error) {
	user, err := dao.GetUser(id)
	if errors.Is(err, sql.ErrNoRows) {
		return user, errorUserNotFound
	}
	return user, err
}

// NewUserService is the user service factory method
func NewUserService(serviceType UserServiceType) (UserService, error) {
	switch serviceType {
	case UserServiceInMemoryType:
		return &InMemoryUserService{}, nil
	}
	return nil, fmt.Errorf("%v is unknown UserServiceType", serviceType)
}

// IsUserNotFound helps to determine whether an error is errorUserNotFound
func IsUserNotFound(err error) bool {
	return err == errorUserNotFound
}
