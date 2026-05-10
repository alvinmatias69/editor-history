package entity

import "time"

type GetDocumentResponse struct {
	ID         string                      `json:"id"`
	ReaderID   string                      `json:"reader_id"`
	IsFinished bool                        `json:"is_finished"`
	Operations []DocumentOperationResponse `json:"operations"`
}

type DocumentOperationResponse struct {
	ID                 string    `json:"id"`
	EditorID           string    `json:"editor_id"`
	PositionStart      uint64    `json:"position_start"`
	PositionEnd        uint64    `json:"position_end"`
	Value              string    `json:"value"`
	Operation          string    `json:"operation"`
	OperationTimestamp time.Time `json:"operation_timestamp"`
}
