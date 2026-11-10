package user

import "errors"

// ErrUserNotFound is returned when a user cannot be found in the database
var ErrUserNotFound = errors.New("user not found")
