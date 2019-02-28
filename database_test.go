package database

import (
	"testing"

	"github.com/RTradeLtd/config"
	"go.uber.org/zap"
)

func TestNew(t *testing.T) {
	t.Run("invalid argument", func(t *testing.T) {
		if _, err := New(nil, Options{}); err == nil {
			t.Error("expected error")
		}
	})
	t.Run("with migrations and logger", func(t *testing.T) {
		db, err := New(&config.TemporalConfig{
			Database: config.Database{
				Name:     "temporal",
				URL:      "127.0.0.1",
				Port:     "5433",
				Username: "postgres",
				Password: "password123",
			},
		}, Options{
			RunMigrations:  true,
			SSLModeDisable: true,
			Logger:         NewZapLogger(LogLevelInfo, zap.NewExample().Sugar())})
		if err != nil {
			t.Fatal(err)
		}
		db.Close()
	})
}
