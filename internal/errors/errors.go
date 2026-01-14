package errors

import "errors"

var ErrConflictOriginalURL = errors.New("data conflict: original url already exists")
var ErrConflictID = errors.New("data conflict: id already exists")

var ErrIDNotFound = errors.New("id not found")
var ErrIDDeleted = errors.New("id was deleted")
var ErrUserIDUnauthorized = errors.New("this user is not authorized for this operation")

var ErrTokenIsNotValid = errors.New("token is not valid")
