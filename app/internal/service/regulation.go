package service

import (
	"context"

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

type WritableRegulationGRPC struct {
	regulationStorage RegulationStorage
	chapterStorage    ChapterStorage
	paragraphStorage  ParagraphStorage
	logging           logging.Logger
	pb.UnimplementedWritableRegulationGRPCServer
}

func NewWritableRegulationGRPCService(regulationStorage RegulationStorage, chapterStorage ChapterStorage, paragraphStorage ParagraphStorage, loging logging.Logger) *WritableRegulationGRPC {
	return &WritableRegulationGRPC{
		regulationStorage: regulationStorage,
		chapterStorage:    chapterStorage,
		paragraphStorage:  paragraphStorage,
		logging:           loging,
	}
}
