package usecase

import (
	"context"
	"github.com/iambakhodir/short-link/domain"
	"time"
)

type tagsUseCase struct {
	tagsRepo       domain.TagsRepository
	contextTimeout time.Duration
}

func NewTagsUseCase(tagsRepo domain.TagsRepository, timeout time.Duration) domain.TagsUseCase {
	return &tagsUseCase{tagsRepo: tagsRepo, contextTimeout: timeout}
}

func (t tagsUseCase) Fetch(ctx context.Context, limit int64) ([]domain.Tags, error) {
	if limit == 0 {
		limit = 10
	} else if limit > 100 {
		limit = 100
	}

	ctx, cancel := context.WithTimeout(ctx, t.contextTimeout)
	defer cancel()

	res, err := t.tagsRepo.Fetch(ctx, limit)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (t tagsUseCase) FetchByLinkId(ctx context.Context, linkId int64) ([]domain.Tags, error) {
	ctx, cancel := context.WithTimeout(ctx, t.contextTimeout)
	defer cancel()

	res, err := t.tagsRepo.FetchByLinkId(ctx, linkId)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (t tagsUseCase) GetById(ctx context.Context, id int64) (domain.Tags, error) {
	ctx, cancel := context.WithTimeout(ctx, t.contextTimeout)
	defer cancel()

	res, err := t.tagsRepo.GetById(ctx, id)
	if err != nil {
		return domain.Tags{}, err
	}

	if res == (domain.Tags{}) {
		return domain.Tags{}, domain.ErrNotFound
	}

	return res, nil
}

func (t tagsUseCase) Update(ctx context.Context, tags domain.Tags) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, t.contextTimeout)
	defer cancel()

	existedTags, err := t.tagsRepo.GetById(ctx, tags.ID)
	if err != nil {
		return 0, err
	}

	if existedTags == (domain.Tags{}) {
		return 0, domain.ErrNotFound
	}

	tags.UpdatedAt = time.Now()

	return t.tagsRepo.Update(ctx, tags)
}

func (t tagsUseCase) Store(ctx context.Context, tags domain.Tags) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, t.contextTimeout)
	defer cancel()

	return t.tagsRepo.Store(ctx, tags)
}

func (t tagsUseCase) Delete(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, t.contextTimeout)
	defer cancel()

	existedTags, err := t.tagsRepo.GetById(ctx, id)
	if err != nil {
		return err
	}

	if existedTags == (domain.Tags{}) {
		return domain.ErrNotFound
	}

	return t.tagsRepo.Delete(ctx, id)
}
