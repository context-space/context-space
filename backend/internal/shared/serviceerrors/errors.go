package serviceerrors

import "errors"

var (
	ErrTranslationNotFound = errors.New("translation not found")
)
