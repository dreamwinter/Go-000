package dao

import (
	"database/sql"
	"fmt"
	"math/rand"

	"github.com/pkg/errors"
)

// User is the data entity to store basic user info
type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Gender   string `json:"gender"`
}

var userStoreInMemory = []User{
	User{
		ID:       int64(1),
		Username: "John",
		Gender:   "M",
	},
	User{
		ID:       int64(2),
		Username: "Tom",
		Gender:   "M",
	},
	User{
		ID:       int64(3),
		Username: "Jerry",
		Gender:   "M",
	},
	User{
		ID:       int64(4),
		Username: "Jane",
		Gender:   "F",
	},
}

// GetUser is the function to find the user with given id
func GetUser(id int64) (User, error) {
	// add random db error, 10% chance to hit this error
	if rand.Intn(10) == 1 {
		return User{}, errors.Wrap(fmt.Errorf("Unknow DB Error"), "Random error happened")
	}
	// simulate sql call
	for _, user := range userStoreInMemory {
		if user.ID == id {
			return user, nil
		}
	}
	return User{}, errors.Wrap(sql.ErrNoRows, fmt.Sprintf("can not find the user with the id:%d", id))
}
