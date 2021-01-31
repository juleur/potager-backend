package auth

import "github.com/jackc/pgx/v4/pgxpool"

type AuthPsql struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) *AuthPsql {
	return &AuthPsql{db: db}
}
