package party

import "errors"

// Enumerate service errors

var (
	ErrAccompanyingGuestsNumberInvalid = errors.New("invalid accompanying guests number")
	ErrGuestAlreadyInList              = errors.New("guest already in list")
	ErrGuestNameRequired               = errors.New("guest name required")
	ErrGuestNotInList                  = errors.New("guest not in list")
	ErrTableNotEnoughSeats             = errors.New("table not enough seats")
	ErrTableNumberInvalid              = errors.New("table number invalid")
	ErrTableNumberNotFound             = errors.New("table number not found")
	ErrTableNumberRequired             = errors.New("table number required")
)
