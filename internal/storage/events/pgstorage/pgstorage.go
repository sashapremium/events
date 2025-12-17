package pgstorage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PGStorage struct {
	db *pgxpool.Pool
}

func NewPGStorge(connString string) (*PGStorage, error) {
	cfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга конфига: %w", err)
	}

	db, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения: %w", err)
	}

	s := &PGStorage{db: db}

	if err := s.initTables(context.Background()); err != nil {
		db.Close()
		return nil, err
	}

	return s, nil
}

func (s *PGStorage) Close() {
	if s.db != nil {
		s.db.Close()
	}
}

func (s *PGStorage) initTables(ctx context.Context) error {
	sql := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			%s      BIGSERIAL PRIMARY KEY,
			%s      TEXT NOT NULL,
			%s      TEXT NOT NULL,
			%s      TEXT NOT NULL,
			%s      TEXT,
			%s      TIMESTAMPTZ NOT NULL DEFAULT now()
		);

		CREATE INDEX IF NOT EXISTS idx_%s_content_id_at
			ON %s (%s, %s);

		CREATE INDEX IF NOT EXISTS idx_%s_type_at
			ON %s (%s, %s);
	`,
		eventsTableName,
		IDColumnName,
		ContentIDColumnName,
		UserHashColumnName,
		TypeColumnName,
		CommentColumnName,
		AtColumnName,

		eventsTableName, eventsTableName, ContentIDColumnName, AtColumnName,
		eventsTableName, eventsTableName, TypeColumnName, AtColumnName,
	)

	if _, err := s.db.Exec(ctx, sql); err != nil {
		return fmt.Errorf("init tables: %w", err)
	}
	return nil
}
