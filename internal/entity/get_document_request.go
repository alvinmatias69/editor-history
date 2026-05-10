package entity

import "github.com/alvinmatias69/editor-history/internal/constants"

type GetDocumentRequest struct {
	DocumentID string
	EditorID   string
	AccessType constants.AccessType
}
