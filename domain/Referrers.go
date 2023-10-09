package domain

import "time"

type Referrers struct {
	ID        int64     `json:"id"`
	Hash      string    `json:"hash"`
	Url       string    `json:"url"`
	CreatedAt time.Time `json:"-" db:"created_at"`
	UpdatedAt time.Time `json:"-" db:"updated_at"`
}
