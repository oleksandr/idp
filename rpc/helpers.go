package rpc

import (
	"github.com/oleksandr/idp/errs"
	"github.com/oleksandr/idp/rpc/generated/services"
)

func errorToServiceError(err *errs.Error) error {
	switch err.Type {
	case errs.ErrorTypeForbidden:
		e := services.NewForbiddenError()
		e.Msg = err.Msg
		if err.Cause != nil {
			e.Cause = err.Cause.Error()
		}
		return e

	case errs.ErrorTypeConflict:
		e := services.NewBadRequestError()
		e.Msg = err.Msg
		if err.Cause != nil {
			e.Cause = err.Cause.Error()
		}
		return e

	case errs.ErrorTypeNotFound:
		e := services.NewNotFoundError()
		e.Msg = err.Msg
		if err.Cause != nil {
			e.Cause = err.Cause.Error()
		}
		return e

	case errs.ErrorTypeOperational:
		fallthrough
	default:
		e := services.NewServerError()
		e.Msg = err.Msg
		if err.Cause != nil {
			e.Cause = err.Cause.Error()
		}
		return e
	}
}
