package models

import "time"

type EventType string

// Информация о посте
type PostInfo struct {
	ID          uint64
	Title       string
	AuthorID    string
	Category    string
	Content     string
	PublishedAt time.Time
}
