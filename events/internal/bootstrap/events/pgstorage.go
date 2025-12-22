package events

import (
	"fmt"
	"log"

	"github.com/sashapremium/events/config"
	"github.com/sashapremium/events/internal/storage/events/pgstorage"
)

func InitPGStorage(cfg *config.Config) *pgstorage.PGStorage {
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

	storage, err := pgstorage.NewPGStorge(dsn)
	if err != nil {
		log.Fatalf("ошибка инициализации БД: %v", err)
	}
	return storage
}
