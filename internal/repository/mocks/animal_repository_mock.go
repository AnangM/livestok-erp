package mocks

import (
	"context"

	"github.com/AnangM/livestok-erp/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockAnimalRepository struct {
	mock.Mock
}

func (m *MockAnimalRepository) Create(ctx context.Context, animal *domain.Animal) error {
	args := m.Called(ctx, animal)
	return args.Error(0)
}

func (m *MockAnimalRepository) List(ctx context.Context, farmId string) ([]domain.Animal, error) {
	args := m.Called(ctx, farmId)
	if animals, ok := args.Get(0).([]domain.Animal); ok {
		return animals, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAnimalRepository) Update(ctx context.Context, id string, animal *domain.Animal) (domain.Animal, error) {
	args := m.Called(ctx, id, animal)
	if res, ok := args.Get(0).(domain.Animal); ok {
		return res, args.Error(1)
	}
	return domain.Animal{}, args.Error(1)
}

func (m *MockAnimalRepository) Get(ctx context.Context, id string) (domain.Animal, error) {
	args := m.Called(ctx, id)
	if res, ok := args.Get(0).(domain.Animal); ok {
		return res, args.Error(1)
	}
	return domain.Animal{}, args.Error(1)
}

func (m *MockAnimalRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
