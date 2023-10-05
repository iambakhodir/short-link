package domain

import (
	"context"
	"time"
)

// Link is representing the Link data struct
type Link struct {
	ID        int64     `json:"id"`
	UserId    int64     `json:"user_id,omitempty"`
	Alias     string    `json:"alias,omitempty"`
	Target    string    `json:"target" validate:"required"`
	DeletedAt time.Time `json:"deleted_at,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// LinkUsecase represent the link's usecases
type LinkUseCase interface {
	Fetch(ctx context.Context, limit int64) ([]Link, error)
	GetById(ctx context.Context, id int64) (Link, error)
	Update(ctx context.Context, link *Link) error
	GetByAlias(ctx context.Context, alias string) (Link, error)
	Store(ctx context.Context, link *Link) error
	Delete(ctx context.Context, id int64) error
}

// LinkRepository represent the link's repository contract
type LinkRepository interface {
	Fetch(ctx context.Context, limit int64) ([]Link, error)
	GetById(ctx context.Context, id int64) (Link, error)
	Update(ctx context.Context, link *Link) error
	GetByAlias(ctx context.Context, alias string) (Link, error)
	Store(ctx context.Context, link *Link) error
	Delete(ctx context.Context, id int64) error
}
