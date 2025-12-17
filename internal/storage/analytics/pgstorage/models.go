package pgstorage

import "time"

const (
	eventsTableName = "content_events"
	totalsTableName = "post_totals"

	IDColumnName        = "id"
	ContentIDColumnName = "content_id" // post_id
	AuthorIDColumnName  = "author_id"
	UserHashColumnName  = "user_hash"
	TypeColumnName      = "type"
	CommentColumnName   = "comment"
	AtColumnName        = "at"

	PostIDColumnName      = "post_id"
	ViewsColumnName       = "views"
	LikesColumnName       = "likes"
	CommentsColumnName    = "comments"
	RepostsColumnName     = "reposts"
	UniqueUsersColumnName = "unique_users"
	UpdatedAtColumnName   = "updated_at"
)

type ContentEvent struct {
	ContentID string
	AuthorID  string
	UserHash  string
	Type      string
	Comment   string
	At        time.Time
}

type TopPostItem struct {
	ContentID string
	Count     uint64
}

type TopAuthorItem struct {
	AuthorID string
	Count    uint64
}

type DailyCount struct {
	Day   time.Time
	Count uint64
}

type PostStats struct {
	ContentID    string
	ViewsTotal   uint64
	ViewsUnique  uint64
	Likes        uint64
	Reposts      uint64
	Comments     uint64
	DailyViews   []DailyCount
	CalculatedAt time.Time
}

type AuthorStats struct {
	AuthorID     string
	ViewsTotal   uint64
	ViewsUnique  uint64
	Likes        uint64
	Reposts      uint64
	Comments     uint64
	CalculatedAt time.Time
}
