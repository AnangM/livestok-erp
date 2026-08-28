package repository

import (
	"context"
	"database/sql"

	"github.com/AnangM/livestok-erp/internal/domain"
)

type AnimalRepository interface {
	Create(ctx context.Context, animal *domain.Animal) error
	List(ctx context.Context, farmId string) ([]domain.Animal, error)
	Update(ctx context.Context, id string, animal *domain.Animal) (domain.Animal, error)
	Get(ctx context.Context, id string) (domain.Animal, error)
	Delete(ctx context.Context, id string) error
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

func (r *PgAnimalRepository) List(ctx context.Context, farmId string) ([]domain.Animal, error) {
	query := `
	SELECT
	id, tag_number, breed, gender, sire_id, dam_id, birth_date, status, created_at
	FROM animals
	WHERE farm_id = $1
	ORDER BY created_at DESC;
	`

	rows, err := r.db.QueryContext(ctx, query, farmId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var animals []domain.Animal

	for rows.Next() {
		var a domain.Animal
		if err := rows.Scan(&a.ID, &a.TagNumber, &a.Breed, &a.Gender, &a.SireID, &a.DamID, &a.BirthDate, &a.Status, &a.CreatedAt); err != nil {
			return nil, err
		}
		animals = append(animals, a)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return animals, nil
}

func (r *PgAnimalRepository) Update(ctx context.Context, id string, a *domain.Animal) (domain.Animal, error) {
	query := `
	UPDATE animals
	SET tag_number = $2, breed = $3, gender = $4, sire_id = $5, dam_id = $6, birth_date = $7, status = $8
	WHERE id = $1
	RETURNING id, farm_id, tag_number, breed, gender, sire_id, dam_id, birth_date, status, created_at
	`
	var updatedAnimal domain.Animal
	err := r.db.QueryRowContext(ctx, query, id, a.TagNumber, a.Breed, a.Gender, a.SireID, a.DamID, a.BirthDate, a.Status).Scan(&updatedAnimal.ID, &updatedAnimal.FarmID, &updatedAnimal.TagNumber, &updatedAnimal.Breed, &updatedAnimal.Gender, &updatedAnimal.SireID, &updatedAnimal.DamID, &updatedAnimal.BirthDate, &updatedAnimal.Status, &updatedAnimal.CreatedAt)
	if err != nil {
		return domain.Animal{}, err
	}
	return updatedAnimal, nil
}

func (r *PgAnimalRepository) Get(ctx context.Context, id string) (domain.Animal, error) {
	query := `
	SELECT
	id, tag_number, breed, gender, sire_id, dam_id, birth_date, status, created_at
	FROM animals
	WHERE id = $1
	`
	var animal domain.Animal
	err := r.db.QueryRowContext(ctx, query, id).Scan(&animal.ID, &animal.TagNumber, &animal.Breed, &animal.Gender, &animal.SireID, &animal.DamID, &animal.BirthDate, &animal.Status, &animal.CreatedAt)
	if err != nil {
		return domain.Animal{}, err
	}
	return animal, nil
}

func (r *PgAnimalRepository) Delete(ctx context.Context, id string) error {
	query := `
	DELETE FROM animals
	WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
