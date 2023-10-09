package domain

import (
	"context"
	"time"
)

type Visits struct {
	ID          int64     `json:"id"`
	LinkId      int64     `json:"link_id"`
	UserAgentId int64     `json:"user_agent_id"`
	ReferrerId  int64     `json:"referrer_id"`
	Ip          int       `json:"ip"`
	Headers     string    `json:"headers"`
	QueryString string    `json:"query_string"`
	CreatedAt   time.Time `json:"-" db:"created_at"`
	UpdatedAt   time.Time `json:"-" db:"updated_at"`
}

// VisitsUseCase represent the link's use-cases
type VisitsUseCase interface {
	Fetch(ctx context.Context, limit int64) ([]Visits, error)
	GetById(ctx context.Context, id int64) (Visits, error)
	Update(ctx context.Context, visit Visits) (int64, error)
	GetByAlias(ctx context.Context, alias string) (Visits, error)
	Store(ctx context.Context, visit Visits) (int64, error)
	Delete(ctx context.Context, id int64) error
}

// VisitsRepository represent the visit's repository contract
type VisitsRepository interface {
	Fetch(ctx context.Context, limit int64) ([]Visits, error)
	GetById(ctx context.Context, id int64) (Visits, error)
	Update(ctx context.Context, visit Visits) (int64, error)
	GetByAlias(ctx context.Context, alias string) (Visits, error)
	Store(ctx context.Context, visit Visits) (int64, error)
	Delete(ctx context.Context, id int64) error
}
