package types

import "errors"

type ErrorResponse struct {
	Message string `json:"message"`
}

var (
	ErrUnauthorized            = errors.New("Unauthorized")
	ErrInsufficientBalance     = errors.New("Insufficient balance")
	ErrStartAfterEndDate       = errors.New("Start date cannot be after end date")
	ErrPromotionNoLongerActive = errors.New("Promotion is no longer active")
	ErrPromotionExpired        = errors.New("Promotion is expired")
	ErrPromotionNotStarted     = errors.New("Promotion did not start yet")
	ErrRequestorIDNotMatching  = errors.New("Requestor ID is not matching path ID")
	ErrPromotionClaimed        = errors.New("Promotion claimed")
)
