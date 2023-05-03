package dto

import "github.com/0B1t322/Documents-Service/sessions/internal/core/operations/elements"

type (
	SaveOperationRequestDto struct {
		InsertStructuralElement elements.InsertParagraphElement
		DeleteStructuralElement elements.DeleteStructuralElement
	}
)
