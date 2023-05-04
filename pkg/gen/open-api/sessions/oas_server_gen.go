// Code generated by ogen, DO NOT EDIT.

package sessions

import (
	"context"
)

// Handler handles operations described by OpenAPI v3 specification.
type Handler interface {
	// GetDocumentHistory implements getDocumentHistory operation.
	//
	// Get history of operations.
	//
	// GET /api/sessions/v1/document/{id}/history
	GetDocumentHistory(ctx context.Context, params GetDocumentHistoryParams) (GetDocumentHistoryRes, error)
	// GetDocumentRevisionId implements getDocumentRevisionId operation.
	//
	// Return current document revision id.
	//
	// GET /api/sessions/v1/documents/{id}/revisionId
	GetDocumentRevisionId(ctx context.Context, params GetDocumentRevisionIdParams) (GetDocumentRevisionIdRes, error)
	// PushOperationToDocument implements pushOperationToDocument operation.
	//
	// Push operation to document.
	//
	// POST /api/sessions/v1/documents/{id}/save
	PushOperationToDocument(ctx context.Context, req *SaveDocumentRequest, params PushOperationToDocumentParams) (PushOperationToDocumentRes, error)
	// SyncDocumentsById implements syncDocumentsById operation.
	//
	// Return all operation that need to apply.
	//
	// GET /api/sessions/v1/documents/{id}/sync
	SyncDocumentsById(ctx context.Context, params SyncDocumentsByIdParams) (SyncDocumentsByIdRes, error)
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
