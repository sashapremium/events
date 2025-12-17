package analyticsService

import "time"

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
