package main

import "errors"

var (
	ErrInvalidCredentials = errors.New("authentication with invalid crendentials")
	ErrInvalidForm        = errors.New("submitted form failed validation")
	ErrInvalidQueryParams = errors.New("invalid query parameters submitted")
	ErrMySQL              = errors.New("failed to execute mySQL statement")
)
