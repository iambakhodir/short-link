package domain

import (
	"context"
	"time"
)

type LinkTag struct {
	ID        int64     `json:"id"`
	LinkId    int64     `json:"link_id"`
	TagId     int64     `json:"tag_id"`
	CreatedAt time.Time `json:"-" db:"created_at"`
	UpdatedAt time.Time `json:"-" db:"updated_at"`
}

// LinkTagUseCase represent the link tag's use-cases
type LinkTagUseCase interface {
	Fetch(ctx context.Context, limit int64) ([]LinkTag, error)
	GetById(ctx context.Context, id int64) (LinkTag, error)
	Update(ctx context.Context, link LinkTag) (int64, error)
	Store(ctx context.Context, link LinkTag) (int64, error)
	Delete(ctx context.Context, id int64) error
}

// LinkTagRepository represent the link tag's repository contract
type LinkTagRepository interface {
	Fetch(ctx context.Context, limit int64) ([]LinkTag, error)
	GetById(ctx context.Context, id int64) (LinkTag, error)
	Update(ctx context.Context, link LinkTag) (int64, error)
	Store(ctx context.Context, link LinkTag) (int64, error)
	Delete(ctx context.Context, id int64) error
}
