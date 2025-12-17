package pgstorage

import (
	"context"
	"fmt"
	"time"

	eventmodel "github.com/sashapremium/events/internal/models"
)

func (s *PGStorage) InsertEvents(ctx context.Context, events []*eventmodel.ContentEvent) error {
	if len(events) == 0 {
		return nil
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	const q = `
		INSERT INTO content_events (content_id, user_hash, type, comment, at)
		VALUES ($1, $2, $3, $4, $5)
	`

	for _, e := range events {
		if e == nil {
			continue
		}

		at := e.At
		if at.IsZero() {
			at = time.Now().UTC()
		}

		var comment *string
		if e.Comment != "" {
			c := e.Comment
			comment = &c
		}

		if _, err := tx.Exec(ctx, q, e.ContentID, e.UserHash, string(e.Type), comment, at); err != nil {
			return fmt.Errorf("exec insert: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}
