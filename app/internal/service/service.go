package service

import (
	"context"
	"fmt"

	pb "github.com/i-b8o/regulations_contracts/pb/writable/v1"

	"github.com/i-b8o/logging"
)

type RegulationStorage interface {
	Create(ctx context.Context, regulation *pb.CreateRegulationRequest) (uint64, error)
	Delete(ctx context.Context, regulationID uint64) error
}

type ChapterStorage interface {
	Create(ctx context.Context, chapter *pb.CreateChapterRequest) (uint64, error)
	DeleteAllForRegulation(ctx context.Context, regulationID uint64) error
	GetAllById(ctx context.Context, regulationID uint64) ([]*pb.WritableChapter, error)
	GetRegulationId(ctx context.Context, chapterID uint64) (uint64, error)
}

type ParagraphStorage interface {
	CreateAll(ctx context.Context, paragraphs []*pb.WritableParagraph) error
	UpdateOne(ctx context.Context, content string, paragraphID uint64) error
	DeleteForChapter(ctx context.Context, chapterID uint64) error
	GetWithHrefs(ctx context.Context, chapterID uint64) ([]*pb.WritableParagraph, error)
}

type WritableRegulationGRPCServce struct {
	regulationStorage RegulationStorage
	chapterStorage    ChapterStorage
	paragraphStorage  ParagraphStorage
	logging           logging.Logger
	pb.UnimplementedWritableRegulationGRPCServer
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
func (s *WritableRegulationGRPCServce) DeleteRegulation(ctx context.Context, req *pb.DeleteRegulationRequest) (*pb.DeleteRegulationResponse, error) {
	ID := req.GetID()
	err := s.regulationStorage.Delete(ctx, ID)
	if err != nil {
		return &pb.DeleteRegulationResponse{Status: fmt.Sprintf("could not delete the regulation %d", ID)}, err
	}
	return &pb.DeleteRegulationResponse{Status: fmt.Sprintf("regulation %d deleted", ID)}, nil
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
func (s *WritableRegulationGRPCServce) DeleteChaptersForRegulation(ctx context.Context, req *pb.DeleteChaptersForRegulationRequest) (*pb.DeleteChaptersForRegulationResponse, error) {
	ID := req.GetID()
	err := s.chapterStorage.DeleteAllForRegulation(ctx, ID)
	if err != nil {
		return &pb.DeleteChaptersForRegulationResponse{Status: fmt.Sprintf("could not delete chapters for the regulation %d", ID)}, err
	}
	return &pb.DeleteChaptersForRegulationResponse{Status: fmt.Sprintf("chapters for the regulation %d deleted", ID)}, nil
}
func (s *WritableRegulationGRPCServce) CreateAllParagraphs(ctx context.Context, req *pb.CreateAllParagraphsRequest) (*pb.CreateAllParagraphsResponse, error) {
	paragraphs := req.GetParagraphs()
	err := s.paragraphStorage.CreateAll(ctx, paragraphs)
	if err != nil {
		return &pb.CreateAllParagraphsResponse{Status: fmt.Sprintf("could not create paragraphs for the chapter %d", paragraphs[0].ChapterID)}, err
	}
	return &pb.CreateAllParagraphsResponse{Status: fmt.Sprintf("paragraphs for the chapter %d created", paragraphs[0].ChapterID)}, nil
}
func (s *WritableRegulationGRPCServce) UpdateOneParagraph(ctx context.Context, req *pb.UpdateOneParagraphRequest) (*pb.UpdateOneParagraphResponse, error) {
	ID := req.GetID()
	content := req.GetContent()
	err := s.paragraphStorage.UpdateOne(ctx, content, ID)
	if err != nil {
		return &pb.UpdateOneParagraphResponse{Status: fmt.Sprintf("could not update paragraph %d", ID)}, err
	}
	return &pb.UpdateOneParagraphResponse{Status: fmt.Sprintf("paragraph %d updated", ID)}, nil
}
func (s *WritableRegulationGRPCServce) DeleteParagraphsForChapter(ctx context.Context, req *pb.DeleteParagraphsForChapterRequest) (*pb.DeleteParagraphsForChapterResponse, error) {
	ID := req.GetID()
	err := s.paragraphStorage.DeleteForChapter(ctx, ID)
	if err != nil {
		return &pb.DeleteParagraphsForChapterResponse{Status: fmt.Sprintf("could not delete paragraph %d", ID)}, err
	}
	return &pb.DeleteParagraphsForChapterResponse{Status: fmt.Sprintf("paragraph %d deleted", ID)}, nil

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
