package models

import (
	"testing"
)

func TestZone(t *testing.T) {
	var zm = NewZoneManager(newTestDB(t, &Zone{}))
	args := struct {
		zoneName           string
		zoneManagerKeyName string
		zonePublicKeyName  string
		ipfshash           string
	}{"testzone", "testzonemanager", "testzonepublic", "testhash"}
	zone1, err := zm.NewZone(
		args.zoneName,
		args.zoneManagerKeyName,
		args.zonePublicKeyName,
		args.ipfshash,
	)
	if err != nil {
		t.Fatal(err)
	}
	defer zm.DB.Delete(zone1)
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
