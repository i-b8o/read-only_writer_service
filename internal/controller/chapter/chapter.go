package chapter_controller

import (
	"context"
	"fmt"

	"github.com/i-b8o/logging"
	pb "github.com/i-b8o/read-only_contracts/pb/writer/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ChapterUsecase interface {
	Create(ctx context.Context, chapter *pb.CreateChapterRequest) (uint64, error)
	GetAllById(ctx context.Context, regulationID uint64) ([]uint64, error)
	GetRegulationId(ctx context.Context, chapterID uint64) (uint64, error)
}

type WriterChapterGrpcController struct {
	chapterUsecase ChapterUsecase
	logging        logging.Logger
	pb.UnimplementedWriterChapterGRPCServer
}

func NewWriterChapterGrpcController(chapterUsecase ChapterUsecase, loging logging.Logger) *WriterChapterGrpcController {
	return &WriterChapterGrpcController{
		chapterUsecase: chapterUsecase,
		logging:        loging,
	}
}

func (c *WriterChapterGrpcController) Create(ctx context.Context, req *pb.CreateChapterRequest) (*pb.CreateChapterResponse, error) {
	ID, err := c.chapterUsecase.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	return &pb.CreateChapterResponse{ID: ID}, nil

}
func (c *WriterChapterGrpcController) GetAll(ctx context.Context, req *pb.GetAllChaptersIdsRequest) (*pb.GetAllChaptersIdsResponse, error) {
	ID := req.GetID()
	IDs, err := c.chapterUsecase.GetAllById(ctx, ID)
	if err != nil {
		return nil, err
	}
	return &pb.GetAllChaptersIdsResponse{IDs: IDs}, nil
}

func (c *WriterChapterGrpcController) GetRegulationId(ctx context.Context, req *pb.GetRegulationIdByChapterIdRequest) (*pb.GetRegulationIdByChapterIdResponse, error) {
	ID := req.GetID()
	regId, err := c.chapterUsecase.GetRegulationId(ctx, ID)
	if err != nil {
		err := status.Errorf(codes.NotFound, fmt.Sprintf("no rows in result set: %d", ID))
		return &pb.GetRegulationIdByChapterIdResponse{}, err
	}
	return &pb.GetRegulationIdByChapterIdResponse{ID: regId}, nil
}
