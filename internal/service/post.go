package service

import (
	"context"
	"time"

	"imman/parser_service/internal/entity"
	"imman/parser_service/protos/protos/parser_pb"
)

type PostRepository interface {
	Create(ctx context.Context, posts []entity.Post) error
}

type PostWebAPI interface {
	GetData(ctx context.Context, page int) ([]entity.Post, error)
}

type PostService struct {
	repo   PostRepository
	webapi PostWebAPI
	parser_pb.UnimplementedPostParserServiceServer
}

func NewPostService(repo PostRepository, webapi PostWebAPI) *PostService {
	return &PostService{
		repo:   repo,
		webapi: webapi,
	}
}

func (s *PostService) ParseData(ctx context.Context, req *parser_pb.Empty) (*parser_pb.ParseDataResponse, error) {
	for i := 1; i <= 50; i++ {
		time.Sleep(100 * time.Millisecond)
		posts, err := s.webapi.GetData(ctx, i)
		if err != nil {
			return &parser_pb.ParseDataResponse{}, err
		}
		err = s.repo.Create(ctx, posts)
		if err != nil {
			return &parser_pb.ParseDataResponse{}, err
		}
	}

	return &parser_pb.ParseDataResponse{}, nil
}
