package oas

import (
	"github.com/0B1t322/Documents-Service/pkg/gen/open-api/sessions"
	"github.com/0B1t322/Documents-Service/sessions/internal/core/models"
	"github.com/0B1t322/Documents-Service/sessions/internal/core/operations/elements"
	"github.com/samber/lo"
)

type mapper struct{}

func (m mapper) SaveOperations(req *sessions.SaveDocumentRequest) (ops []models.RevisionOperation) {
	for _, command := range req.Commands {
		ops = append(ops, m.commandToOperation(command))
	}

	return
}

func (mapper) commandToOperation(op sessions.DocumentOperation) models.RevisionOperation {
	switch op.Type {
	case sessions.InsertTextDocumentOperation:
		return models.RevisionOperation{
			Insert: lo.ToPtr(
				elements.NewInsertStructuralElement(
					op.InsertText.StructuralElementIndex, op.InsertText.ParagraphElementIndex,
					op.InsertText.InsertBefore, op.InsertText.Content,
				),
			),
		}
	case sessions.DeleteTextDocumentOperation:
		return models.RevisionOperation{
			Delete: lo.ToPtr(
				elements.NewDeleteStructuralElement(
					op.DeleteText.StructuralElementIndex, op.DeleteText.ParagraphElementIndex,
					op.DeleteText.DeleteAfter, op.DeleteText.Content,
				),
			),
		}
	}

	panic("Unknown operation type")
}

func (m mapper) DocumentHistory(documentRevisions []models.DocumentRevision) *sessions.GetDocumentHistoryOKApplicationJSON {
	var revs []sessions.DocumentRevision
	{
		for _, rev := range documentRevisions {
			revs = append(revs, m.DocumentRevision(rev))
		}
	}

	return (*sessions.GetDocumentHistoryOKApplicationJSON)(&revs)
}

func (m mapper) DocumentRevision(rev models.DocumentRevision) sessions.DocumentRevision {
	return sessions.DocumentRevision{
		RevisionId: rev.RevisionID,
		Commands: lo.Map(
			rev.Operations, func(item models.RevisionOperation, index int) sessions.DocumentOperation {
				return m.DocumentOperation(item)
			},
		),
	}
}

func (m mapper) DocumentsOperations(ops []models.RevisionOperation) *sessions.SyncDocumentsByIdOKApplicationJSON {
	var slice []sessions.DocumentOperation
	for _, op := range ops {
		slice = append(slice, m.DocumentOperation(op))
	}

	return (*sessions.SyncDocumentsByIdOKApplicationJSON)(&slice)
}

func (m mapper) DocumentOperation(command models.RevisionOperation) sessions.DocumentOperation {
	switch command.Type() {
	case models.RevisionOperationTypeInsert:
		return sessions.DocumentOperation{
			Type: sessions.InsertTextDocumentOperation,
			InsertText: sessions.InsertText{
				StructuralElementIndex: command.StructuralElementIndex(),
				ParagraphElementIndex:  command.ParagraphElementIndex(),
				InsertBefore:           command.LocalIndex(),
				Content:                command.Insert.Text,
			},
		}
	case models.RevisionOperationTypeDelete:
		return sessions.DocumentOperation{
			Type: sessions.DeleteTextDocumentOperation,
			DeleteText: sessions.DeleteText{
				StructuralElementIndex: command.StructuralElementIndex(),
				ParagraphElementIndex:  command.ParagraphElementIndex(),
				DeleteAfter:            command.LocalIndex(),
				Content:                command.Delete.Text,
			},
		}
	}

	panic("Unknown operation type")
}
