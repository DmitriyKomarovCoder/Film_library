package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Postgres struct {
	maxPoolSize int32
	Pool        *pgxpool.Pool
}

func New(host, userName, password, dbName string, port int, poolSize int32) (*Postgres, error) {
	pg := &Postgres{
		maxPoolSize: poolSize,
	}

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, userName, password, dbName)

	poolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - pgxpool.ParseConfig: %w", err)
	}

	poolConfig.MaxConns = pg.maxPoolSize

	pg.Pool, err = pgxpool.ConnectConfig(context.Background(), poolConfig)

	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - connAttempts == 0: %w", err)
	}

	return pg, nil
}

func (p *Postgres) Close(ctx context.Context) error {
	if p.Pool != nil {
		p.Pool.Close()
	}
	return nil
}
