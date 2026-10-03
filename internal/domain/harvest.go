package domain

import "time"

type Harvest struct {
	ID          string    `json:"id"`
	AnimalID    string    `json:"animal_id"`
	FarmID      string    `json:"farm_id"`
	HarvestDate time.Time `json:"harvest_date"`
	Weight      float64   `json:"weight"`
	CreatedAt   time.Time `json:"created_at"`
}
