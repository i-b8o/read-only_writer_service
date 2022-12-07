package paragraph_usecase

import (
	"context"

	pb "github.com/i-b8o/read-only_contracts/pb/writer/v1"
)

type ParagraphService interface {
	CreateAll(ctx context.Context, paragraphs []*pb.WriterParagraph) error
	DeleteForChapter(ctx context.Context, chapterID uint64) error
	GetOne(ctx context.Context, paragraphID uint64) (*pb.WriterParagraph, error)
	GetWithHrefs(ctx context.Context, chapterID uint64) ([]*pb.WriterParagraph, error)
	UpdateOne(ctx context.Context, content string, paragraphID uint64) error
}

type paragraphUsecase struct {
	service ParagraphService
}

func NewParagraphUsecase(paragraphService ParagraphService) *paragraphUsecase {
	return &paragraphUsecase{service: paragraphService}
}

func (u *paragraphUsecase) CreateAll(ctx context.Context, paragraphs []*pb.WriterParagraph) error {
	return u.service.CreateAll(ctx, paragraphs)
}

func (u *paragraphUsecase) DeleteForChapter(ctx context.Context, chapterID uint64) error {
	return u.service.DeleteForChapter(ctx, chapterID)
}

func (u *paragraphUsecase) GetOne(ctx context.Context, paragraphID uint64) (*pb.WriterParagraph, error) {
	return u.service.GetOne(ctx, paragraphID)
}

func (u *paragraphUsecase) GetWithHrefs(ctx context.Context, chapterID uint64) ([]*pb.WriterParagraph, error) {
	return u.service.GetWithHrefs(ctx, chapterID)
}

func (u *paragraphUsecase) UpdateOne(ctx context.Context, content string, paragraphID uint64) error {
	return u.service.UpdateOne(ctx, content, paragraphID)
}
