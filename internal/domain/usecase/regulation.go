package usecase

import (
	"context"

	pb "github.com/i-b8o/read-only_contracts/pb/writer/v1"
)

type RegulationService interface {
	Create(ctx context.Context, regulation *pb.CreateRegulationRequest) (uint64, error)
	Delete(ctx context.Context, regulationID uint64) error
}

type regulationUsecase struct {
	service RegulationService
}

func NewRegulationUsecase(service RegulationService) *regulationUsecase {
	return &regulationUsecase{service: service}
}

func (s *regulationUsecase) Create(ctx context.Context, regulation *pb.CreateRegulationRequest) (uint64, error) {
	return s.service.Create(ctx, regulation)
}
func (s *regulationUsecase) Delete(ctx context.Context, regulationID uint64) error {
	return s.service.Delete(ctx, regulationID)
}
