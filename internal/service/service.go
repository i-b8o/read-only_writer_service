package service

import (
	"context"

	pb "github.com/i-b8o/regulations_contracts/pb/writer/v1"

	"github.com/i-b8o/logging"
)

type RegulationStorage interface {
	Create(ctx context.Context, regulation *pb.CreateRegulationRequest) (uint64, error)
	Delete(ctx context.Context, regulationID uint64) error
}

type ChapterStorage interface {
	Create(ctx context.Context, chapter *pb.CreateChapterRequest) (uint64, error)
	DeleteAllForRegulation(ctx context.Context, regulationID uint64) error
	GetAllById(ctx context.Context, regulationID uint64) ([]*pb.WriterChapter, error)
	GetRegulationId(ctx context.Context, chapterID uint64) (uint64, error)
}

type ParagraphStorage interface {
	CreateAll(ctx context.Context, paragraphs []*pb.WriterParagraph) error
	UpdateOne(ctx context.Context, content string, paragraphID uint64) error
	DeleteForChapter(ctx context.Context, chapterID uint64) error
	GetWithHrefs(ctx context.Context, chapterID uint64) ([]*pb.WriterParagraph, error)
}

type WritableRegulationGRPCServce struct {
	regulationStorage RegulationStorage
	chapterStorage    ChapterStorage
	paragraphStorage  ParagraphStorage
	logging           logging.Logger
	pb.UnimplementedWriterGRPCServer
}

func NewWritableRegulationGRPCService(regulationStorage RegulationStorage, chapterStorage ChapterStorage, paragraphStorage ParagraphStorage, loging logging.Logger) *WritableRegulationGRPCServce {
	return &WritableRegulationGRPCServce{
		regulationStorage: regulationStorage,
		chapterStorage:    chapterStorage,
		paragraphStorage:  paragraphStorage,
		logging:           loging,
	}
}

func (s *WritableRegulationGRPCServce) CreateRegulation(ctx context.Context, req *pb.CreateRegulationRequest) (*pb.CreateRegulationResponse, error) {
	id, err := s.regulationStorage.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	return &pb.CreateRegulationResponse{ID: id}, nil
}
func (s *WritableRegulationGRPCServce) DeleteRegulation(ctx context.Context, req *pb.DeleteRegulationRequest) (*pb.Empty, error) {
	ID := req.GetID()
	err := s.regulationStorage.Delete(ctx, ID)
	return &pb.Empty{}, err
}
func (s *WritableRegulationGRPCServce) CreateChapter(ctx context.Context, req *pb.CreateChapterRequest) (*pb.CreateChapterResponse, error) {
	ID, err := s.chapterStorage.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	return &pb.CreateChapterResponse{ID: ID}, nil

}
func (s *WritableRegulationGRPCServce) GetAllChapters(ctx context.Context, req *pb.GetAllChaptersRequest) (*pb.GetAllChaptersResponse, error) {
	ID := req.GetID()
	chapters, err := s.chapterStorage.GetAllById(ctx, ID)
	if err != nil {
		return nil, err
	}
	return &pb.GetAllChaptersResponse{Chapters: chapters}, nil

}
func (s *WritableRegulationGRPCServce) DeleteChaptersForRegulation(ctx context.Context, req *pb.DeleteChaptersForRegulationRequest) (*pb.Empty, error) {
	ID := req.GetID()
	err := s.chapterStorage.DeleteAllForRegulation(ctx, ID)
	return &pb.Empty{}, err
}
func (s *WritableRegulationGRPCServce) CreateAllParagraphs(ctx context.Context, req *pb.CreateAllParagraphsRequest) (*pb.Empty, error) {
	paragraphs := req.GetParagraphs()
	err := s.paragraphStorage.CreateAll(ctx, paragraphs)
	return &pb.Empty{}, err
}
func (s *WritableRegulationGRPCServce) UpdateOneParagraph(ctx context.Context, req *pb.UpdateOneParagraphRequest) (*pb.Empty, error) {
	ID := req.GetID()
	content := req.GetContent()
	err := s.paragraphStorage.UpdateOne(ctx, content, ID)
	return &pb.Empty{}, err

}
func (s *WritableRegulationGRPCServce) DeleteParagraphsForChapter(ctx context.Context, req *pb.DeleteParagraphsForChapterRequest) (*pb.Empty, error) {
	ID := req.GetID()
	err := s.paragraphStorage.DeleteForChapter(ctx, ID)
	return &pb.Empty{}, err
}

func (s *WritableRegulationGRPCServce) GetParagraphsWithHrefs(ctx context.Context, req *pb.GetParagraphsWithHrefsRequest) (*pb.GetParagraphsWithHrefsResponse, error) {
	ID := req.GetID()
	paragraphs, err := s.paragraphStorage.GetWithHrefs(ctx, ID)
	if err != nil {
		return &pb.GetParagraphsWithHrefsResponse{}, err
	}
	return &pb.GetParagraphsWithHrefsResponse{Paragraphs: paragraphs}, nil

}

func (s *WritableRegulationGRPCServce) GetRegulationIdByChapterId(ctx context.Context, req *pb.GetRegulationIdByChapterIdRequest) (*pb.GetRegulationIdByChapterIdResponse, error) {
	ID := req.GetID()
	regId, err := s.chapterStorage.GetRegulationId(ctx, ID)
	if err != nil {
		return &pb.GetRegulationIdByChapterIdResponse{}, err
	}
	return &pb.GetRegulationIdByChapterIdResponse{ID: regId}, nil

}
