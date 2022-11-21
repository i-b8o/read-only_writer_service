package service

import (
	"context"

	pb "github.com/i-b8o/read-only_contracts/pb/writable/v1"
)

type ChapterStorage interface {
	Create(ctx context.Context, chapter *pb.CreateChapterRequest) (uint64, error)
	DeleteAllForRegulation(ctx context.Context, regulationID uint64) error
	GetAllById(ctx context.Context, regulationID uint64) ([]uint64, error)
	GetRegulationId(ctx context.Context, chapterID uint64) (uint64, error)
}

type chapterService struct {
	storage ChapterStorage
}

func NewChapterService(storage ChapterStorage) *chapterService {
	return &chapterService{storage: storage}
}

func (s *chapterService) Create(ctx context.Context, chapter *pb.CreateChapterRequest) (uint64, error) {
	return s.storage.Create(ctx, chapter)
}
func (s *chapterService) DeleteAllForRegulation(ctx context.Context, regulationID uint64) error {
	return s.storage.DeleteAllForRegulation(ctx, regulationID)
}
func (s *chapterService) GetAllById(ctx context.Context, regulationID uint64) ([]uint64, error) {
	return s.storage.GetAllById(ctx, regulationID)
}
func (s *chapterService) GetRegulationId(ctx context.Context, chapterID uint64) (uint64, error) {
	return s.storage.GetRegulationId(ctx, chapterID)
}
