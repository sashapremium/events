package models

import "time"

type EventType string

const (
	EventView    EventType = "view"
	EventLike    EventType = "like"
	EventRepost  EventType = "repost"
	EventComment EventType = "comment"
)

// Событие по посту
type ContentEvent struct {
	ContentID string    `json:"content_id"`
	AuthorID  uint64    `json:"author_id"`
	UserHash  string    `json:"user_hash"`
	Type      EventType `json:"type"`
	Comment   string    `json:"comment,omitempty"`
	At        time.Time `json:"at"`
}
