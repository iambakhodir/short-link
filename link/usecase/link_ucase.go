package usecase

import (
	"context"
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

func (lu linkUseCase) Fetch(ctx context.Context, limit int64) ([]domain.Link, error) {
	if limit == 0 {
		limit = 10
	} else if limit > 100 {
		limit = 100
	}

	ctx, cancel := context.WithTimeout(ctx, lu.contextTimeout)
	defer cancel()

	res, err := lu.linkRepo.Fetch(ctx, limit)
	if err != nil {
		return nil, err
	}

	return res, nil

}

func (lu linkUseCase) GetById(ctx context.Context, id int64) (domain.Link, error) {
	ctx, cancel := context.WithTimeout(ctx, lu.contextTimeout)
	defer cancel()

	res, err := lu.linkRepo.GetById(ctx, id)
	if err != nil {
		return domain.Link{}, err
	}

	return res, nil
}

func (lu linkUseCase) Update(ctx context.Context, link domain.Link) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, lu.contextTimeout)
	defer cancel()

	link.UpdatedAt = time.Now()

	return lu.linkRepo.Update(ctx, link)
}

func (lu linkUseCase) GetByAlias(ctx context.Context, alias string) (domain.Link, error) {
	ctx, cancel := context.WithTimeout(ctx, lu.contextTimeout)
	defer cancel()

	res, err := lu.linkRepo.GetByAlias(ctx, alias)
	if err != nil {
		return domain.Link{}, err
	}

	return res, nil
}

func (lu linkUseCase) Store(ctx context.Context, link domain.Link) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, lu.contextTimeout)
	defer cancel()

	return lu.linkRepo.Store(ctx, link)
}

func (lu linkUseCase) Delete(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, lu.contextTimeout)
	defer cancel()

	existedLink, err := lu.linkRepo.GetById(ctx, id)
	if err != nil {
		return err
	}

	if existedLink == (domain.Link{}) {
		return domain.ErrNotFound
	}

	return lu.linkRepo.Delete(ctx, id)
}
