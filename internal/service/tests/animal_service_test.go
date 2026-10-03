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

func TestAnimalService_RegisterNewAnimal(t *testing.T) {
	pastDate := time.Now().Add(-24 * time.Hour)
	futureDate := time.Now().Add(24 * time.Hour)

	tests := []struct {
		name        string
		animal      *domain.Animal
		setupMock   func(m *mocks.MockAnimalRepository)
		expectedErr string
		checkAnimal func(t *testing.T, a *domain.Animal)
	}{
		{
			name: "Validation error - empty tag number",
			animal: &domain.Animal{
				TagNumber: "   ",
				FarmID:    "farm-1",
				BirthDate: pastDate,
			},
			setupMock:   func(m *mocks.MockAnimalRepository) {},
			expectedErr: "tag number is required",
		},
		{
			name: "Validation error - empty farm id",
			animal: &domain.Animal{
				TagNumber: "TAG-001",
				FarmID:    "",
				BirthDate: pastDate,
			},
			setupMock:   func(m *mocks.MockAnimalRepository) {},
			expectedErr: "farm_id is required",
		},
		{
			name: "Validation error - birth date in future",
			animal: &domain.Animal{
				TagNumber: "TAG-001",
				FarmID:    "farm-1",
				BirthDate: futureDate,
			},
			setupMock:   func(m *mocks.MockAnimalRepository) {},
			expectedErr: "birth date cannot be in the future",
		},
		{
			name: "Repository error - create fails",
			animal: &domain.Animal{
				TagNumber: "TAG-001",
				FarmID:    "farm-1",
				BirthDate: pastDate,
			},
			setupMock: func(m *mocks.MockAnimalRepository) {
				m.On("Create", mock.Anything, mock.MatchedBy(func(a *domain.Animal) bool {
					return a.TagNumber == "TAG-001" && a.Status == "Active"
				})).Return(errors.New("db insert error")).Once()
			},
			expectedErr: "db insert error",
		},
		{
			name: "Success - status defaulted to Active",
			animal: &domain.Animal{
				TagNumber: "TAG-001",
				FarmID:    "farm-1",
				Breed:     "Angus",
				Gender:    "Male",
				BirthDate: pastDate,
				Status:    "",
			},
			setupMock: func(m *mocks.MockAnimalRepository) {
				m.On("Create", mock.Anything, mock.MatchedBy(func(a *domain.Animal) bool {
					return a.TagNumber == "TAG-001" && a.Status == "Active"
				})).Return(nil).Once()
			},
			expectedErr: "",
			checkAnimal: func(t *testing.T, a *domain.Animal) {
				assert.Equal(t, "Active", a.Status)
			},
		},
		{
			name: "Success - existing status preserved",
			animal: &domain.Animal{
				TagNumber: "TAG-002",
				FarmID:    "farm-1",
				BirthDate: pastDate,
				Status:    "Quarantine",
			},
			setupMock: func(m *mocks.MockAnimalRepository) {
				m.On("Create", mock.Anything, mock.MatchedBy(func(a *domain.Animal) bool {
					return a.TagNumber == "TAG-002" && a.Status == "Quarantine"
				})).Return(nil).Once()
			},
			expectedErr: "",
			checkAnimal: func(t *testing.T, a *domain.Animal) {
				assert.Equal(t, "Quarantine", a.Status)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockAnimalRepository)
			tt.setupMock(mockRepo)

			svc := service.NewAnimalService(mockRepo)
			err := svc.RegisterNewAnimal(context.Background(), tt.animal)

			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}

			if tt.checkAnimal != nil {
				tt.checkAnimal(t, tt.animal)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAnimalService_ListAnimals(t *testing.T) {
	tests := []struct {
		name        string
		farmID      string
		setupMock   func(m *mocks.MockAnimalRepository)
		expectedLen int
		expectedErr string
	}{
		{
			name:   "Repository error - list fails",
			farmID: "farm-1",
			setupMock: func(m *mocks.MockAnimalRepository) {
				m.On("List", mock.Anything, "farm-1").Return(nil, errors.New("query failed")).Once()
			},
			expectedLen: 0,
			expectedErr: "query failed",
		},
		{
			name:   "Success - returns animals",
			farmID: "farm-1",
			setupMock: func(m *mocks.MockAnimalRepository) {
				m.On("List", mock.Anything, "farm-1").Return([]domain.Animal{
					{ID: "animal-1", TagNumber: "TAG-1", FarmID: "farm-1"},
					{ID: "animal-2", TagNumber: "TAG-2", FarmID: "farm-1"},
				}, nil).Once()
			},
			expectedLen: 2,
			expectedErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockAnimalRepository)
			tt.setupMock(mockRepo)

			svc := service.NewAnimalService(mockRepo)
			animals, err := svc.ListAnimals(context.Background(), tt.farmID)

			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				assert.Nil(t, animals)
			} else {
				assert.NoError(t, err)
				assert.Len(t, animals, tt.expectedLen)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAnimalService_GetAnimal(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		setupMock   func(m *mocks.MockAnimalRepository)
		expectedErr string
		expectedID  string
	}{
		{
			name:        "Validation error - empty id",
			id:          "",
			setupMock:   func(m *mocks.MockAnimalRepository) {},
			expectedErr: "animal id is required",
		},
		{
			name: "Repository error - get fails",
			id:   "animal-1",
			setupMock: func(m *mocks.MockAnimalRepository) {
				m.On("Get", mock.Anything, "animal-1").Return(domain.Animal{}, errors.New("not found in db")).Once()
			},
			expectedErr: "not found in db",
		},
		{
			name: "Success - animal found",
			id:   "animal-1",
			setupMock: func(m *mocks.MockAnimalRepository) {
				m.On("Get", mock.Anything, "animal-1").Return(domain.Animal{
					ID:        "animal-1",
					TagNumber: "TAG-1",
					FarmID:    "farm-1",
				}, nil).Once()
			},
			expectedErr: "",
			expectedID:  "animal-1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockAnimalRepository)
			tt.setupMock(mockRepo)

			svc := service.NewAnimalService(mockRepo)
			animal, err := svc.GetAnimal(context.Background(), tt.id)

			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, animal.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAnimalService_UpdateAnimal(t *testing.T) {
	pastDate := time.Now().Add(-24 * time.Hour)
	futureDate := time.Now().Add(24 * time.Hour)

	tests := []struct {
		name        string
		id          string
		animal      *domain.Animal
		setupMock   func(m *mocks.MockAnimalRepository)
		expectedErr string
		expectedRes domain.Animal
	}{
		{
			name: "Validation error - empty tag number",
			id:   "animal-1",
			animal: &domain.Animal{
				TagNumber: "   ",
				FarmID:    "farm-1",
				BirthDate: pastDate,
			},
			setupMock:   func(m *mocks.MockAnimalRepository) {},
			expectedErr: "tag number is required",
		},
		{
			name: "Validation error - empty farm id",
			id:   "animal-1",
			animal: &domain.Animal{
				TagNumber: "TAG-001",
				FarmID:    "",
				BirthDate: pastDate,
			},
			setupMock:   func(m *mocks.MockAnimalRepository) {},
			expectedErr: "farm_id is required",
		},
		{
			name: "Validation error - future birth date",
			id:   "animal-1",
			animal: &domain.Animal{
				TagNumber: "TAG-001",
				FarmID:    "farm-1",
				BirthDate: futureDate,
			},
			setupMock:   func(m *mocks.MockAnimalRepository) {},
			expectedErr: "birth date cannot be in the future",
		},
		{
			name: "Repository error - get fails",
			id:   "animal-1",
			animal: &domain.Animal{
				TagNumber: "TAG-001",
				FarmID:    "farm-1",
				BirthDate: pastDate,
			},
			setupMock: func(m *mocks.MockAnimalRepository) {
				m.On("Get", mock.Anything, "animal-1").Return(domain.Animal{}, errors.New("db error")).Once()
			},
			expectedErr: "db error",
		},
		{
			name: "Error - animal not found (empty ID returned)",
			id:   "animal-1",
			animal: &domain.Animal{
				TagNumber: "TAG-001",
				FarmID:    "farm-1",
				BirthDate: pastDate,
			},
			setupMock: func(m *mocks.MockAnimalRepository) {
				m.On("Get", mock.Anything, "animal-1").Return(domain.Animal{}, nil).Once()
			},
			expectedErr: "animal not found",
		},
		{
			name: "Repository error - update fails",
			id:   "animal-1",
			animal: &domain.Animal{
				TagNumber: "TAG-001",
				FarmID:    "farm-1",
				BirthDate: pastDate,
			},
			setupMock: func(m *mocks.MockAnimalRepository) {
				m.On("Get", mock.Anything, "animal-1").Return(domain.Animal{ID: "animal-1"}, nil).Once()
				m.On("Update", mock.Anything, "animal-1", mock.Anything).Return(domain.Animal{}, errors.New("update err")).Once()
			},
			expectedErr: "update err",
		},
		{
			name: "Success - updates existing animal",
			id:   "animal-1",
			animal: &domain.Animal{
				TagNumber: "TAG-001-NEW",
				FarmID:    "farm-1",
				Breed:     "Angus",
				Gender:    "Female",
				BirthDate: pastDate,
				Status:    "", // should default to Active
			},
			setupMock: func(m *mocks.MockAnimalRepository) {
				m.On("Get", mock.Anything, "animal-1").Return(domain.Animal{
					ID:        "animal-1",
					FarmID:    "farm-1",
					TagNumber: "TAG-001-OLD",
				}, nil).Once()
				m.On("Update", mock.Anything, "animal-1", mock.MatchedBy(func(a *domain.Animal) bool {
					return a.TagNumber == "TAG-001-NEW" && a.Status == "Active" && a.Breed == "Angus"
				})).Return(domain.Animal{
					ID:        "animal-1",
					TagNumber: "TAG-001-NEW",
					Breed:     "Angus",
					Status:    "Active",
				}, nil).Once()
			},
			expectedErr: "",
			expectedRes: domain.Animal{
				ID:        "animal-1",
				TagNumber: "TAG-001-NEW",
				Breed:     "Angus",
				Status:    "Active",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockAnimalRepository)
			tt.setupMock(mockRepo)

			svc := service.NewAnimalService(mockRepo)
			res, err := svc.UpdateAnimal(context.Background(), tt.id, tt.animal)

			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedRes.TagNumber, res.TagNumber)
				assert.Equal(t, tt.expectedRes.Status, res.Status)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAnimalService_DeleteAnimal(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		setupMock   func(m *mocks.MockAnimalRepository)
		expectedErr string
	}{
		{
			name:        "Validation error - empty id",
			id:          "",
			setupMock:   func(m *mocks.MockAnimalRepository) {},
			expectedErr: "animal id is required",
		},
		{
			name: "Repository error - get fails",
			id:   "animal-1",
			setupMock: func(m *mocks.MockAnimalRepository) {
				m.On("Get", mock.Anything, "animal-1").Return(domain.Animal{}, errors.New("db error")).Once()
			},
			expectedErr: "db error",
		},
		{
			name: "Error - animal not found",
			id:   "animal-1",
			setupMock: func(m *mocks.MockAnimalRepository) {
				m.On("Get", mock.Anything, "animal-1").Return(domain.Animal{}, nil).Once()
			},
			expectedErr: "animal not found",
		},
		{
			name: "Repository error - delete fails",
			id:   "animal-1",
			setupMock: func(m *mocks.MockAnimalRepository) {
				m.On("Get", mock.Anything, "animal-1").Return(domain.Animal{ID: "animal-1"}, nil).Once()
				m.On("Delete", mock.Anything, "animal-1").Return(errors.New("delete failed")).Once()
			},
			expectedErr: "delete failed",
		},
		{
			name: "Success - animal deleted",
			id:   "animal-1",
			setupMock: func(m *mocks.MockAnimalRepository) {
				m.On("Get", mock.Anything, "animal-1").Return(domain.Animal{ID: "animal-1"}, nil).Once()
				m.On("Delete", mock.Anything, "animal-1").Return(nil).Once()
			},
			expectedErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockAnimalRepository)
			tt.setupMock(mockRepo)

			svc := service.NewAnimalService(mockRepo)
			err := svc.DeleteAnimal(context.Background(), tt.id)

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
