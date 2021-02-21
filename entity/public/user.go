package public

import "time"

type User struct {
	ID           int        `json:"id,omitempty"`
	Username     string     `json:"username,omitempty"`
	Email        string     `json:"email,omitempty"`
	RegisteredAt *time.Time `json:"registeredAt,omitempty"`
}
