package database

import (
	"testing"

	"github.com/RTradeLtd/config/v2"
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
				URL:      "localhost",
				Port:     "5432",
				Username: "temporaladmin",
				Password: "temporaladmin",
			},
		}, Options{
			RunMigrations: true,
			Logger:        NewZapLogger(LogLevelInfo, zap.NewExample().Sugar())})
		if err != nil {
			t.Fatal(err)
		}
		db.Close()
	})
}
