package pgstorage

import "time"

type EventRow struct {
	ID        int64
	ContentID string
	UserHash  string
	Type      string
	Comment   *string
	At        time.Time
}

const (
	eventsTableName     = "content_events"
	IDColumnName        = "id"
	ContentIDColumnName = "content_id"
	UserHashColumnName  = "user_hash"
	TypeColumnName      = "type"
	CommentColumnName   = "comment"
	AtColumnName        = "at"
)
