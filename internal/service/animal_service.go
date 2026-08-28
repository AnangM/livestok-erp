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
	ListAnimals(ctx context.Context, farmId string) ([]domain.Animal, error)
	UpdateAnimal(ctx context.Context, id string, animal *domain.Animal) (domain.Animal, error)
	GetAnimal(ctx context.Context, id string) (domain.Animal, error)
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

func (s *AnimalService) ListAnimals(ctx context.Context, farmId string) ([]domain.Animal, error) {
	return s.repo.List(ctx, farmId)
}

func (s *AnimalService) UpdateAnimal(ctx context.Context, id string, a *domain.Animal) (domain.Animal, error) {

	a.TagNumber = strings.TrimSpace(a.TagNumber)
	if a.TagNumber == "" {
		return domain.Animal{}, errors.New("tag number is required")
	}

	if a.FarmID == "" {
		return domain.Animal{}, errors.New("farm_id is required")
	}

	if a.Status == "" {
		a.Status = "Active"
	}

	if a.BirthDate.After(time.Now()) {
		return domain.Animal{}, errors.New("birth date cannot be in the future")
	}

	animal, err := s.repo.Get(ctx, id)

	if err != nil {
		return domain.Animal{}, err
	}

	if animal.ID == "" {
		return domain.Animal{}, errors.New("animal not found")
	}

	animal.TagNumber = a.TagNumber
	animal.Breed = a.Breed
	animal.Gender = a.Gender
	animal.SireID = a.SireID
	animal.DamID = a.DamID
	animal.BirthDate = a.BirthDate
	animal.Status = a.Status

	return s.repo.Update(ctx, id, &animal)
}

func (s *AnimalService) GetAnimal(ctx context.Context, id string) (domain.Animal, error) {

	if id == "" {
		return domain.Animal{}, errors.New("animal id is required")
	}

	return s.repo.Get(ctx, id)
}
