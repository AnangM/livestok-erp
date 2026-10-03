package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/AnangM/livestok-erp/internal/domain"
	"github.com/AnangM/livestok-erp/internal/repository/mocks"
	"github.com/AnangM/livestok-erp/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHarvestService_RegisterNewHarvest(t *testing.T) {
	pastDate := time.Now().Add(-24 * time.Hour)
	futureDate := time.Now().Add(24 * time.Hour)

	tests := []struct {
		name        string
		harvest     *domain.Harvest
		setupMock   func(m *mocks.MockHarvestRepository)
		expectedErr string
	}{
		{
			name: "Validation error - empty farm id",
			harvest: &domain.Harvest{
				FarmID:      "",
				AnimalID:    "animal-1",
				Weight:      150.5,
				HarvestDate: pastDate,
			},
			setupMock:   func(m *mocks.MockHarvestRepository) {},
			expectedErr: "farm_id is required",
		},
		{
			name: "Validation error - empty animal id",
			harvest: &domain.Harvest{
				FarmID:      "farm-1",
				AnimalID:    "",
				Weight:      150.5,
				HarvestDate: pastDate,
			},
			setupMock:   func(m *mocks.MockHarvestRepository) {},
			expectedErr: "animal_id is required",
		},
		{
			name: "Validation error - zero weight",
			harvest: &domain.Harvest{
				FarmID:      "farm-1",
				AnimalID:    "animal-1",
				Weight:      0,
				HarvestDate: pastDate,
			},
			setupMock:   func(m *mocks.MockHarvestRepository) {},
			expectedErr: "weight is required",
		},
		{
			name: "Validation error - future harvest date",
			harvest: &domain.Harvest{
				FarmID:      "farm-1",
				AnimalID:    "animal-1",
				Weight:      150.5,
				HarvestDate: futureDate,
			},
			setupMock:   func(m *mocks.MockHarvestRepository) {},
			expectedErr: "harvest date cannot be in the future",
		},
		{
			name: "Repository error - create fails",
			harvest: &domain.Harvest{
				FarmID:      "farm-1",
				AnimalID:    "animal-1",
				Weight:      150.5,
				HarvestDate: pastDate,
			},
			setupMock: func(m *mocks.MockHarvestRepository) {
				m.On("Create", mock.Anything, mock.MatchedBy(func(h *domain.Harvest) bool {
					return h.AnimalID == "animal-1" && h.Weight == 150.5
				})).Return(errors.New("db insert harvest error")).Once()
			},
			expectedErr: "db insert harvest error",
		},
		{
			name: "Success - harvest registered",
			harvest: &domain.Harvest{
				FarmID:      "farm-1",
				AnimalID:    "animal-1",
				Weight:      150.5,
				HarvestDate: pastDate,
			},
			setupMock: func(m *mocks.MockHarvestRepository) {
				m.On("Create", mock.Anything, mock.MatchedBy(func(h *domain.Harvest) bool {
					return h.AnimalID == "animal-1" && h.Weight == 150.5
				})).Return(nil).Once()
			},
			expectedErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockHarvestRepository)
			tt.setupMock(mockRepo)

			svc := service.NewHarvestService(mockRepo)
			err := svc.RegisterNewHarvest(context.Background(), tt.harvest)

			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestHarvestService_ListHarvests(t *testing.T) {
	tests := []struct {
		name        string
		farmID      string
		setupMock   func(m *mocks.MockHarvestRepository)
		expectedLen int
		expectedErr string
	}{
		{
			name:   "Repository error - list fails",
			farmID: "farm-1",
			setupMock: func(m *mocks.MockHarvestRepository) {
				m.On("List", mock.Anything, "farm-1").Return(nil, errors.New("query harvests error")).Once()
			},
			expectedLen: 0,
			expectedErr: "query harvests error",
		},
		{
			name:   "Success - returns harvest list",
			farmID: "farm-1",
			setupMock: func(m *mocks.MockHarvestRepository) {
				m.On("List", mock.Anything, "farm-1").Return([]domain.Harvest{
					{ID: "harvest-1", FarmID: "farm-1", AnimalID: "animal-1", Weight: 120.0},
					{ID: "harvest-2", FarmID: "farm-1", AnimalID: "animal-2", Weight: 210.5},
				}, nil).Once()
			},
			expectedLen: 2,
			expectedErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockHarvestRepository)
			tt.setupMock(mockRepo)

			svc := service.NewHarvestService(mockRepo)
			harvests, err := svc.ListHarvests(context.Background(), tt.farmID)

			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, harvests)
			} else {
				assert.NoError(t, err)
				assert.Len(t, harvests, tt.expectedLen)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestHarvestService_GetHarvest(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		setupMock   func(m *mocks.MockHarvestRepository)
		expectedErr string
		expectedID  string
	}{
		{
			name:        "Validation error - empty id",
			id:          "",
			setupMock:   func(m *mocks.MockHarvestRepository) {},
			expectedErr: "harvest_id is required",
		},
		{
			name: "Repository error - get fails",
			id:   "harvest-1",
			setupMock: func(m *mocks.MockHarvestRepository) {
				m.On("Get", mock.Anything, "harvest-1").Return(domain.Harvest{}, errors.New("db error")).Once()
			},
			expectedErr: "db error",
		},
		{
			name: "Success - returns harvest",
			id:   "harvest-1",
			setupMock: func(m *mocks.MockHarvestRepository) {
				m.On("Get", mock.Anything, "harvest-1").Return(domain.Harvest{
					ID:       "harvest-1",
					FarmID:   "farm-1",
					AnimalID: "animal-1",
					Weight:   175.0,
				}, nil).Once()
			},
			expectedErr: "",
			expectedID:  "harvest-1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockHarvestRepository)
			tt.setupMock(mockRepo)

			svc := service.NewHarvestService(mockRepo)
			res, err := svc.GetHarvest(context.Background(), tt.id)

			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, res.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestHarvestService_UpdateHarvest(t *testing.T) {
	pastDate := time.Now().Add(-24 * time.Hour)
	futureDate := time.Now().Add(24 * time.Hour)

	tests := []struct {
		name        string
		id          string
		harvest     *domain.Harvest
		setupMock   func(m *mocks.MockHarvestRepository)
		expectedErr string
		expectedRes domain.Harvest
	}{
		{
			name: "Validation error - empty farm id",
			id:   "harvest-1",
			harvest: &domain.Harvest{
				FarmID:      "",
				AnimalID:    "animal-1",
				Weight:      180.0,
				HarvestDate: pastDate,
			},
			setupMock:   func(m *mocks.MockHarvestRepository) {},
			expectedErr: "farm_id is required",
		},
		{
			name: "Validation error - empty animal id",
			id:   "harvest-1",
			harvest: &domain.Harvest{
				FarmID:      "farm-1",
				AnimalID:    "",
				Weight:      180.0,
				HarvestDate: pastDate,
			},
			setupMock:   func(m *mocks.MockHarvestRepository) {},
			expectedErr: "animal_id is required",
		},
		{
			name: "Validation error - zero weight",
			id:   "harvest-1",
			harvest: &domain.Harvest{
				FarmID:      "farm-1",
				AnimalID:    "animal-1",
				Weight:      0,
				HarvestDate: pastDate,
			},
			setupMock:   func(m *mocks.MockHarvestRepository) {},
			expectedErr: "weight is required",
		},
		{
			name: "Validation error - future harvest date",
			id:   "harvest-1",
			harvest: &domain.Harvest{
				FarmID:      "farm-1",
				AnimalID:    "animal-1",
				Weight:      180.0,
				HarvestDate: futureDate,
			},
			setupMock:   func(m *mocks.MockHarvestRepository) {},
			expectedErr: "harvest date cannot be in the future",
		},
		{
			name: "Repository error - get fails",
			id:   "harvest-1",
			harvest: &domain.Harvest{
				FarmID:      "farm-1",
				AnimalID:    "animal-1",
				Weight:      180.0,
				HarvestDate: pastDate,
			},
			setupMock: func(m *mocks.MockHarvestRepository) {
				m.On("Get", mock.Anything, "harvest-1").Return(domain.Harvest{}, errors.New("db error")).Once()
			},
			expectedErr: "db error",
		},
		{
			name: "Error - harvest not found",
			id:   "harvest-1",
			harvest: &domain.Harvest{
				FarmID:      "farm-1",
				AnimalID:    "animal-1",
				Weight:      180.0,
				HarvestDate: pastDate,
			},
			setupMock: func(m *mocks.MockHarvestRepository) {
				m.On("Get", mock.Anything, "harvest-1").Return(domain.Harvest{}, nil).Once()
			},
			expectedErr: "harvest not found",
		},
		{
			name: "Repository error - update fails",
			id:   "harvest-1",
			harvest: &domain.Harvest{
				FarmID:      "farm-1",
				AnimalID:    "animal-1",
				Weight:      180.0,
				HarvestDate: pastDate,
			},
			setupMock: func(m *mocks.MockHarvestRepository) {
				m.On("Get", mock.Anything, "harvest-1").Return(domain.Harvest{ID: "harvest-1"}, nil).Once()
				m.On("Update", mock.Anything, "harvest-1", mock.Anything).Return(domain.Harvest{}, errors.New("update failed")).Once()
			},
			expectedErr: "update failed",
		},
		{
			name: "Success - harvest updated",
			id:   "harvest-1",
			harvest: &domain.Harvest{
				FarmID:      "farm-1",
				AnimalID:    "animal-1",
				Weight:      195.5,
				HarvestDate: pastDate,
			},
			setupMock: func(m *mocks.MockHarvestRepository) {
				m.On("Get", mock.Anything, "harvest-1").Return(domain.Harvest{
					ID:       "harvest-1",
					FarmID:   "farm-1",
					AnimalID: "animal-1",
					Weight:   150.0,
				}, nil).Once()
				m.On("Update", mock.Anything, "harvest-1", mock.MatchedBy(func(h *domain.Harvest) bool {
					return h.ID == "harvest-1" && h.Weight == 195.5
				})).Return(domain.Harvest{
					ID:       "harvest-1",
					FarmID:   "farm-1",
					AnimalID: "animal-1",
					Weight:   195.5,
				}, nil).Once()
			},
			expectedErr: "",
			expectedRes: domain.Harvest{
				ID:       "harvest-1",
				FarmID:   "farm-1",
				AnimalID: "animal-1",
				Weight:   195.5,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockHarvestRepository)
			tt.setupMock(mockRepo)

			svc := service.NewHarvestService(mockRepo)
			res, err := svc.UpdateHarvest(context.Background(), tt.id, tt.harvest)

			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedRes.Weight, res.Weight)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestHarvestService_DeleteHarvest(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		setupMock   func(m *mocks.MockHarvestRepository)
		expectedErr string
	}{
		{
			name:        "Validation error - empty id",
			id:          "",
			setupMock:   func(m *mocks.MockHarvestRepository) {},
			expectedErr: "harvest_id is required",
		},
		{
			name: "Repository error - get fails",
			id:   "harvest-1",
			setupMock: func(m *mocks.MockHarvestRepository) {
				m.On("Get", mock.Anything, "harvest-1").Return(domain.Harvest{}, errors.New("db error")).Once()
			},
			expectedErr: "db error",
		},
		{
			name: "Error - harvest not found",
			id:   "harvest-1",
			setupMock: func(m *mocks.MockHarvestRepository) {
				m.On("Get", mock.Anything, "harvest-1").Return(domain.Harvest{}, nil).Once()
			},
			expectedErr: "harvest not found",
		},
		{
			name: "Repository error - delete fails",
			id:   "harvest-1",
			setupMock: func(m *mocks.MockHarvestRepository) {
				m.On("Get", mock.Anything, "harvest-1").Return(domain.Harvest{ID: "harvest-1"}, nil).Once()
				m.On("Delete", mock.Anything, "harvest-1").Return(errors.New("delete failed")).Once()
			},
			expectedErr: "delete failed",
		},
		{
			name: "Success - harvest deleted",
			id:   "harvest-1",
			setupMock: func(m *mocks.MockHarvestRepository) {
				m.On("Get", mock.Anything, "harvest-1").Return(domain.Harvest{ID: "harvest-1"}, nil).Once()
				m.On("Delete", mock.Anything, "harvest-1").Return(nil).Once()
			},
			expectedErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockHarvestRepository)
			tt.setupMock(mockRepo)

			svc := service.NewHarvestService(mockRepo)
			err := svc.DeleteHarvest(context.Background(), tt.id)

			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
