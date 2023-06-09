// Code generated by ogen, DO NOT EDIT.

package documents

import (
	"context"
)

// Handler handles operations described by OpenAPI v3 specification.
type Handler interface {
	// CreateDocument implements createDocument operation.
	//
	// Create document.
	//
	// POST /api/documents/v1/documents
	CreateDocument(ctx context.Context, req *CreateUpdateDocumentView) (CreateDocumentRes, error)
	// CreateDocumentStyle implements createDocumentStyle operation.
	//
	// Create document style.
	//
	// POST /api/documents/v1/documents/{id}/styles
	CreateDocumentStyle(ctx context.Context, req *CreateUpdateStyle, params CreateDocumentStyleParams) (CreateDocumentStyleRes, error)
	// CreateElement implements createElement operation.
	//
	// Create structural element in document.
	//
	// POST /api/documents/v1/documents/{id}/elements
	CreateElement(ctx context.Context, req *CreateUpdateStructuralElement, params CreateElementParams) (CreateElementRes, error)
	// CreateParagraphElement implements createParagraphElement operation.
	//
	// Create paragraph element.
	//
	// POST /api/documents/v1/documents/{id}/elements/{seId}/element/paragraphs
	CreateParagraphElement(ctx context.Context, req *CreateUpdateParagraphElement, params CreateParagraphElementParams) (CreateParagraphElementRes, error)
	// DeleteParagraphElement implements deleteParagraphElement operation.
	//
	// Delete paragraph element.
	//
	// DELETE /api/documents/v1/documents/{id}/elements/{seId}/element/paragraphs/{elementId}
	DeleteParagraphElement(ctx context.Context, params DeleteParagraphElementParams) (DeleteParagraphElementRes, error)
	// DeleteStructuralElementByID implements deleteStructuralElementByID operation.
	//
	// Delete structural element.
	//
	// DELETE /api/documents/v1/documents/{id}/elements/{seId}
	DeleteStructuralElementByID(ctx context.Context, params DeleteStructuralElementByIDParams) (DeleteStructuralElementByIDRes, error)
	// DeleteStyleById implements deleteStyleById operation.
	//
	// Delete style by id.
	//
	// DELETE /api/documents/v1/documents/{id}/styles/{styleId}
	DeleteStyleById(ctx context.Context, params DeleteStyleByIdParams) (DeleteStyleByIdRes, error)
	// GetDocumentById implements getDocumentById operation.
	//
	// Get document by id.
	//
	// GET /api/documents/v1/documents/{id}
	GetDocumentById(ctx context.Context, params GetDocumentByIdParams) (GetDocumentByIdRes, error)
	// GetDocumentStyles implements getDocumentStyles operation.
	//
	// Get document styles.
	//
	// GET /api/documents/v1/documents/{id}/styles
	GetDocumentStyles(ctx context.Context, params GetDocumentStylesParams) (GetDocumentStylesRes, error)
	// GetDocuments implements getDocuments operation.
	//
	// Return paginated dto.
	//
	// GET /api/documents/v1/documents
	GetDocuments(ctx context.Context, params GetDocumentsParams) (GetDocumentsRes, error)
	// GetElements implements getElements operation.
	//
	// Get structural elements in document.
	//
	// GET /api/documents/v1/documents/{id}/elements
	GetElements(ctx context.Context, params GetElementsParams) (GetElementsRes, error)
	// GetParagraphElementByIndexes implements getParagraphElementByIndexes operation.
	//
	// Get paragraphs elements by indexes.
	//
	// GET /api/documents/v1/documents/{id}/elements/{structuralElementIndex}/paragraphs/elements/{paragraphElementIndex}
	GetParagraphElementByIndexes(ctx context.Context, params GetParagraphElementByIndexesParams) (GetParagraphElementByIndexesRes, error)
	// GetParagraphElements implements getParagraphElements operation.
	//
	// Get elements.
	//
	// GET /api/documents/v1/documents/{id}/elements/{seId}
	GetParagraphElements(ctx context.Context, params GetParagraphElementsParams) (GetParagraphElementsRes, error)
	// UpdateDocumentById implements updateDocumentById operation.
	//
	// Update document by id.
	//
	// PUT /api/documents/v1/documents/{id}
	UpdateDocumentById(ctx context.Context, req *CreateUpdateDocumentView, params UpdateDocumentByIdParams) (UpdateDocumentByIdRes, error)
	// UpdateParagraphElement implements updateParagraphElement operation.
	//
	// Update paragraph element.
	//
	// PUT /api/documents/v1/documents/{id}/elements/{seId}/element/paragraphs/{elementId}
	UpdateParagraphElement(ctx context.Context, req *UpdateParagraphElement, params UpdateParagraphElementParams) (UpdateParagraphElementRes, error)
	// UpdateParagraphElementByIndexes implements updateParagraphElementByIndexes operation.
	//
	// Update paragraph element by indexes.
	//
	// PUT /api/documents/v1/documents/{id}/elements/{structuralElementIndex}/paragraphs/elements/{paragraphElementIndex}
	UpdateParagraphElementByIndexes(ctx context.Context, req *UpdateParagraphElement, params UpdateParagraphElementByIndexesParams) (UpdateParagraphElementByIndexesRes, error)
	// UpdateStructuralElement implements updateStructuralElement operation.
	//
	// Update structural element.
	//
	// PUT /api/documents/v1/documents/{id}/elements/{seId}
	UpdateStructuralElement(ctx context.Context, req *UpdateStyleOfStructuralElement, params UpdateStructuralElementParams) (UpdateStructuralElementRes, error)
	// UpdateStyleById implements updateStyleById operation.
	//
	// Update style by id.
	//
	// PUT /api/documents/v1/documents/{id}/styles/{styleId}
	UpdateStyleById(ctx context.Context, req *CreateUpdateStyle, params UpdateStyleByIdParams) (UpdateStyleByIdRes, error)
}

// Server implements http server based on OpenAPI v3 specification and
// calls Handler to handle requests.
type Server struct {
	h Handler
	baseServer
}

// NewServer creates new Server.
func NewServer(h Handler, opts ...ServerOption) (*Server, error) {
	s, err := newServerConfig(opts...).baseServer()
	if err != nil {
		return nil, err
	}
	return &Server{
		h:          h,
		baseServer: s,
	}, nil
}
