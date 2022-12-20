package service

import (
	"context"

	pb "github.com/i-b8o/read-only_contracts/pb/writer/v1"
)

type ChapterStorage interface {
	Create(ctx context.Context, chapter *pb.CreateChapterRequest) (uint64, error)
	DeleteAllForDoc(ctx context.Context, docID uint64) error
	GetAllById(ctx context.Context, docID uint64) ([]uint64, error)
	GetDocId(ctx context.Context, chapterID uint64) (uint64, error)
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

func (s *chapterService) DeleteAllForDoc(ctx context.Context, docID uint64) error {
	return s.storage.DeleteAllForDoc(ctx, docID)
}

func (s *chapterService) GetAllById(ctx context.Context, docID uint64) ([]uint64, error) {
	return s.storage.GetAllById(ctx, docID)
}

func (s *chapterService) GetDocId(ctx context.Context, chapterID uint64) (uint64, error) {
	return s.storage.GetDocId(ctx, chapterID)
}
