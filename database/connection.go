package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
	"os"
	"time"

	"backend/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

func NewConnection(cfg *config.Config) (*DB, error) {
	config, err := pgxpool.ParseConfig(cfg.DatabaseURL())
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	// Configure connection pool
	config.MaxConns = 30
	config.MinConns = 5

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test the connection
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to database")

	if err := runInitSQL(context.Background(), pool, "./scripts/init.sql"); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to run init.sql: %w", err)
	}

	return &DB{Pool: pool}, nil
}

func runInitSQL(ctx context.Context, pool *pgxpool.Pool, path string) error {
	// читаем файл
	b, err := os.ReadFile(path)
	if err != nil {
		// если файла нет — просто выходим без ошибки
		if os.IsNotExist(err) {
			fmt.Println(1123123)
			return nil
		}
		return fmt.Errorf("read init.sql: %w", err)
	}
	sqlText := string(b)

	// таймаут на инициализацию
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	// глобальная блокировка, чтобы несколько инстансов не запускали init одновременно
	if _, err := tx.Exec(ctx, `SELECT pg_advisory_lock(7242025)`); err != nil {
		return fmt.Errorf("advisory_lock: %w", err)
	}
	defer tx.Exec(ctx, `SELECT pg_advisory_unlock(7242025)`)

	// маркер инициализации
	if _, err := tx.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS app_bootstrap (
		  id          int PRIMARY KEY DEFAULT 1,
		  runned_at   timestamptz NOT NULL DEFAULT now()
		);
	`); err != nil {
		return fmt.Errorf("create marker: %w", err)
	}

	var already bool
	if err := tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM app_bootstrap WHERE id=1)`).Scan(&already); err != nil {
		return fmt.Errorf("check marker: %w", err)
	}
	if already {
		return tx.Commit(ctx)
	}

	// ВАЖНО: благодаря Simple Protocol (выше) можно исполнить несколько стейтментов за раз
	if _, err := tx.Exec(ctx, sqlText); err != nil {
		return fmt.Errorf("exec init.sql: %w", err)
	}

	if _, err := tx.Exec(ctx, `INSERT INTO app_bootstrap (id) VALUES (1) ON CONFLICT (id) DO NOTHING`); err != nil {
		return fmt.Errorf("insert marker: %w", err)
	}

	return tx.Commit(ctx)
}

func (db *DB) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}

func (db *DB) Health() error {
	return db.Pool.Ping(context.Background())
}
