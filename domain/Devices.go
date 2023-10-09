package domain

import "time"

type Devices struct {
	ID        int64     `json:"id" db:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"-" db:"created_at"`
	UpdatedAt time.Time `json:"-" db:"updated_at"`
}
