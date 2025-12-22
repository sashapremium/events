package pgstorage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	analyticsService "github.com/sashapremium/events/analytics/internal/services/analyticsService"
)

type PGStorage struct {
	db *pgxpool.Pool
}

func (s *PGStorage) UpsertPostTotals(ctx context.Context, postID uint64, delta analyticsService.TotalsDelta) error {
	//TODO implement me
	panic("implement me")
}

func NewPGStorage(connString string) (*PGStorage, error) {
	cfg, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("parse pg config: %w", err)
	}

	db, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		return nil, fmt.Errorf("connect postgres: %w", err)
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
			%s BIGINT PRIMARY KEY,
			%s BIGINT NOT NULL DEFAULT 0,
			%s BIGINT NOT NULL DEFAULT 0,
			%s BIGINT NOT NULL DEFAULT 0,
			%s BIGINT NOT NULL DEFAULT 0,
			%s BIGINT NOT NULL DEFAULT 0,
			%s TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);

		CREATE TABLE IF NOT EXISTS %s (
			%s BIGSERIAL PRIMARY KEY,
			%s BIGINT NOT NULL,
			%s TEXT,
			%s TEXT NOT NULL,
			%s TEXT NOT NULL,
			%s TEXT,
			%s TIMESTAMPTZ NOT NULL
		);

		CREATE INDEX IF NOT EXISTS idx_%s_content_at
			ON %s (%s, %s);

		CREATE INDEX IF NOT EXISTS idx_%s_author_type_at
			ON %s (%s, %s, %s);
	`,
		totalsTableName,
		PostIDColumnName,
		ViewsColumnName,
		LikesColumnName,
		CommentsColumnName,
		RepostsColumnName,
		UniqueUsersColumnName,
		UpdatedAtColumnName,

		eventsTableName,
		IDColumnName,
		ContentIDColumnName,
		AuthorIDColumnName,
		UserHashColumnName,
		TypeColumnName,
		CommentColumnName,
		AtColumnName,

		eventsTableName, eventsTableName, ContentIDColumnName, AtColumnName,
		eventsTableName, eventsTableName, AuthorIDColumnName, TypeColumnName, AtColumnName,
	)

	if _, err := s.db.Exec(ctx, sql); err != nil {
		return fmt.Errorf("init tables: %w", err)
	}
	return nil
}
