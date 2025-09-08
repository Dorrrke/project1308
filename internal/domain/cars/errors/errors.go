package errors

import "errors"

var (
	ErrCarsNotFound     = errors.New("sorry, we have no cars =(")
	ErrCarNotFound      = errors.New("car not found")
	ErrNotAvailableCars = errors.New("not available cars, try again later")
	ErrCarNotAvailable  = errors.New("car not available")
)
