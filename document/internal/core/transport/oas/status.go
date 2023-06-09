package oas

import "github.com/0B1t322/Documents-Service/pkg/gen/open-api/documents"

func Status(httpStatus int, message string) *documents.ErrorStatusCode {
	return &documents.ErrorStatusCode{
		StatusCode: httpStatus,
		Response: documents.Error{
			Status:  httpStatus,
			Message: message,
		},
	}
}
