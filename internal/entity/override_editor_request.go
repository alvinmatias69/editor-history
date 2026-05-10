package entity

type SessionRequest struct {
	DocumentID string `json:"document_id"`
	EditorID   string `json:"editor_id"`
}
