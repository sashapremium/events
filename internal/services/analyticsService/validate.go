package analyticsService

func (s *Service) validateMetric(metric string) error {
	switch metric {
	case "views", "likes", "comments", "reposts":
		return nil
	default:
		return ErrInvalidMetric
	}
}
