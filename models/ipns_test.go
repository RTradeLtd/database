package models_test

import (
	"testing"
	"time"

	"github.com/RTradeLtd/config"
	"github.com/RTradeLtd/database/models"
)

const (
	newIpfsHash = "newHash"
)

func TestIpnsManager_NewEntry(t *testing.T) {
	cfg, err := config.LoadConfig(testCfgPath)
	if err != nil {
		t.Fatal(err)
	}
	db, err := openDatabaseConnection(t, cfg)
	if err != nil {
		t.Fatal(err)
	}
	im := models.NewIPNSManager(db)
	type args struct {
		ipnsHash    string
		ipfsHash    string
		key         string
		networkName string
		lifetime    time.Duration
		ttl         time.Duration
		userName    string
	}
	tests := []struct {
		name string
		args args
	}{
		{"Test1", args{"ipnsHash", "ipfsHash", "key", "public", time.Hour, time.Hour, "username"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry, err := im.CreateEntry(
				tt.args.ipnsHash,
				tt.args.ipfsHash,
				tt.args.key,
				tt.args.networkName,
				tt.args.userName,
				tt.args.lifetime,
				tt.args.ttl,
			)
			if err != nil {
				t.Fatal(err)
			}
			db.Delete(entry)
		})
	}
}

func TestIpnsManager_UpdateEntry(t *testing.T) {
	cfg, err := config.LoadConfig(testCfgPath)
	if err != nil {
		t.Fatal(err)
	}
	db, err := openDatabaseConnection(t, cfg)
	if err != nil {
		t.Fatal(err)
	}
	im := models.NewIPNSManager(db)
	type args struct {
		ipnsHash    string
		ipfsHash    string
		key         string
		networkName string
		lifetime    time.Duration
		ttl         time.Duration
		userName    string
	}
	tests := []struct {
		name string
		args args
	}{
		{"Test1", args{"ipnsHash", "ipfsHash", "key", "public", time.Hour, time.Hour, "username"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry, err := im.CreateEntry(
				tt.args.ipnsHash,
				tt.args.ipfsHash,
				tt.args.key,
				tt.args.networkName,
				tt.args.userName,
				tt.args.lifetime,
				tt.args.ttl,
			)
			if err != nil {
				t.Fatal(err)
			}
			defer db.Delete(entry)
			entryCopy, err := im.UpdateIPNSEntry(
				tt.args.ipnsHash,
				newIpfsHash,
				tt.args.key,
				tt.args.networkName,
				tt.args.userName,
				tt.args.lifetime,
				tt.args.ttl,
			)
			if err != nil {
				t.Fatal(err)
			}
			if entryCopy.IPNSHash != entry.IPNSHash {
				t.Fatal("failed to update correct ipns record")
			}
		})
	}
}

func TestIpnsManager_FindByIPNSHash(t *testing.T) {
	cfg, err := config.LoadConfig(testCfgPath)
	if err != nil {
		t.Fatal(err)
	}
	db, err := openDatabaseConnection(t, cfg)
	if err != nil {
		t.Fatal(err)
	}
	im := models.NewIPNSManager(db)
	type args struct {
		ipnsHash    string
		ipfsHash    string
		key         string
		networkName string
		lifetime    time.Duration
		ttl         time.Duration
		userName    string
	}
	tests := []struct {
		name string
		args args
	}{
		{"Test1", args{"ipnsHash", "ipfsHash", "key", "public", time.Hour, time.Hour, "username"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry, err := im.CreateEntry(
				tt.args.ipnsHash,
				tt.args.ipfsHash,
				tt.args.key,
				tt.args.networkName,
				tt.args.userName,
				tt.args.lifetime,
				tt.args.ttl,
			)
			if err != nil {
				t.Fatal(err)
			}
			defer db.Delete(entry)
			entryCopy, err := im.FindByIPNSHash(tt.args.ipnsHash)
			if err != nil {
				t.Fatal(err)
			}
			if entryCopy.CurrentIPFSHash != entry.CurrentIPFSHash {
				t.Fatal("failed to recover correct entry")
			}
		})
	}
}

func TestIpnsManager_FindByUser(t *testing.T) {
	cfg, err := config.LoadConfig(testCfgPath)
	if err != nil {
		t.Fatal(err)
	}
	db, err := openDatabaseConnection(t, cfg)
	if err != nil {
		t.Fatal(err)
	}
	im := models.NewIPNSManager(db)
	type args struct {
		ipnsHash    string
		ipfsHash    string
		key         string
		networkName string
		lifetime    time.Duration
		ttl         time.Duration
		userName    string
	}
	tests := []struct {
		name string
		args args
	}{
		{"Test1", args{"ipnsHash", "ipfsHash", "key", "public", time.Hour, time.Hour, "username"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry, err := im.CreateEntry(
				tt.args.ipnsHash,
				tt.args.ipfsHash,
				tt.args.key,
				tt.args.networkName,
				tt.args.userName,
				tt.args.lifetime,
				tt.args.ttl,
			)
			if err != nil {
				t.Fatal(err)
			}
			defer db.Delete(entry)
			if _, err := im.FindByUserName(tt.args.userName); err != nil {
				t.Fatal(err)
			}
		})
	}
}
