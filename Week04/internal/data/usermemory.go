package data

import (
	"context"
	"fmt"

	"xll.com/go-000/Week04/internal/biz"

	"github.com/google/uuid"
	xerrors "github.com/pkg/errors"
)

// UserMemoryRepository is the concrete memory implementation of biz.UserRepository
type UserMemoryRepository struct {
	userMap map[string]biz.UserDO
}

// CreateUser adds user into memory repo
func (repo *UserMemoryRepository) CreateUser(ctx context.Context, user biz.UserDO) (biz.UserDO, error) {
	if user.ID == "" {
		exist := true
		for exist {
			// ID collision, regenerate until no collision found
			uid, err := uuidWraper()
			if err != nil {
				return biz.UserDO{}, xerrors.Wrap(err, "Cannot generate uuid")
			}
			user.ID = uid
			_, exist = repo.userMap[user.ID]
		}
	} else {
		if _, exist := repo.userMap[user.ID]; exist {
			return biz.UserDO{}, xerrors.Wrap(errUserAlreadyExisted, fmt.Sprintf("%v id already exists", user.ID))
		}
	}
	repo.userMap[user.ID] = user

	return user, nil
}

// GetUser retrieves the user with the requested id
func (repo *UserMemoryRepository) GetUser(ctx context.Context, id string) (biz.UserDO, error) {
	if u, ok := repo.userMap[id]; ok {
		return u, nil
	}
	return biz.UserDO{}, xerrors.Wrap(errUserNotFound, fmt.Sprintf("%v id cannot be found", id))
}

// UpdateUser updates the user with requested user  data
func (repo *UserMemoryRepository) UpdateUser(ctx context.Context, user biz.UserDO) (biz.UserDO, error) {
	if _, ok := repo.userMap[user.ID]; ok {
		repo.userMap[user.ID] = user
		return user, nil
	}
	return biz.UserDO{}, xerrors.Wrap(errUserNotFound, fmt.Sprintf("%v id cannot be found", user.ID))
}

// DeleteUser deletes the user with requested id
func (repo *UserMemoryRepository) DeleteUser(ctx context.Context, id string) error {
	if _, ok := repo.userMap[id]; ok {
		delete(repo.userMap, id)
		return nil
	}
	return xerrors.Wrap(errUserNotFound, fmt.Sprintf("%v id cannot be found", id))
}

// ListUsers retrieves all the users
func (repo *UserMemoryRepository) ListUsers(ctx context.Context) ([]biz.UserDO, error) {
	userList := make([]biz.UserDO, len(repo.userMap))
	idx := 0
	for _, value := range repo.userMap {
		userList[idx] = value
		idx++
	}
	return userList, nil
}

// NewUserMemoryRepository is the constructor of UserMemoryRepository
func NewUserMemoryRepository() biz.UserRepository {
	userMap := make(map[string]biz.UserDO)
	return &UserMemoryRepository{
		userMap: userMap,
	}
}

var uuidWraper = func() (string, error) {
	uid, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return uid.String(), nil
}
