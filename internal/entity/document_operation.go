package entity

import "time"

type DocumentOperation struct {
	ID                 string    `db:"id"`
	DocumentID         string    `db:"document_id"`
	EditorID           string    `db:"editor_id"`
	PositionStart      uint64    `db:"position_start"`
	PositionEnd        uint64    `db:"position_end"`
	Value              string    `db:"val"`
	Operation          string    `db:"op"`
	OperationTimestamp time.Time `db:"operation_time"`
}
