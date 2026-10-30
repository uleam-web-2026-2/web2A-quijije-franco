package db

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

//go:embed migrations/001_init.sql
var schema embed.FS

func OpenAndMigrate(ctx context.Context, url string) (*sql.DB, error) {
	if url == "" {
		return nil, fmt.Errorf("DATABASE_URL no configurada")
	}
	cfg, err := pgx.ParseConfig(url)
	if err != nil {
		return nil, err
	}
	db := sql.OpenDB(stdlib.GetConnector(*cfg))
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}
	b, err := schema.ReadFile("migrations/001_init.sql")
	if err != nil {
		db.Close()
		return nil, err
	}
	if _, err = db.ExecContext(ctx, string(b)); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
