package controller

import (
	"context"

	"github.com/alvinmatias69/editor-history/internal/entity"
)

type documentRepository interface {
	Create(context.Context, entity.Document) error
	Finish(context.Context, string) error
	GetForRead(context.Context, string) (entity.Document, error)
	GetForWrite(context.Context, string) (entity.Document, error)
}

type documentOperationsRepository interface {
	Get(context.Context, string) ([]entity.DocumentOperation, error)
	Add(context.Context, []entity.DocumentOperation) (uint64, error)
}

type sessionRepository interface {
	Get(context.Context, string) (string, error)
	Set(context.Context, entity.SessionRequest) error
}
