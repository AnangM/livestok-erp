package mocks

import (
	"context"

	"github.com/AnangM/livestok-erp/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockHarvestRepository struct {
	mock.Mock
}

func (m *MockHarvestRepository) Create(ctx context.Context, harvest *domain.Harvest) error {
	args := m.Called(ctx, harvest)
	return args.Error(0)
}

func (m *MockHarvestRepository) List(ctx context.Context, farmId string) ([]domain.Harvest, error) {
	args := m.Called(ctx, farmId)
	if harvests, ok := args.Get(0).([]domain.Harvest); ok {
		return harvests, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockHarvestRepository) Update(ctx context.Context, id string, harvest *domain.Harvest) (domain.Harvest, error) {
	args := m.Called(ctx, id, harvest)
	if res, ok := args.Get(0).(domain.Harvest); ok {
		return res, args.Error(1)
	}
	return domain.Harvest{}, args.Error(1)
}

func (m *MockHarvestRepository) Get(ctx context.Context, id string) (domain.Harvest, error) {
	args := m.Called(ctx, id)
	if res, ok := args.Get(0).(domain.Harvest); ok {
		return res, args.Error(1)
	}
	return domain.Harvest{}, args.Error(1)
}

func (m *MockHarvestRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
