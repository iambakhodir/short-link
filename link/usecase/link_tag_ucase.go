package usecase

import (
	"context"
	"github.com/iambakhodir/short-link/domain"
	"time"
)

type linkTagUseCase struct {
	linkTagRepo    domain.LinkTagRepository
	contextTimeout time.Duration
}

func NewLinkTagUseCase(linkTagRepo domain.LinkTagRepository, timeout time.Duration) domain.LinkTagUseCase {
	return &linkTagUseCase{linkTagRepo: linkTagRepo, contextTimeout: timeout}
}

func (lt linkTagUseCase) Fetch(ctx context.Context, limit int64) ([]domain.LinkTag, error) {
	if limit == 0 {
		limit = 10
	} else if limit > 100 {
		limit = 100
	}

	ctx, cancel := context.WithTimeout(ctx, lt.contextTimeout)
	defer cancel()

	res, err := lt.linkTagRepo.Fetch(ctx, limit)
	if err != nil {
		return nil, err
	}

	return res, nil

}

func (lt linkTagUseCase) GetById(ctx context.Context, id int64) (domain.LinkTag, error) {
	ctx, cancel := context.WithTimeout(ctx, lt.contextTimeout)
	defer cancel()

	res, err := lt.linkTagRepo.GetById(ctx, id)
	if err != nil {
		return domain.LinkTag{}, err
	}

	if res == (domain.LinkTag{}) {
		return domain.LinkTag{}, domain.ErrNotFound
	}

	return res, nil
}

func (lt linkTagUseCase) Update(ctx context.Context, linkTag domain.LinkTag) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, lt.contextTimeout)
	defer cancel()

	existedTags, err := lt.linkTagRepo.GetById(ctx, linkTag.ID)
	if err != nil {
		return 0, err
	}

	if existedTags == (domain.LinkTag{}) {
		return 0, domain.ErrNotFound
	}

	linkTag.UpdatedAt = time.Now()

	return lt.linkTagRepo.Update(ctx, linkTag)
}

func (lt linkTagUseCase) Store(ctx context.Context, linkTag domain.LinkTag) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, lt.contextTimeout)
	defer cancel()

	return lt.linkTagRepo.Store(ctx, linkTag)
}

func (lt linkTagUseCase) Delete(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, lt.contextTimeout)
	defer cancel()

	existedTags, err := lt.linkTagRepo.GetById(ctx, id)
	if err != nil {
		return err
	}

	if existedTags == (domain.LinkTag{}) {
		return domain.ErrNotFound
	}

	return lt.linkTagRepo.Delete(ctx, id)
}
