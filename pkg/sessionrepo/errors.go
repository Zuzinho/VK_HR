package sessionrepo

import (
	"fmt"
)

type InvalidTokenMapClaimsError struct {
	token string
}

func newInvalidTokenMapClaimsError(token string) InvalidTokenMapClaimsError {
	return InvalidTokenMapClaimsError{
		token: token,
	}
}

func (err InvalidTokenMapClaimsError) Error() string {
	return fmt.Sprintf("invalid map claims for token '%s'", err.token)
}

type InvalidTokenError struct {
	token string
}

func newInvalidTokenError(token string) InvalidTokenError {
	return InvalidTokenError{
		token: token,
	}
}

func (err InvalidTokenError) Error() string {
	return fmt.Sprintf("invalid token '%s'", err.token)
}

type NoPayloadKeyError struct {
	key string
}

func newNoPayloadKeyError(key string) NoPayloadKeyError {
	return NoPayloadKeyError{
		key: key,
	}
}

func (err NoPayloadKeyError) Error() string {
	return fmt.Sprintf("no key '%s' in payload", err.key)
}

type InvalidPayloadValueError struct {
	value any
}

func newInvalidPayloadValueError(value any) InvalidPayloadValueError {
	return InvalidPayloadValueError{
		value: value,
	}
}

func (err InvalidPayloadValueError) Error() string {
	return fmt.Sprintf("invalid type of value '%v'", err.value)
}

var (
	NoPayloadKeyErr        = NoPayloadKeyError{}
	InvalidPayloadValueErr = InvalidPayloadValueError{}
)
