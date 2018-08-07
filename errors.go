package progimage

import (
	"errors"
)

// Common errors across the code base
var (
	ErrInvalidArgument    = errors.New("invalid argument")
	ErrInvalidImageFormat = errors.New("invalid image format")
	ErrSavingImage        = errors.New("unexpected error when saving your image")
)
