package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/AnangM/livestok-erp/internal/domain"
	"github.com/AnangM/livestok-erp/internal/repository"
)

type AnimalServiceInterface interface {
	RegisterNewAnimal(ctx context.Context, animal *domain.Animal) error
}

type AnimalService struct {
	repo repository.AnimalRepository
}

func NewAnimalService(repo repository.AnimalRepository) AnimalServiceInterface {
	return &AnimalService{repo: repo}
}

func (s *AnimalService) RegisterNewAnimal(ctx context.Context, a *domain.Animal) error {
	a.TagNumber = strings.TrimSpace(a.TagNumber)
	if a.TagNumber == "" {
		return errors.New("tag number is required")
	}

	if a.FarmID == "" {
		return errors.New("farm_id is required")
	}

	if a.Status == "" {
		a.Status = "Active"
	}

	if a.BirthDate.After(time.Now()) {
		return errors.New("birth date cannot be in the future")
	}

	return s.repo.Create(ctx, a)
}
