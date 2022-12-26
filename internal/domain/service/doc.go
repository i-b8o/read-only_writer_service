package service

import (
	"context"

	pb "github.com/i-b8o/read-only_contracts/pb/writer/v1"
)

type DocStorage interface {
	Create(ctx context.Context, doc *pb.CreateDocRequest, subtype_id uint64) (uint64, error)
	Delete(ctx context.Context, docID uint64) error
	GetAll(ctx context.Context) (docs []*pb.WriterDoc, err error)
}

type docService struct {
	storage DocStorage
}

func NewDocService(storage DocStorage) *docService {
	return &docService{storage: storage}
}

func (s *docService) Create(ctx context.Context, doc *pb.CreateDocRequest, subtype_id uint64) (uint64, error) {
	return s.storage.Create(ctx, doc, subtype_id)
}
func (s *docService) Delete(ctx context.Context, docID uint64) error {
	return s.storage.Delete(ctx, docID)
}
func (s *docService) GetAll(ctx context.Context) (docs []*pb.WriterDoc, err error) {
	return s.storage.GetAll(ctx)
}
