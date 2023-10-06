package domain

import (
	"context"
	"database/sql"
	"time"
)

// Link is representing the Link data struct
type Link struct {
	ID        int64        `json:"id" db:"id"`
	UserId    int64        `json:"user_id,omitempty" db:"user_id"`
	Alias     string       `json:"alias,omitempty" db:"alias"`
	Target    string       `json:"target" validate:"required" db:"target"`
	CreatedAt time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt sql.NullTime `json:"-" db:"updated_at"`
	DeletedAt sql.NullTime `json:"-" db:"deleted_at"`
}

type LinkRequest struct {
	Target      string   `json:"target" validate:"required,url"`
	Alias       string   `json:"alias,omitempty"`
	Length      int      `json:"length,omitempty" validate:"omitempty,gte=3,lte=10"`
	Description string   `json:"description,omitempty"`
	Tags        []string `json:"tags,omitempty" validate:"dive,required"`
}

// LinkUseCase represent the link's use-cases
type LinkUseCase interface {
	Fetch(ctx context.Context, limit int64) ([]Link, error)
	GetById(ctx context.Context, id int64) (Link, error)
	Update(ctx context.Context, link Link) (int64, error)
	GetByAlias(ctx context.Context, alias string) (Link, error)
	Store(ctx context.Context, link Link) (int64, error)
	Delete(ctx context.Context, id int64) error
}

// LinkRepository represent the link's repository contract
type LinkRepository interface {
	Fetch(ctx context.Context, limit int64) ([]Link, error)
	GetById(ctx context.Context, id int64) (Link, error)
	Update(ctx context.Context, link Link) (int64, error)
	GetByAlias(ctx context.Context, alias string) (Link, error)
	Store(ctx context.Context, link Link) (int64, error)
	Delete(ctx context.Context, id int64) error
}
