package domain

import "time"

type UserAgents struct {
	ID        int64     `json:"id"`
	BrowserId int64     `json:"browser_id"`
	DeviceId  int64     `json:"device_id"`
	Hash      string    `json:"hash"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"-" db:"created_at"`
	UpdatedAt time.Time `json:"-" db:"updated_at"`
}
