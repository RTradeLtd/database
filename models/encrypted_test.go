package models_test

import (
	"testing"

	"github.com/RTradeLtd/config"
	"github.com/RTradeLtd/database/models"
)

func TestMigration_EncryptedUpload(t *testing.T) {
	cfg, err := config.LoadConfig(testCfgPath)
	if err != nil {
		t.Fatal(err)
	}
	db, err := openDatabaseConnection(t, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if check := db.AutoMigrate(&models.EncryptedUpload{}); check.Error != nil {
		t.Fatal(err)
	}
}

func TestEncryptedUploads(t *testing.T) {
	cfg, err := config.LoadConfig(testCfgPath)
	if err != nil {
		t.Fatal(err)
	}
	db, err := openDatabaseConnection(t, cfg)
	if err != nil {
		t.Fatal(err)
	}
	ecm := models.NewEncryptedUploadManager(db)
	type args struct {
		user    string
		file    string
		network string
		hash    string
	}
	tests := []struct {
		name        string
		args        args
		wantUploads bool
		wantErr     bool
	}{
		{"Success", args{"user", "suchfilemuchspaceverydisk", "public", "dathashdoe"}, true, false},
		{"Failure", args{"notarealuser", "notarealfile", "notarealnetwork", "notarealhash"}, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Success" {
				if _, err := ecm.NewUpload(
					tt.args.user, tt.args.file, tt.args.network, tt.args.hash,
				); (err != nil) != tt.wantErr {
					t.Fatalf("NewUpload err = %v, wantErr %v", err, tt.wantErr)
				}
			}
			uploads, err := ecm.FindUploadsByUser(tt.args.user)
			if err != nil {
				t.Fatal(err)
			}
			if (len(*uploads) != 0) != tt.wantUploads {
				t.Fatalf("FineUploadsByUser uploads = %v, wantUploads %v", len(*uploads) != 0, tt.wantUploads)
			}
		})
	}
}
