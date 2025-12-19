package analytics_service_api

import "fmt"

func eventTypeFromMetric(metric string) (string, error) {
	switch metric {
	case "views":
		return "view", nil
	case "likes":
		return "like", nil
	case "comments":
		return "comment", nil
	case "reposts":
		return "repost", nil
	default:
		return "", fmt.Errorf("unknown metric: %s", metric)
	}
}
