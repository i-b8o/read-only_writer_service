package chapter_usecase

import (
	"context"

	pb "github.com/i-b8o/read-only_contracts/pb/writer/v1"
)

type ChapterService interface {
	Create(ctx context.Context, chapter *pb.CreateChapterRequest) (uint64, error)
	DeleteAllForRegulation(ctx context.Context, regulationID uint64) error
	GetAllById(ctx context.Context, regulationID uint64) ([]uint64, error)
	GetRegulationId(ctx context.Context, chapterID uint64) (uint64, error)
}

type ParagraphService interface {
	DeleteForChapter(ctx context.Context, chapterID uint64) error
}

type chapterUsecase struct {
	chapterService   ChapterService
	paragraphService ParagraphService
}

func NewChapterUsecase(chapterService ChapterService, paragraphService ParagraphService) *chapterUsecase {
	return &chapterUsecase{chapterService: chapterService, paragraphService: paragraphService}
}

func (u *chapterUsecase) Create(ctx context.Context, chapter *pb.CreateChapterRequest) (uint64, error) {
	return u.chapterService.Create(ctx, chapter)
}

func (u *chapterUsecase) GetAllById(ctx context.Context, regulationID uint64) ([]uint64, error) {
	return u.chapterService.GetAllById(ctx, regulationID)
}

func (u *chapterUsecase) DeleteAllForRegulation(ctx context.Context, ID uint64) error {
	chIDs, err := u.chapterService.GetAllById(ctx, ID)
	if err != nil {
		return err
	}
	for _, chID := range chIDs {
		err := u.paragraphService.DeleteForChapter(ctx, chID)
		if err != nil {
			return err
		}
	}

	err = u.chapterService.DeleteAllForRegulation(ctx, ID)
	if err != nil {
		return err
	}
	return nil

}

func (u *chapterUsecase) GetRegulationId(ctx context.Context, chapterID uint64) (uint64, error) {
	return u.chapterService.GetRegulationId(ctx, chapterID)
}
