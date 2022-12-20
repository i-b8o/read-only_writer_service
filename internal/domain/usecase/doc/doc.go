package doc_usecase

import (
	"context"

	pb "github.com/i-b8o/read-only_contracts/pb/writer/v1"
)

type DocService interface {
	Create(ctx context.Context, doc *pb.CreateDocRequest) (uint64, error)
	Delete(ctx context.Context, docID uint64) error
	GetAll(ctx context.Context) (docs []*pb.WriterDoc, err error)
}

type ChapterService interface {
	GetAllById(ctx context.Context, docID uint64) ([]uint64, error)
	DeleteAllForDoc(ctx context.Context, docID uint64) error
}

type ParagraphService interface {
	DeleteForChapter(ctx context.Context, chapterID uint64) error
}

type docUsecase struct {
	docService       DocService
	chapterService   ChapterService
	paragraphService ParagraphService
}

func NewDocUsecase(docService DocService, chapterService ChapterService, paragraphService ParagraphService) *docUsecase {
	return &docUsecase{docService: docService, chapterService: chapterService, paragraphService: paragraphService}
}
func (u *docUsecase) GetAll(ctx context.Context) (docs []*pb.WriterDoc, err error) {
	return u.docService.GetAll(ctx)
}
func (u *docUsecase) Create(ctx context.Context, doc *pb.CreateDocRequest) (uint64, error) {
	return u.docService.Create(ctx, doc)
}
func (u *docUsecase) Delete(ctx context.Context, docID uint64) error {
	// get all doc`s chapters
	chIDs, err := u.chapterService.GetAllById(ctx, docID)
	if err != nil {
		return err
	}

	// delete all paragraphs
	for _, chID := range chIDs {
		err := u.paragraphService.DeleteForChapter(ctx, chID)
		if err != nil {
			return err
		}
	}

	// delete all chapters
	err = u.chapterService.DeleteAllForDoc(ctx, docID)
	if err != nil {
		return err
	}

	// delete a doc data
	err = u.docService.Delete(ctx, docID)
	if err != nil {
		return err
	}
	return nil
}
