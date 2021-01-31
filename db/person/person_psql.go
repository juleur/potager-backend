package person

import (
	"github.com/jackc/pgx/v4/pgxpool"
)

type PersonPsql struct {
	db *pgxpool.Pool
}

func NewPersonRepository(db *pgxpool.Pool) *PersonPsql {
	return &PersonPsql{db: db}
}
