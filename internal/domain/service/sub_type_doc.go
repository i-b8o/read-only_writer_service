package service

import (
	"context"
)

type SubTypeDocStorage interface {
	Create(ctx context.Context, subTypeID, docID uint64) error
	Delete(ctx context.Context, docID uint64) error
}

type subTypeDocService struct {
	storage SubTypeDocStorage
}

func NewSubTypeDocService(storage SubTypeDocStorage) *subTypeDocService {
	return &subTypeDocService{storage: storage}
}

func (s *subTypeDocService) Create(ctx context.Context, subTypeID, docID uint64) error {
	return s.storage.Create(ctx, subTypeID, docID)
}
func (s *subTypeDocService) Delete(ctx context.Context, docID uint64) error {
	return s.storage.Delete(ctx, docID)
}
