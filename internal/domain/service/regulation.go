package service

import (
	"context"

	pb "github.com/i-b8o/read-only_contracts/pb/writer/v1"
)

type RegulationStorage interface {
	Create(ctx context.Context, regulation *pb.CreateRegulationRequest) (uint64, error)
	Delete(ctx context.Context, regulationID uint64) error
	GetAll(ctx context.Context) (regulations []*pb.WriterRegulation, err error)
}

type regulationService struct {
	storage RegulationStorage
}

func NewRegulationService(storage RegulationStorage) *regulationService {
	return &regulationService{storage: storage}
}

func (s *regulationService) Create(ctx context.Context, regulation *pb.CreateRegulationRequest) (uint64, error) {
	return s.storage.Create(ctx, regulation)
}
func (s *regulationService) Delete(ctx context.Context, regulationID uint64) error {
	return s.storage.Delete(ctx, regulationID)
}
func (s *regulationService) GetAll(ctx context.Context) (regulations []*pb.WriterRegulation, err error) {
	return s.storage.GetAll(ctx)
}
