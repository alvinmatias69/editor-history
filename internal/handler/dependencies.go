package handler

import (
	"context"

	"github.com/alvinmatias69/editor-history/internal/entity"
)

type documentController interface {
	Create(context.Context) (entity.Document, error)
	Finish(context.Context, entity.SessionRequest) error
	Get(context.Context, entity.GetDocumentRequest) (entity.GetDocumentResponse, error)
	AddOperations(context.Context, entity.AddOperationsRequest) error
	OverrideEditor(context.Context, entity.SessionRequest) error
}
