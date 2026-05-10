package entity

import "time"

type AddOperationsRequest struct {
	DocumentID        string             `json:"document_id"`
	EditorID          string             `json:"editor_id"`
	OperationRequests []OperationRequest `json:"operation_requests"`
}

type OperationRequest struct {
	ID            string    `json:"id"`
	PositionStart uint64    `json:"position_start"`
	PositionEnd   uint64    `json:"position_end"`
	Value         string    `json:"value"`
	Operation     string    `json:"operation"`
	Timestamp     time.Time `json:"timestamp"`
}
