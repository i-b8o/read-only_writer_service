package usecase

type ChapterService interface {
}

type ParagraphService interface {
}

type chapterUsecase struct {
	chapterService   ChapterService
	paragraphService ParagraphService
}

func NewChapterUsecase(chapterService ChapterService, paragraphService ParagraphService) *chapterUsecase {
	return &chapterUsecase{chapterService: chapterService, paragraphService: paragraphService}
}

func (u *chapterUsecase) Create(ctx, req) {

}
func (u *chapterUsecase) GetAllById(ctx, ID) {

}
func (u *chapterUsecase) DeleteAllForRegulation(ctx, ID) {

}
func (u *chapterUsecase) GetRegulationId(ctx, ID) {

}
