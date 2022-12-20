package paragraph_controller

import (
	"context"

	"github.com/i-b8o/logging"
	pb "github.com/i-b8o/read-only_contracts/pb/writer/v1"
)

type ParagraphUsecase interface {
	CreateAll(ctx context.Context, paragraphs []*pb.WriterParagraph) error
	UpdateOne(ctx context.Context, content string, paragraphID uint64) error
	GetOne(ctx context.Context, paragraphID uint64) (*pb.WriterParagraph, error)
	DeleteForChapter(ctx context.Context, chapterID uint64) error
	GetWithHrefs(ctx context.Context, chapterID uint64) ([]*pb.WriterParagraph, error)
}

type WriterParagraphGrpcController struct {
	paragraphUsecase ParagraphUsecase
	logging          logging.Logger
	pb.UnimplementedWriterParagraphGRPCServer
}

func NewWritableDocGRPCService(paragraphStorage ParagraphUsecase, loging logging.Logger) *WriterParagraphGrpcController {
	return &WriterParagraphGrpcController{
		paragraphUsecase: paragraphStorage,
		logging:          loging,
	}
}

func (c *WriterParagraphGrpcController) CreateAll(ctx context.Context, req *pb.CreateAllParagraphsRequest) (*pb.Empty, error) {
	paragraphs := req.GetParagraphs()
	err := c.paragraphUsecase.CreateAll(ctx, paragraphs)
	return &pb.Empty{}, err
}

func (c *WriterParagraphGrpcController) Update(ctx context.Context, req *pb.UpdateOneParagraphRequest) (*pb.Empty, error) {
	ID := req.GetID()
	content := req.GetContent()
	err := c.paragraphUsecase.UpdateOne(ctx, content, ID)
	return &pb.Empty{}, err

}

func (c *WriterParagraphGrpcController) GetOne(ctx context.Context, req *pb.GetOneParagraphRequest) (*pb.GetOneParagraphResponse, error) {
	ID := req.GetID()
	paragraphs, err := c.paragraphUsecase.GetOne(ctx, ID)
	if err != nil {
		return &pb.GetOneParagraphResponse{}, err
	}
	return &pb.GetOneParagraphResponse{Content: paragraphs.Content}, nil

}

func (c *WriterParagraphGrpcController) GetWithHrefs(ctx context.Context, req *pb.GetParagraphsWithHrefsRequest) (*pb.GetParagraphsWithHrefsResponse, error) {
	ID := req.GetID()
	paragraphs, err := c.paragraphUsecase.GetWithHrefs(ctx, ID)
	if err != nil {
		return &pb.GetParagraphsWithHrefsResponse{}, err
	}
	return &pb.GetParagraphsWithHrefsResponse{Paragraphs: paragraphs}, nil

}
