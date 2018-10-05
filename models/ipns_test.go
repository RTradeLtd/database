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
		{"Test1", args{"12D3KooWSev8mmycrPbCMs4Awe4AFGkUQKPh7CTuifh51U8iFEr8", "QmQxXGDe84eUjCg2ZspvduEZxjWZk5DCB2N7bwPjXahoXE", "key", "public", time.Hour, time.Hour, "username"}},
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
			im.DB.Unscoped().Delete(entry)
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
		{"Test1", args{"12D3KooWSev8mmycrPbCMs4Awe4AFGkUQKPh7CTuifh51U8iFEr8", "QmQxXGDe84eUjCg2ZspvduEZxjWZk5DCB2N7bwPjXahoXE", "key", "public", time.Hour, time.Hour, "username"}},
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
			defer im.DB.Unscoped().Delete(entry)
			entryCopy, err := im.UpdateIPNSEntry(
				tt.args.ipnsHash,
				newIpfsHash,
				tt.args.networkName,
				tt.args.userName,
				tt.args.lifetime,
				tt.args.ttl,
			)
			if err != nil {
				t.Fatal(err)
			}
			if entryCopy.Sequence <= entry.Sequence {
				t.Fatal("failed to update sequence")
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
		{"Test1", args{"12D3KooWSev8mmycrPbCMs4Awe4AFGkUQKPh7CTuifh51U8iFEr8", "QmQxXGDe84eUjCg2ZspvduEZxjWZk5DCB2N7bwPjXahoXE", "key", "public", time.Hour, time.Hour, "username"}},
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
			defer im.DB.Unscoped().Delete(entry)
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
		{"Test1", args{"12D3KooWSev8mmycrPbCMs4Awe4AFGkUQKPh7CTuifh51U8iFEr8", "QmQxXGDe84eUjCg2ZspvduEZxjWZk5DCB2N7bwPjXahoXE", "key", "public", time.Hour, time.Hour, "username"}},
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
			defer im.DB.Unscoped().Delete(entry)
			if _, err := im.FindByUserName(tt.args.userName); err != nil {
				t.Fatal(err)
			}
		})
	}
}
