package hikaku

import "errors"

var (
	ErrValuesInvalid        = errors.New("values are invalid")
	ErrFieldNotExported     = errors.New("field is not exported")
	ErrFieldHasNoName       = errors.New("field has no name")
	ErrUnkownKind           = errors.New("kind of value is unknown")
	ErrContextValueNotFound = errors.New("context value not found")
)
