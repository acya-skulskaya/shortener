package errors

import "errors"

var ErrConflictOriginalURL = errors.New("data conflict: original url already exists")
var ErrConflictID = errors.New("data conflict: id already exists")

var ErrTokenIsNotValid = errors.New("token is not valid")
