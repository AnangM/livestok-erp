package repository

import (
	"context"
	"database/sql"

	"github.com/AnangM/livestok-erp/internal/domain"
)

type AnimalRepository interface {
	Create(ctx context.Context, animal *domain.Animal) error
}

type PgAnimalRepository struct {
	db *sql.DB
}

func NewPostgresAnimalRepository(db *sql.DB) AnimalRepository {
	return &PgAnimalRepository{db: db}
}

func (r *PgAnimalRepository) Create(ctx context.Context, animal *domain.Animal) error {
	query := `
	INSERT INTO animals (farm_id, tag_number, breed, gender, sire_id, dam_id, birth_date, status)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING id, created_at
	`

	err := r.db.QueryRowContext(ctx,
		query,
		animal.FarmID,
		animal.TagNumber,
		animal.Breed,
		animal.Gender,
		animal.SireID,
		animal.DamID,
		animal.BirthDate,
		animal.Status,
	).Scan(&animal.ID, &animal.CreatedAt)

	if err != nil {
		return err
	}
	return nil
}
