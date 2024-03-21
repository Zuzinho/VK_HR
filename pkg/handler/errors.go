package handler

import (
	"VK_HR/pkg/loginformrepo"
	"fmt"
)

type IncorrectLoginPasswordError struct {
	form loginformrepo.LoginForm
}

func newIncorrectLoginPasswordError(form loginformrepo.LoginForm) IncorrectLoginPasswordError {
	return IncorrectLoginPasswordError{
		form: form,
	}
}

func (err IncorrectLoginPasswordError) Error() string {
	return fmt.Sprintf("incorrec login '%s' or password '%s'", err.form.Login, err.form.Password)
}

type NoUserRoleError struct {
}

func (NoUserRoleError) Error() string {
	return "no user role in context"
}

type NoRequiredParamError struct {
	paramName string
}

func newNoRequiredParamError(paramName string) NoRequiredParamError {
	return NoRequiredParamError{
		paramName: paramName,
	}
}

func (err NoRequiredParamError) Error() string {
	return fmt.Sprintf("no required parameter '%s'", err.paramName)
}

type NoRequiredAccessError struct {
}

func (NoRequiredAccessError) Error() string {
	return "tried admin query, but not Admin role"
}

var (
	IncorrectLoginPasswordErr = IncorrectLoginPasswordError{}
	NoUserRoleErr             = NoUserRoleError{}
	NoRequiredParamErr        = NoRequiredParamError{}
	NoRequiredAccessErr       = NoRequiredAccessError{}
)
