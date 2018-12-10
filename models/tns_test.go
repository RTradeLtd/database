package models_test

import (
	"testing"

	"github.com/RTradeLtd/config"
	"github.com/RTradeLtd/database/models"
)

func TestMigration_TNS(t *testing.T) {
	cfg, err := config.LoadConfig(testCfgPath)
	if err != nil {
		t.Fatal(err)
	}
	db, err := openDatabaseConnection(t, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if err = db.AutoMigrate(&models.Zone{}).Error; err != nil {
		t.Fatal(err)
	}
}

func TestZone(t *testing.T) {
	cfg, err := config.LoadConfig(testCfgPath)
	if err != nil {
		t.Fatal(err)
	}
	db, err := openDatabaseConnection(t, cfg)
	if err != nil {
		t.Fatal(err)
	}
	args := struct {
		zoneName           string
		zoneManagerKeyName string
		zonePublicKeyName  string
		ipfshash           string
	}{"testzone", "testzonemanager", "testzonepublic", "testhash"}
	zm := models.NewZoneManager(db)
	zone1, err := zm.NewZone(
		args.zoneName,
		args.zoneManagerKeyName,
		args.zonePublicKeyName,
		args.ipfshash,
	)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Delete(zone1)
	zone2, err := zm.FindZoneByName(args.zoneName)
	if err != nil {
		t.Fatal(err)
	}
	if zone2.LatestIPFSHash != zone1.LatestIPFSHash {
		t.Fatal("bad hash recovered")
	}
	zone3, err := zm.UpdateLatestIPFSHashForZone(args.zoneName, "newhash")
	if err != nil {
		t.Fatal(err)
	}
	if zone3.LatestIPFSHash != "newhash" {
		t.Fatal("bad hash recovered")
	}
	zone4, err := zm.AddRecordForZone(args.zoneName, "testrecord1")
	if err != nil {
		t.Fatal(err)
	}
	if len(zone4.RecordNames) != 1 {
		t.Fatal("bad record count recovered")
	}
	zone5, err := zm.AddRecordForZone(args.zoneName, "testrecord2")
	if err != nil {
		t.Fatal(err)
	}
	if len(zone5.RecordNames) != 2 {
		t.Fatal("bad record count recovered")
	}
}
