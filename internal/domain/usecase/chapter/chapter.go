package chapter_usecase

import (
	"context"

	pb "github.com/i-b8o/read-only_contracts/pb/writer/v1"
)

type ChapterService interface {
	Create(ctx context.Context, chapter *pb.CreateChapterRequest) (uint64, error)
	GetAllById(ctx context.Context, regulationID uint64) ([]uint64, error)
	GetRegulationId(ctx context.Context, chapterID uint64) (uint64, error)
}

type chapterUsecase struct {
	chapterService ChapterService
}

func NewChapterUsecase(chapterService ChapterService) *chapterUsecase {
	return &chapterUsecase{chapterService: chapterService}
}

func (u *chapterUsecase) Create(ctx context.Context, chapter *pb.CreateChapterRequest) (uint64, error) {
	return u.chapterService.Create(ctx, chapter)
}

func (u *chapterUsecase) GetAllById(ctx context.Context, regulationID uint64) ([]uint64, error) {
	return u.chapterService.GetAllById(ctx, regulationID)
}

func (u *chapterUsecase) GetRegulationId(ctx context.Context, chapterID uint64) (uint64, error) {
	return u.chapterService.GetRegulationId(ctx, chapterID)
}
