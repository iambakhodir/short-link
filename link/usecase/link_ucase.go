package usecase

import (
	"context"
	"database/sql"
	"github.com/iambakhodir/short-link/domain"
	"time"
)

type linkUseCase struct {
	linkRepo       domain.LinkRepository
	contextTimeout time.Duration
}

func NewLinkUseCase(linkRepo domain.LinkRepository, timeout time.Duration) domain.LinkUseCase {
	return &linkUseCase{linkRepo: linkRepo, contextTimeout: timeout}
}

func (l linkUseCase) Fetch(ctx context.Context, limit int64) ([]domain.Link, error) {
	if limit == 0 {
		limit = 10
	} else if limit > 100 {
		limit = 100
	}

	ctx, cancel := context.WithTimeout(ctx, l.contextTimeout)
	defer cancel()

	res, err := l.linkRepo.Fetch(ctx, limit)
	if err != nil {
		return nil, err
	}

	return res, nil

}

func (l linkUseCase) GetById(ctx context.Context, id int64) (domain.Link, error) {
	ctx, cancel := context.WithTimeout(ctx, l.contextTimeout)
	defer cancel()

	res, err := l.linkRepo.GetById(ctx, id)
	if err != nil {
		return domain.Link{}, err
	}

	return res, nil
}

func (l linkUseCase) Update(ctx context.Context, link domain.Link) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, l.contextTimeout)
	defer cancel()

	link.UpdatedAt = sql.NullTime{Time: time.Now(), Valid: true}

	return l.linkRepo.Update(ctx, link)
}

func (l linkUseCase) GetByAlias(ctx context.Context, alias string) (domain.Link, error) {
	ctx, cancel := context.WithTimeout(ctx, l.contextTimeout)
	defer cancel()

	res, err := l.linkRepo.GetByAlias(ctx, alias)
	if err != nil {
		return domain.Link{}, err
	}

	return res, nil
}

func (l linkUseCase) Store(ctx context.Context, link domain.Link) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, l.contextTimeout)
	defer cancel()

	return l.linkRepo.Store(ctx, link)
}

func (l linkUseCase) Delete(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, l.contextTimeout)
	defer cancel()

	existedLink, err := l.linkRepo.GetById(ctx, id)
	if err != nil {
		return err
	}

	if existedLink == (domain.Link{}) {
		return domain.ErrNotFound
	}

	return l.linkRepo.Delete(ctx, id)
}
