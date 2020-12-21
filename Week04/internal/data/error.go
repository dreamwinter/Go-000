package data

import "errors"

var (
	errUserAlreadyExisted = errors.New("User already existed")
	errUserNotFound       = errors.New("Cannot find the user")
)

// IsErrUserNotFound checks whether the given err is errUserNotFound
func IsErrUserNotFound(err error) bool {
	if errors.Is(err, errUserNotFound) {
		return true
	}
	return false
}

// IsErrUserAlreadyExisted checks whether the given err is errUserAlreadyExisted
func IsErrUserAlreadyExisted(err error) bool {
	if errors.Is(err, errUserAlreadyExisted) {
		return true
	}
	return false
}
