package pgstorage

import (
	"context"
	"fmt"
	"time"

	eventmodel "github.com/sashapremium/events/internal/models"
)

func (s *PGStorage) GetEventsByContentID(ctx context.Context, contentID string) ([]*eventmodel.ContentEvent, error) {
	rows, err := s.db.Query(ctx, `
		SELECT id, content_id, user_hash, type, comment, at
		FROM content_events
		WHERE content_id = $1
		ORDER BY at ASC, id ASC
	`, contentID)
	if err != nil {
		return nil, fmt.Errorf("select events: %w", err)
	}
	defer rows.Close()

	var result []*eventmodel.ContentEvent

	for rows.Next() {
		var (
			id       int64
			cid      string
			userHash string
			typ      string
			comment  *string
			at       time.Time
		)

		if err := rows.Scan(&id, &cid, &userHash, &typ, &comment, &at); err != nil {
			return nil, fmt.Errorf("scan event row: %w", err)
		}

		ev := &eventmodel.ContentEvent{
			ContentID: cid,
			UserHash:  userHash,
			Type:      eventmodel.EventType(typ),
			At:        at,
		}
		if comment != nil {
			ev.Comment = *comment
		}

		result = append(result, ev)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows err: %w", err)
	}

	return result, nil
}
