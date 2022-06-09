package repository

import (
	"context"

	"imman/parser_service/internal/entity"

	"github.com/uptrace/bun"
)

type PostRepository struct {
	db *bun.DB
}

func NewPostRepository(db *bun.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r PostRepository) Create(ctx context.Context, posts []entity.Post) error {
	_, err := r.db.NewInsert().Model(&posts).Exec(ctx)
	return err
}
