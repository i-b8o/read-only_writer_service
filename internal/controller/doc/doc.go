package doc_controller

import (
	"context"

	"github.com/i-b8o/logging"
	pb "github.com/i-b8o/read-only_contracts/pb/writer/v1"
)

type DocUsecase interface {
	Create(ctx context.Context, doc *pb.CreateDocRequest) (uint64, error)
	Delete(ctx context.Context, docID uint64) error
	GetAll(ctx context.Context) (docs []*pb.WriterDoc, err error)
}

type WriterDocGrpcController struct {
	docUsecase DocUsecase

	logging logging.Logger
	pb.UnimplementedWriterDocGRPCServer
}

func NewWriterDocGrpcController(docUsecase DocUsecase, loging logging.Logger) *WriterDocGrpcController {
	return &WriterDocGrpcController{
		docUsecase: docUsecase,
		logging:    loging,
	}
}

func (c *WriterDocGrpcController) Create(ctx context.Context, req *pb.CreateDocRequest) (*pb.CreateDocResponse, error) {
	id, err := c.docUsecase.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	return &pb.CreateDocResponse{ID: id}, nil
}

func (c *WriterDocGrpcController) Delete(ctx context.Context, req *pb.DeleteDocRequest) (*pb.Empty, error) {
	ID := req.GetID()
	err := c.docUsecase.Delete(ctx, ID)
	return &pb.Empty{}, err
}

func (c *WriterDocGrpcController) GetAll(ctx context.Context, req *pb.Empty) (*pb.GetDocsResponse, error) {
	docs, err := c.docUsecase.GetAll(ctx)
	return &pb.GetDocsResponse{Docs: docs}, err
}
