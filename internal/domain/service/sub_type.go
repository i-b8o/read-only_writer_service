package service

import (
	"context"
)

type SubTypeStorage interface {
	Create(ctx context.Context, name string, typeID uint64) (uint64, error)
}

type subTypeService struct {
	storage SubTypeStorage
}

func NewSubTypeService(storage SubTypeStorage) *subTypeService {
	return &subTypeService{storage: storage}
}

func (s *subTypeService) Create(ctx context.Context, subTypeName string, typeID uint64) (uint64, error) {
	return s.storage.Create(ctx, subTypeName, typeID)
}
