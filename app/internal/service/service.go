package service

import (
	"context"
	"fmt"

	"regulations_writable_service/internal/pb"

	"github.com/i-b8o/logging"
)

type RegulationStorage interface {
	Create(ctx context.Context, regulation *pb.Regulation) (uint64, error)
	Delete(ctx context.Context, regulationID uint64) error
}

type ChapterStorage interface {
	Create(ctx context.Context, chapter *pb.Chapter) (uint64, error)
	DeleteAllForRegulation(ctx context.Context, regulationID uint64) error
}

type ParagraphStorage interface {
	CreateAll(ctx context.Context, paragraphs []*pb.Paragraph) error
	UpdateOne(ctx context.Context, content string, paragraphID uint64) error
	DeleteForChapter(ctx context.Context, chapterID uint64) error
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

func (s *WritableRegulationGRPCServce) CreateRegulation(ctx context.Context, req *pb.Regulation) (*pb.ID, error) {
	id, err := s.regulationStorage.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	return &pb.ID{ID: id}, nil
}
func (s *WritableRegulationGRPCServce) DeleteRegulation(ctx context.Context, req *pb.ID) (*pb.Status, error) {
	ID := req.GetID()
	err := s.regulationStorage.Delete(ctx, ID)
	if err != nil {
		return &pb.Status{Status: fmt.Sprintf("could not delete the regulation %d", ID)}, err
	}
	return &pb.Status{Status: fmt.Sprintf("regulation %d deleted", ID)}, nil
}
func (s *WritableRegulationGRPCServce) CreateChapter(ctx context.Context, req *pb.Chapter) (*pb.ID, error) {
	ID, err := s.chapterStorage.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	return &pb.ID{ID: ID}, nil

}
func (s *WritableRegulationGRPCServce) DeleteChaptersForRegulation(ctx context.Context, req *pb.ID) (*pb.Status, error) {
	ID := req.GetID()
	err := s.chapterStorage.DeleteAllForRegulation(ctx, ID)
	if err != nil {
		return &pb.Status{Status: fmt.Sprintf("could not delete chapters for the regulation %d", ID)}, err
	}
	return &pb.Status{Status: fmt.Sprintf("chapters for the regulation %d deleted", ID)}, nil
}
func (s *WritableRegulationGRPCServce) CreateAllParagraphs(ctx context.Context, req *pb.Paragraphs) (*pb.Status, error) {
	paragraphs := req.GetParagraphs()
	err := s.paragraphStorage.CreateAll(ctx, paragraphs)
	if err != nil {
		return &pb.Status{Status: fmt.Sprintf("could not create paragraphs for the chapter %d", paragraphs[0].ChapterID)}, err
	}
	return &pb.Status{Status: fmt.Sprintf("paragraphs for the chapter %d created", paragraphs[0].ChapterID)}, nil
}
func (s *WritableRegulationGRPCServce) UpdateOneParagraph(ctx context.Context, req *pb.UpdateOneRequestMesssage) (*pb.Status, error) {
	ID := req.GetID()
	content := req.GetContent()
	err := s.paragraphStorage.UpdateOne(ctx, content, ID)
	if err != nil {
		return &pb.Status{Status: fmt.Sprintf("could not update paragraph %d", ID)}, err
	}
	return &pb.Status{Status: fmt.Sprintf("paragraph %d updated", ID)}, nil
}
func (s *WritableRegulationGRPCServce) DeleteParagraphsForChapter(ctx context.Context, req *pb.ID) (*pb.Status, error) {
	ID := req.GetID()
	err := s.paragraphStorage.DeleteForChapter(ctx, ID)
	if err != nil {
		return &pb.Status{Status: fmt.Sprintf("could not delete paragraph %d", ID)}, err
	}
	return &pb.Status{Status: fmt.Sprintf("paragraph %d deleted", ID)}, nil

}
