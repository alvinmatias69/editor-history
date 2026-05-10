package entity

type Document struct {
	ID         string `json:"id" db:"id"`
	ReadOnlyID string `json:"read_only_id" db:"read_only_id"`
	IsFinished bool   `json:"is_finished" db:"is_finished"`
}
