package analytics

import (
	"fmt"
	"log"

	"github.com/sashapremium/events/config"
	analyticsPG "github.com/sashapremium/events/internal/storage/analytics/pgstorage"
)

func InitPGStorage(cfg *config.Config) *analyticsPG.PGStorage {
	sslmode := cfg.Postgres.SSLMode
	if sslmode == "" {
		sslmode = "disable"
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.DBName,
		sslmode,
	)

	storage, err := analyticsPG.NewPGStorage(dsn)
	if err != nil {
		log.Fatalf("ошибка инициализации БД analytics: %v", err)
	}
	return storage
}
