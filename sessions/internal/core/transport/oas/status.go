package oas

import (
	"github.com/0B1t322/Documents-Service/pkg/gen/open-api/sessions"
)

func Status(httpStatus int, message string) *sessions.ErrorStatusCode {
	return &sessions.ErrorStatusCode{
		StatusCode: httpStatus,
		Response: sessions.Error{
			Status:  httpStatus,
			Message: message,
		},
	}
}
