package service

import (
	"context"

	pb "github.com/i-b8o/read-only_contracts/pb/writable/v1"
)

type ParagraphStorage interface {
	CreateAll(ctx context.Context, paragraphs []*pb.WriterParagraph) error
	UpdateOne(ctx context.Context, content string, paragraphID uint64) error
	DeleteForChapter(ctx context.Context, chapterID uint64) error
	GetWithHrefs(ctx context.Context, chapterID uint64) ([]*pb.WriterParagraph, error)
}

type paragraphService struct {
	storage ParagraphStorage
}

func NewParagraphService(storage ParagraphStorage) *paragraphService {
	return &paragraphService{storage: storage}
}

func (s *paragraphService) CreateAll(ctx context.Context, paragraphs []*pb.WriterParagraph) error {
	return s.storage.CreateAll(ctx, paragraphs)
}
func (s *paragraphService) UpdateOne(ctx context.Context, content string, paragraphID uint64) error {
	return s.storage.UpdateOne(ctx, content, paragraphID)
}
func (s *paragraphService) DeleteForChapter(ctx context.Context, chapterID uint64) error {
	return s.storage.DeleteForChapter(ctx, chapterID)
}
func (s *paragraphService) GetWithHrefs(ctx context.Context, chapterID uint64) ([]*pb.WriterParagraph, error) {
	return s.storage.GetWithHrefs(ctx, chapterID)
}
