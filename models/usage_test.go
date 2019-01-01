package models_test

import (
	"testing"

	"github.com/RTradeLtd/config"
	"github.com/RTradeLtd/database/models"
)

func TestMigration_Billing(t *testing.T) {
	cfg, err := config.LoadConfig(testCfgPath)
	if err != nil {
		t.Fatal(err)
	}
	db, err := openDatabaseConnection(t, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if check := db.AutoMigrate(&models.Usage{}); check.Error != nil {
		t.Fatal(err)
	}
}

func TestNewUsage(t *testing.T) {
	cfg, err := config.LoadConfig(testCfgPath)
	if err != nil {
		t.Fatal(err)
	}

	db, err := openDatabaseConnection(t, cfg)
	if err != nil {
		t.Fatal(err)
	}
	// create Usage entry
	Usage := &models.Usage{
		UserName: username,
		// default data limit of 3GB
		MonthlyDataLimitGB:      3.0,
		CurrentDataUsedGB:       0.0,
		PrivateNetworkTrialUsed: false,
		// if this is 0, they have not started
		// otherwise it is in unix nano
		TrialEndTime: 0,
		Tier:         models.Free,
	}
	if err := db.Create(Usage).Error; err != nil {
		t.Fatal(err)
	}
	bm := models.NewUsageManager(db)
	b, err := bm.FindByUserName(username)
	if err != nil {
		t.Fatal(err)
	}
	defer bm.DB.Delete(b)
}
