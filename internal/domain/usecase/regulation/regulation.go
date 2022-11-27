package regulation_usecase

import (
	"context"

	pb "github.com/i-b8o/read-only_contracts/pb/writer/v1"
)

type RegulationService interface {
	Create(ctx context.Context, regulation *pb.CreateRegulationRequest) (uint64, error)
	Delete(ctx context.Context, regulationID uint64) error
}

type ChapterService interface {
	GetAllById(ctx context.Context, regulationID uint64) ([]uint64, error)
	DeleteAllForRegulation(ctx context.Context, regulationID uint64) error
}

type ParagraphService interface {
	DeleteForChapter(ctx context.Context, chapterID uint64) error
}

type regulationUsecase struct {
	regulationService RegulationService
	chapterService    ChapterService
	paragraphService  ParagraphService
}

func NewRegulationUsecase(regulationService RegulationService, chapterService ChapterService, paragraphService ParagraphService) *regulationUsecase {
	return &regulationUsecase{regulationService: regulationService, chapterService: chapterService, paragraphService: paragraphService}
}

func (u *regulationUsecase) Create(ctx context.Context, regulation *pb.CreateRegulationRequest) (uint64, error) {
	return u.regulationService.Create(ctx, regulation)
}
func (u *regulationUsecase) Delete(ctx context.Context, regulationID uint64) error {
	chIDs, err := u.chapterService.GetAllById(ctx, regulationID)
	if err != nil {
		return err
	}
	for _, chID := range chIDs {
		err := u.paragraphService.DeleteForChapter(ctx, chID)
		if err != nil {
			return err
		}
	}

	err = u.chapterService.DeleteAllForRegulation(ctx, regulationID)
	if err != nil {
		return err
	}
	err = u.regulationService.Delete(ctx, regulationID)
	if err != nil {
		return err
	}
	return nil
}
