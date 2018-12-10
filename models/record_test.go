package models_test

import (
	"testing"

	"github.com/RTradeLtd/config"
	"github.com/RTradeLtd/database/models"
)

func TestMigration_Record(t *testing.T) {
	cfg, err := config.LoadConfig(testCfgPath)
	if err != nil {
		t.Fatal(err)
	}
	db, err := openDatabaseConnection(t, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if err = db.AutoMigrate(&models.Record{}).Error; err != nil {
		t.Fatal(err)
	}
}

func TestRecord(t *testing.T) {
	cfg, err := config.LoadConfig(testCfgPath)
	if err != nil {
		t.Fatal(err)
	}
	db, err := openDatabaseConnection(t, cfg)
	if err != nil {
		t.Fatal(err)
	}
	rm := models.NewRecordManager(db)
	type args struct {
		recordName    string
		recordKeyName string
		zoneName      string
		ipfsHash      string
		metadata      map[string]interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{"NoMetaData", args{"testrecord1", "testkey1", "testzone1", "testhash1", nil}},
		{"YesMetaData", args{"testrecord2", "testkey2", "testzone2", "testhash2", map[string]interface{}{
			"food": "pizza",
			"pet":  "dog",
		}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			record1, err := rm.AddRecord(
				tt.args.recordName,
				tt.args.recordKeyName,
				tt.args.zoneName,
				tt.args.metadata,
			)
			if err != nil {
				t.Fatal(err)
			}
			defer db.Delete(record1)
			if record1.LatestIPFSHash != "" {
				t.Fatal("latest ipfs hash should be empty")
			}
			record2, err := rm.UpdateLatestIPFSHash(
				tt.args.recordName,
				tt.args.ipfsHash,
			)
			if err != nil {
				t.Fatal(err)
			}
			if record2.LatestIPFSHash != tt.args.ipfsHash {
				t.Fatal("bad ipfs hash set")
			}
			record3, err := rm.FindRecordByName(tt.args.recordName)
			if err != nil {
				t.Fatal(err)
			}
			if record3.LatestIPFSHash != tt.args.ipfsHash {
				t.Fatal("bad record recovered")
			}
			records, err := rm.FindRecordsByZone(tt.args.zoneName)
			if err != nil {
				t.Fatal(err)
			}
			if len(*records) != 1 {
				t.Fatal("bad amount of records recovered")
			}
		})
	}
}
