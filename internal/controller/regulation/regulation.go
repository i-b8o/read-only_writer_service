package regulation_controller

import (
	"context"

	"github.com/i-b8o/logging"
	pb "github.com/i-b8o/read-only_contracts/pb/writer/v1"
)

type RegulationUsecase interface {
	Create(ctx context.Context, regulation *pb.CreateRegulationRequest) (uint64, error)
	Delete(ctx context.Context, regulationID uint64) error
	GetAll(ctx context.Context) (regulations []*pb.WriterRegulation, err error)
}

type WriterRegulationGrpcController struct {
	regulationUsecase RegulationUsecase

	logging logging.Logger
	pb.UnimplementedWriterRegulationGRPCServer
}

func NewWriterRegulationGrpcController(regulationUsecase RegulationUsecase, loging logging.Logger) *WriterRegulationGrpcController {
	return &WriterRegulationGrpcController{
		regulationUsecase: regulationUsecase,
		logging:           loging,
	}
}

func (c *WriterRegulationGrpcController) Create(ctx context.Context, req *pb.CreateRegulationRequest) (*pb.CreateRegulationResponse, error) {
	id, err := c.regulationUsecase.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	return &pb.CreateRegulationResponse{ID: id}, nil
}

func (c *WriterRegulationGrpcController) Delete(ctx context.Context, req *pb.DeleteRegulationRequest) (*pb.Empty, error) {
	ID := req.GetID()
	err := c.regulationUsecase.Delete(ctx, ID)
	return &pb.Empty{}, err
}

func (c *WriterRegulationGrpcController) GetAll(ctx context.Context, req *pb.Empty) (*pb.GetRegulationsResponse, error) {
	regulations, err := c.regulationUsecase.GetAll(ctx)
	return &pb.GetRegulationsResponse{Regulations: regulations}, err
}
