package models

import (
	"fmt"
	"testing"

	"github.com/RTradeLtd/config/v2"
	"github.com/jinzhu/gorm"
)

type testLogger struct{ t *testing.T }

func (t *testLogger) Print(args ...interface{}) { t.t.Log(args...) }

func newTestDB(t *testing.T, model interface{}) *gorm.DB {
	cfg, err := config.LoadConfig("../testenv/config.json")
	if err != nil {
		t.Fatal(err)
	}

	db, err := gorm.Open("postgres",
		fmt.Sprintf("host=%s port=%s user=postgres dbname=temporal password=%s sslmode=disable",
			cfg.Database.URL, cfg.Database.Port, cfg.Database.Password))
	if err != nil {
		t.Fatal(err)
	}
	db.SetLogger(&testLogger{t})
	db.LogMode(true)

	if model != nil {
		if check := db.AutoMigrate(model); check.Error != nil {
			t.Fatalf("could not execute migration for model '%+v': %s",
				model, err.Error())
		}
	}

	return db
}

func TestAutoMigrate(t *testing.T) {
	type args struct {
		model interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{"encrypted upload", args{&EncryptedUpload{}}},
		{"ipfs networks", args{&HostedNetwork{}}},
		{"ipns", args{&IPNS{}}},
		{"payment", args{&Payments{}}},
		{"record", args{&Record{}}},
		{"tns zone", args{&Zone{}}},
		{"upload", args{&Upload{}}},
		{"usage", args{&Usage{}}},
		{"user", args{&User{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newTestDB(t, tt.args.model)
			db.Close()
		})
	}
}
