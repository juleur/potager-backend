package farmer

import (
	"github.com/jackc/pgx/v4/pgxpool"
)

type FarmerPsql struct {
	db *pgxpool.Pool
}

func NewFarmerRepository(db *pgxpool.Pool) *FarmerPsql {
	return &FarmerPsql{db: db}
}
