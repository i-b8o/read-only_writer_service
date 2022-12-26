package service

import (
	"context"
)

type TypeStorage interface {
	Create(ctx context.Context, name string) (uint64, error)
}

type typeService struct {
	storage TypeStorage
}

func NewTypeService(storage TypeStorage) *typeService {
	return &typeService{storage: storage}
}

func (s *typeService) Create(ctx context.Context, typeName string) (uint64, error) {
	return s.storage.Create(ctx, typeName)
}
