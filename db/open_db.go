package db

import (
	"context"
	"fmt"
	"log"
	"npp_backend/config"

	"github.com/jackc/pgx/v4/pgxpool"
)

func OpenDB() *pgxpool.Pool {
	connString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", config.DB_HOST, config.DB_USER, config.DB_PASSWORD, config.DB_DATABASE)
	conn, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		log.Fatalln(err)
	}
	return conn
}
