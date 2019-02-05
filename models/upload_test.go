package models_test

import (
	"testing"
	"time"

	"github.com/RTradeLtd/config"
	"github.com/RTradeLtd/database/models"
)

func TestMigration_Upload(t *testing.T) {
	cfg, err := config.LoadConfig(testCfgPath)
	if err != nil {
		t.Fatal(err)
	}
	db, err := openDatabaseConnection(t, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if check := db.AutoMigrate(&models.Upload{}); check.Error != nil {
		t.Fatal(err)
	}
}

func TestUpload(t *testing.T) {
	cfg, err := config.LoadConfig(testCfgPath)
	if err != nil {
		t.Fatal(err)
	}
	db, err := openDatabaseConnection(t, cfg)
	if err != nil {
		t.Fatal(err)
	}
	um := models.NewUploadManager(db)
	type args struct {
		hash       string
		uploadType string
		network    string
		holdTime   int64
		userName1  string
		userName2  string
		gcd        time.Time
		newGCD     time.Time
		encrypted  bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"User1-Hash1", args{"hash1", "file", "public", 5, "user1", "user2", time.Now(), time.Now().Add(time.Hour * 24), false}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			upload1, err := um.NewUpload(
				tt.args.hash,
				tt.args.uploadType,
				models.UploadOptions{
					NetworkName:      tt.args.network,
					Username:         tt.args.userName1,
					HoldTimeInMonths: tt.args.holdTime,
					Encrypted:        tt.args.encrypted,
				},
			)
			if err != nil {
				t.Fatal(err)
			}
			defer um.DB.Unscoped().Delete(upload1)
			upload2, err := um.NewUpload(
				tt.args.hash,
				tt.args.uploadType,
				models.UploadOptions{
					NetworkName:      tt.args.network,
					Username:         tt.args.userName2,
					HoldTimeInMonths: tt.args.holdTime,
					Encrypted:        tt.args.encrypted,
				},
			)
			if err != nil {
				t.Fatal(err)
			}
			defer um.DB.Unscoped().Delete(upload2)
			if _, err := um.NewUpload(
				tt.args.hash,
				tt.args.uploadType,
				models.UploadOptions{
					NetworkName:      tt.args.network,
					Username:         tt.args.userName2,
					HoldTimeInMonths: tt.args.holdTime,
					Encrypted:        tt.args.encrypted,
				},
			); err == nil {
				t.Fatal("expected error")
			} else if err.Error() != models.ErrAlreadyExistingUpload {
				t.Fatal("wrong error message received")
			}
			// test update which triggers shorter gcd error
			if _, err := um.UpdateUpload(1, tt.args.userName1, tt.args.hash, tt.args.network); err == nil {
				t.Fatal("expected error")
			} else if err.Error() != models.ErrShorterGCD {
				t.Fatal("wrong error returned")
			}
			// test update which passes
			if _, err := um.UpdateUpload(10, tt.args.userName1, tt.args.hash, tt.args.network); err != nil {
				t.Fatal(err)
			}
			// test finding uploads by network
			uploads, err := um.FindUploadsByNetwork(tt.args.network)
			if err != nil {
				t.Fatal(err)
			}
			var (
				user1Found bool
				user2Found bool
			)
			for _, upld := range uploads {
				if upld.UserName == tt.args.userName1 && upld.Hash == tt.args.hash {
					user1Found = true
				} else if upld.UserName == tt.args.userName2 && upld.Hash == tt.args.hash {
					user2Found = true
				}
			}
			if !user1Found || !user2Found {
				t.Fatal("failed to find uploads")
			}
			// test finding uploads by hash
			uploads, err = um.FindUploadsByHash(tt.args.hash)
			if err != nil {
				t.Fatal(err)
			}
			user1Found = false
			user2Found = false
			for _, upld := range uploads {
				if upld.UserName == tt.args.userName1 && upld.Hash == tt.args.hash {
					user1Found = true
				} else if upld.UserName == tt.args.userName2 && upld.Hash == tt.args.hash {
					user2Found = true
				}
			}
			if !user1Found || !user2Found {
				t.Fatal("failed to find uploads")
			}
			upload, err := um.FindUploadByHashAndUserAndNetwork(tt.args.userName1, tt.args.hash, tt.args.network)
			if err != nil {
				t.Fatal(err)
			}
			if upload.Hash != tt.args.hash {
				t.Fatal("failed to find correct hash")
			}
			uploads, err = um.GetUploadByHashForUser(tt.args.hash, tt.args.userName1)
			if err != nil {
				t.Fatal(err)
			}
			if uploads[0].Hash != tt.args.hash {
				t.Fatal("bad hash found")
			}
			user1Found = false
			user2Found = false
			uploads, err = um.GetUploads()
			if err != nil {
				t.Fatal(err)
			}
			for _, upld := range uploads {
				if upld.UserName == tt.args.userName1 && upld.Hash == tt.args.hash {
					user1Found = true
				} else if upld.UserName == tt.args.userName2 && upld.Hash == tt.args.hash {
					user2Found = true
				}
			}
			if !user1Found || !user2Found {
				t.Fatal("failed to find uploads")
			}
			uploads, err = um.GetUploadsForUser(tt.args.userName1)
			if err != nil {
				t.Fatal(err)
			}
			if uploads[0].Hash != tt.args.hash {
				t.Fatal("bad upload found")
			}
		})
	}
}
