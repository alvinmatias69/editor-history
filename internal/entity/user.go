package entity

type User struct {
	ID       uint64 `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
}
