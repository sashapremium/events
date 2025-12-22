package analyticsService

import (
	"errors"
	"time"
)

func (s *Service) SetSyncInterval(d time.Duration) {
	if d > 0 {
		s.syncInterval = d
	}
}

func (s *Service) SetSyncBatch(n int) {
	if n > 0 {
		s.syncBatch = n
	}
}

type Service struct {
	storage Storage
	cache   Cache

	syncInterval time.Duration
	syncBatch    int
}

func New(storage Storage, cache Cache) *Service {
	return &Service{
		storage:      storage,
		cache:        cache,
		syncInterval: 3 * time.Second,
		syncBatch:    200,
	}
}

type TotalsDelta struct {
	Views, Likes, Comments, Reposts, UniqueUsers int64
}

var ErrInvalidMetric = errors.New("некорректная метрика (ожидается views|likes|comments|reposts)")
