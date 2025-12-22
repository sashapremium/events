package pgstorage

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/sashapremium/events/analytics/internal/models"
)

func (s *PGStorage) InsertEvents(ctx context.Context, events []*models.ContentEvent) error {
	if len(events) == 0 {
		return nil
	}

	const q = `
        INSERT INTO content_events (content_id, author_id, user_hash, type, comment, at)
        VALUES ($1, $2, $3, $4, $5, $6)
    `

	batch := &pgx.Batch{}
	for _, ev := range events {
		if ev == nil {
			continue
		}
		batch.Queue(q, ev.ContentID, ev.AuthorID, ev.UserHash, string(ev.Type), ev.Comment, ev.At)
	}

	br := s.db.SendBatch(ctx, batch)
	defer br.Close()

	for i := 0; i < batch.Len(); i++ {
		if _, err := br.Exec(); err != nil {
			return err
		}
	}

	return nil
}
