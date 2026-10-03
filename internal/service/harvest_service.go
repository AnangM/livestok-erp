package service

import (
	"context"
	"errors"
	"time"

	"github.com/AnangM/livestok-erp/internal/domain"
	"github.com/AnangM/livestok-erp/internal/repository"
)

type HarvestServiceInterface interface {
	RegisterNewHarvest(ctx context.Context, harvest *domain.Harvest) error
	ListHarvests(ctx context.Context, farmId string) ([]domain.Harvest, error)
	UpdateHarvest(ctx context.Context, id string, harvest *domain.Harvest) (domain.Harvest, error)
	GetHarvest(ctx context.Context, id string) (domain.Harvest, error)
	DeleteHarvest(ctx context.Context, id string) error
}

type HarvestService struct {
	repo repository.HarvestRepository
}

// DeleteHarvest implements [HarvestServiceInterface].
func (s *HarvestService) DeleteHarvest(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("harvest_id is required")
	}

	h, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	if h.ID == "" {
		return errors.New("harvest not found")
	}
	return s.repo.Delete(ctx, id)
}

// GetHarvest implements [HarvestServiceInterface].
func (s *HarvestService) GetHarvest(ctx context.Context, id string) (domain.Harvest, error) {
	if id == "" {
		return domain.Harvest{}, errors.New("harvest_id is required")
	}
	return s.repo.Get(ctx, id)
}

// ListHarvests implements [HarvestServiceInterface].
func (s *HarvestService) ListHarvests(ctx context.Context, farmId string) ([]domain.Harvest, error) {
	return s.repo.List(ctx, farmId)
}

// UpdateHarvest implements [HarvestServiceInterface].
func (s *HarvestService) UpdateHarvest(ctx context.Context, id string, h *domain.Harvest) (domain.Harvest, error) {
	if h.FarmID == "" {
		return domain.Harvest{}, errors.New("farm_id is required")
	}
	if h.AnimalID == "" {
		return domain.Harvest{}, errors.New("animal_id is required")
	}
	if h.Weight == 0 {
		return domain.Harvest{}, errors.New("weight is required")
	}
	if h.HarvestDate.After(time.Now()) {
		return domain.Harvest{}, errors.New("harvest date cannot be in the future")
	}

	_, err := s.repo.Get(ctx, id)
	if err != nil {
		return domain.Harvest{}, err
	}
	if h.ID == "" {
		return domain.Harvest{}, errors.New("harvest not found")
	}
	h.FarmID = h.FarmID
	h.AnimalID = h.AnimalID
	h.Weight = h.Weight
	h.HarvestDate = h.HarvestDate
	return s.repo.Update(ctx, id, h)
}

func NewHarvestService(repo repository.HarvestRepository) HarvestServiceInterface {
	return &HarvestService{repo: repo}
}

func (s *HarvestService) RegisterNewHarvest(ctx context.Context, h *domain.Harvest) error {
	if h.FarmID == "" {
		return errors.New("farm_id is required")
	}

	if h.AnimalID == "" {
		return errors.New("animal_id is required")
	}

	if h.Weight == 0 {
		return errors.New("weight is required")
	}

	if h.HarvestDate.After(time.Now()) {
		return errors.New("harvest date cannot be in the future")
	}

	return s.repo.Create(ctx, h)
}
