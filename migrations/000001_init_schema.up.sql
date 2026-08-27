-- Enable UUID extension if not already active
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Core Animals Table
CREATE TABLE animals (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    farm_id UUID NOT NULL, -- Maps to Supabase auth.users.id
    tag_number VARCHAR(50) NOT NULL,
    breed VARCHAR(50), -- e.g., 'New Zealand White' 
    gender VARCHAR(10) CHECK (gender IN ('Buck', 'Doe', 'Unknown')),
    sire_id UUID REFERENCES animals(id) ON DELETE SET NULL, -- Father
    dam_id UUID REFERENCES animals(id) ON DELETE SET NULL,  -- Mother
    birth_date DATE NOT NULL,
    status VARCHAR(20) DEFAULT 'Active' CHECK (status IN ('Active', 'Harvested', 'Sold', 'Deceased')),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    
    -- Ensure tag numbers are unique per farm
    UNIQUE(farm_id, tag_number) 
);

-- Indexing foreign keys is critical for recursive lineage queries
CREATE INDEX idx_animals_sire ON animals(sire_id);
CREATE INDEX idx_animals_dam ON animals(dam_id);
CREATE INDEX idx_animals_farm ON animals(farm_id);

-- Production Tracking (Meat Yields)
CREATE TABLE harvests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    animal_id UUID NOT NULL REFERENCES animals(id) ON DELETE CASCADE,
    farm_id UUID NOT NULL,
    harvest_date DATE NOT NULL,
    weight_kg DECIMAL(5,2) NOT NULL, 
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_harvests_farm_date ON harvests(farm_id, harvest_date);