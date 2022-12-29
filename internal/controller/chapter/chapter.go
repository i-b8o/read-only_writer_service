package chapter_controller

import (
	"context"
	"fmt"

	"github.com/i-b8o/logging"
	pb "github.com/i-b8o/read-only_contracts/pb/writer/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ChapterService interface {
	Create(ctx context.Context, chapter *pb.CreateChapterRequest) (uint64, error)
	GetAllById(ctx context.Context, DocID uint64) ([]uint64, error)
	GetDocId(ctx context.Context, chapterID uint64) (uint64, error)
}

type WriterChapterGrpcController struct {
	chapterService ChapterService
	logging        logging.Logger
	pb.UnimplementedWriterChapterGRPCServer
}

func NewWriterChapterGrpcController(chapterService ChapterService, loging logging.Logger) *WriterChapterGrpcController {
	return &WriterChapterGrpcController{
		chapterService: chapterService,
		logging:        loging,
	}
}

func (c *WriterChapterGrpcController) Create(ctx context.Context, req *pb.CreateChapterRequest) (*pb.CreateChapterResponse, error) {
	ID, err := c.chapterService.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	return &pb.CreateChapterResponse{ID: ID}, nil

}
func (c *WriterChapterGrpcController) GetAll(ctx context.Context, req *pb.GetAllChaptersIdsRequest) (*pb.GetAllChaptersIdsResponse, error) {
	ID := req.GetID()
	IDs, err := c.chapterService.GetAllById(ctx, ID)
	if err != nil {
		return nil, err
	}
	return &pb.GetAllChaptersIdsResponse{IDs: IDs}, nil
}

func (c *WriterChapterGrpcController) GetDocId(ctx context.Context, req *pb.GetDocIdByChapterIdRequest) (*pb.GetDocIdByChapterIdResponse, error) {
	ID := req.GetID()
	regId, err := c.chapterService.GetDocId(ctx, ID)
	if err != nil {
		err := status.Errorf(codes.NotFound, fmt.Sprintf("no rows in result set: %d", ID))
		return &pb.GetDocIdByChapterIdResponse{}, err
	}
	return &pb.GetDocIdByChapterIdResponse{ID: regId}, nil
}
