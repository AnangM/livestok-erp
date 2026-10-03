package repository

import (
	"context"
	"database/sql"

	"github.com/AnangM/livestok-erp/internal/domain"
)

type HarvestRepository interface {
	Create(ctx context.Context, harvest *domain.Harvest) error
	List(ctx context.Context, farmId string) ([]domain.Harvest, error)
	Update(ctx context.Context, id string, harvest *domain.Harvest) (domain.Harvest, error)
	Get(ctx context.Context, id string) (domain.Harvest, error)
	Delete(ctx context.Context, id string) error
}

type PgHarvestRepository struct {
	db *sql.DB
}

// Create implements [HarvestRepository].
func (p *PgHarvestRepository) Create(ctx context.Context, harvest *domain.Harvest) error {
	query := `
		INSERT INTO harvests (animal_id, farm_id, harvest_date, weight_kg)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	err := p.db.QueryRowContext(ctx,
		query,
		harvest.AnimalID,
		harvest.FarmID,
		harvest.HarvestDate,
		harvest.Weight,
	).Scan(&harvest.ID, &harvest.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

// Delete implements [HarvestRepository].
func (p *PgHarvestRepository) Delete(ctx context.Context, id string) error {
	query := `
		DELETE FROM harvests
		WHERE id = $1
	`
	_, err := p.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}

// Get implements [HarvestRepository].
func (p *PgHarvestRepository) Get(ctx context.Context, id string) (domain.Harvest, error) {
	query := `
		SELECT 
			id, animal_id, farm_id, harvest_date, weight_kg, created_at
		FROM harvests
		WHERE id = $1
	`
	var harvest domain.Harvest
	err := p.db.QueryRowContext(ctx, query, id).Scan(&harvest.ID, &harvest.AnimalID, &harvest.FarmID, &harvest.HarvestDate, &harvest.Weight, &harvest.CreatedAt)
	if err != nil {
		return domain.Harvest{}, err
	}
	return harvest, nil
}

// List implements [HarvestRepository].
func (p *PgHarvestRepository) List(ctx context.Context, farmId string) ([]domain.Harvest, error) {
	query := `
		SELECT
		id, animal_id, farm_id, harvest_date, weight_kg, created_at
		FROM harvests
		WHERE farm_id = $1
		ORDER BY harvest_date DESC
	`
	rows, err := p.db.QueryContext(ctx, query, farmId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var harvests []domain.Harvest
	for rows.Next() {
		var h domain.Harvest
		if err := rows.Scan(&h.ID, &h.AnimalID, &h.FarmID, &h.HarvestDate, &h.Weight, &h.CreatedAt); err != nil {
			return nil, err
		}
		harvests = append(harvests, h)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return harvests, nil
}

// Update implements [HarvestRepository].
func (p *PgHarvestRepository) Update(ctx context.Context, id string, harvest *domain.Harvest) (domain.Harvest, error) {
	query := `
		UPDATE harvests
		SET animal_id = $2, farm_id = $3, harvest_date = $4, weight_kg = $5
		WHERE id = $1
		RETURNING id, animal_id, farm_id, harvest_date, weight_kg, created_at
	`
	var updatedHarvest domain.Harvest
	err := p.db.QueryRowContext(ctx, query, id, harvest.AnimalID, harvest.FarmID, harvest.HarvestDate, harvest.Weight).Scan(&updatedHarvest.ID, &updatedHarvest.AnimalID, &updatedHarvest.FarmID, &updatedHarvest.HarvestDate, &updatedHarvest.Weight, &updatedHarvest.CreatedAt)
	if err != nil {
		return domain.Harvest{}, err
	}
	return updatedHarvest, nil
}

func NewPostgresHarvestRepository(db *sql.DB) HarvestRepository {
	return &PgHarvestRepository{db: db}
}
