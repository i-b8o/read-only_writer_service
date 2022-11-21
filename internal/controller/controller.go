package controller

import (
	"context"

	pb "github.com/i-b8o/regulations_contracts/pb/writer/v1"

	"github.com/i-b8o/logging"
)

type RegulationUsecase interface {
	Create(ctx context.Context, regulation *pb.CreateRegulationRequest) (uint64, error)
	Delete(ctx context.Context, regulationID uint64) error
}

type ChapterUsecase interface {
	Create(ctx context.Context, chapter *pb.CreateChapterRequest) (uint64, error)
	DeleteAllForRegulation(ctx context.Context, regulationID uint64) error
	GetAllById(ctx context.Context, regulationID uint64) ([]uint64, error)
	GetRegulationId(ctx context.Context, chapterID uint64) (uint64, error)
}

type ParagraphUsecase interface {
	CreateAll(ctx context.Context, paragraphs []*pb.WriterParagraph) error
	UpdateOne(ctx context.Context, content string, paragraphID uint64) error
	DeleteForChapter(ctx context.Context, chapterID uint64) error
	GetWithHrefs(ctx context.Context, chapterID uint64) ([]*pb.WriterParagraph, error)
}

type WritableRegulationGRPCServce struct {
	regulationUsecase RegulationUsecase
	chapterUsecase    ChapterUsecase
	paragraphUsecase  ParagraphUsecase
	logging           logging.Logger
	pb.UnimplementedWriterGRPCServer
}

func NewWritableRegulationGRPCService(regulationUsecase RegulationUsecase, chapterUsecase ChapterUsecase, paragraphStorage ParagraphUsecase, loging logging.Logger) *WritableRegulationGRPCServce {
	return &WritableRegulationGRPCServce{
		regulationUsecase: regulationUsecase,
		chapterUsecase:    chapterUsecase,
		paragraphUsecase:  paragraphStorage,
		logging:           loging,
	}
}

func (s *WritableRegulationGRPCServce) CreateRegulation(ctx context.Context, req *pb.CreateRegulationRequest) (*pb.CreateRegulationResponse, error) {
	id, err := s.regulationUsecase.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	return &pb.CreateRegulationResponse{ID: id}, nil
}
func (s *WritableRegulationGRPCServce) DeleteRegulation(ctx context.Context, req *pb.DeleteRegulationRequest) (*pb.Empty, error) {
	ID := req.GetID()
	err := s.regulationUsecase.Delete(ctx, ID)
	return &pb.Empty{}, err
}
func (s *WritableRegulationGRPCServce) CreateChapter(ctx context.Context, req *pb.CreateChapterRequest) (*pb.CreateChapterResponse, error) {
	ID, err := s.chapterUsecase.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	return &pb.CreateChapterResponse{ID: ID}, nil

}
func (s *WritableRegulationGRPCServce) GetAllChaptersIds(ctx context.Context, req *pb.GetAllChaptersIdsRequest) (*pb.GetAllChaptersIdsResponse, error) {
	ID := req.GetID()
	IDs, err := s.chapterUsecase.GetAllById(ctx, ID)
	if err != nil {
		return nil, err
	}
	return &pb.GetAllChaptersIdsResponse{IDs: IDs}, nil

}
func (s *WritableRegulationGRPCServce) DeleteChaptersForRegulation(ctx context.Context, req *pb.DeleteChaptersForRegulationRequest) (*pb.Empty, error) {
	ID := req.GetID()
	err := s.chapterUsecase.DeleteAllForRegulation(ctx, ID)
	return &pb.Empty{}, err
}
func (s *WritableRegulationGRPCServce) CreateAllParagraphs(ctx context.Context, req *pb.CreateAllParagraphsRequest) (*pb.Empty, error) {
	paragraphs := req.GetParagraphs()
	err := s.paragraphUsecase.CreateAll(ctx, paragraphs)
	return &pb.Empty{}, err
}
func (s *WritableRegulationGRPCServce) UpdateOneParagraph(ctx context.Context, req *pb.UpdateOneParagraphRequest) (*pb.Empty, error) {
	ID := req.GetID()
	content := req.GetContent()
	err := s.paragraphUsecase.UpdateOne(ctx, content, ID)
	return &pb.Empty{}, err

}
func (s *WritableRegulationGRPCServce) DeleteParagraphsForChapter(ctx context.Context, req *pb.DeleteParagraphsForChapterRequest) (*pb.Empty, error) {
	ID := req.GetID()
	err := s.paragraphUsecase.DeleteForChapter(ctx, ID)
	return &pb.Empty{}, err
}

func (s *WritableRegulationGRPCServce) GetParagraphsWithHrefs(ctx context.Context, req *pb.GetParagraphsWithHrefsRequest) (*pb.GetParagraphsWithHrefsResponse, error) {
	ID := req.GetID()
	paragraphs, err := s.paragraphUsecase.GetWithHrefs(ctx, ID)
	if err != nil {
		return &pb.GetParagraphsWithHrefsResponse{}, err
	}
	return &pb.GetParagraphsWithHrefsResponse{Paragraphs: paragraphs}, nil

}

func (s *WritableRegulationGRPCServce) GetRegulationIdByChapterId(ctx context.Context, req *pb.GetRegulationIdByChapterIdRequest) (*pb.GetRegulationIdByChapterIdResponse, error) {
	ID := req.GetID()
	regId, err := s.chapterUsecase.GetRegulationId(ctx, ID)
	if err != nil {
		return &pb.GetRegulationIdByChapterIdResponse{}, err
	}
	return &pb.GetRegulationIdByChapterIdResponse{ID: regId}, nil

}
