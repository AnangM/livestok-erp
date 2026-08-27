package domain

import (
	"time"
)

type Animal struct {
	ID        string    `json:"id"`
	FarmID    string    `json:"farm_id"` // Comes from the JWT context
	TagNumber string    `json:"tag_number"`
	Breed     string    `json:"breed"`
	Gender    string    `json:"gender"`
	SireID    *string   `json:"sire_id,omitempty"` // Pointer for nullable fields
	DamID     *string   `json:"dam_id,omitempty"`
	BirthDate time.Time `json:"birth_date"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
