package domain

import (
	"context"
	"time"
)

type Tags struct {
	ID        int64     `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"-" db:"created_at"`
	UpdatedAt time.Time `json:"-" db:"updated_at"`
}

// TagsUseCase represent the link's use-cases
type TagsUseCase interface {
	Fetch(ctx context.Context, limit int64) ([]Tags, error)
	GetById(ctx context.Context, id int64) (Tags, error)
	Update(ctx context.Context, link Tags) (int64, error)
	Store(ctx context.Context, link Tags) (int64, error)
	Delete(ctx context.Context, id int64) error
	FetchByLinkId(ctx context.Context, linkId int64) ([]Tags, error)
}

// TagsRepository represent the link's repository contract
type TagsRepository interface {
	Fetch(ctx context.Context, limit int64) ([]Tags, error)
	GetById(ctx context.Context, id int64) (Tags, error)
	Update(ctx context.Context, link Tags) (int64, error)
	Store(ctx context.Context, link Tags) (int64, error)
	Delete(ctx context.Context, id int64) error
	FetchByLinkId(ctx context.Context, linkId int64) ([]Tags, error)
}
