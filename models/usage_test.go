package models_test

import (
	"testing"

	"github.com/c2h5oh/datasize"

	"github.com/RTradeLtd/config"
	"github.com/RTradeLtd/database/models"
)

func TestMigration_Usage(t *testing.T) {
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

func TestUsage(t *testing.T) {
	cfg, err := config.LoadConfig(testCfgPath)
	if err != nil {
		t.Fatal(err)
	}
	db, err := openDatabaseConnection(t, cfg)
	if err != nil {
		t.Fatal(err)
	}
	bm := models.NewUsageManager(db)
	type args struct {
		username       string
		tier           models.DataUsageTier
		testUploadSize float64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Free", args{"free", models.Free, datasize.GB.GBytes()}, false},
		{"Partner", args{"partner", models.Partner, datasize.GB.GBytes() * 10}, false},
		{"Light", args{"light", models.Light, datasize.GB.GBytes() * 100}, false},
		{"Plus", args{"plus", models.Plus, datasize.GB.GBytes() * 10}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// test create usage
			usage, err := bm.NewUsageEntry(tt.args.username, tt.args.tier)
			if (err != nil) != tt.wantErr {
				t.Fatalf("NewUsage() err = %v, wantErr %v", err, tt.wantErr)
			}
			defer bm.DB.Delete(usage)
			// test find by username
			if _, err := bm.FindByUserName(tt.args.username); (err != nil) != tt.wantErr {
				t.Fatalf("FindByUserName() err = %v, wantErr %v", err, tt.wantErr)
			}
			// test get upload price
			if price, err := bm.GetUploadPricePerGB(tt.args.username); (err != nil) != tt.wantErr {
				t.Fatalf("GetUploadPricePerGB() err = %v, wantErr %v", err, tt.wantErr)
			} else if !tt.wantErr && price != usage.Tier.PricePerGB() {
				t.Fatal("failed to get correct price per gb")
			}
			// test ipns publish check
			if canPub, err := bm.CanPublishIPNS(tt.args.username); (err != nil) != tt.wantErr {
				t.Fatalf("CanPublishIPNS() err = %v, wantErr %v", err, tt.wantErr)
			} else if !canPub {
				t.Fatal("error occured validating ipns publish")
			}
			// test ipns publish check
			if canPub, err := bm.CanPublishPubSub(tt.args.username); (err != nil) != tt.wantErr {
				t.Fatalf("CanPublishPubSub() err = %v, wantErr %v", err, tt.wantErr)
			} else if !canPub {
				t.Fatal("error occured validating ipns publish")
			}
			// test update data usage
			if err := bm.UpdateDataUsage(tt.args.username, tt.args.testUploadSize); (err != nil) != tt.wantErr {
				t.Fatalf("UpdateDataUsage() err = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.name == "Light" && !tt.wantErr {
				// validate that the tier was upgraded
				usage, err = bm.FindByUserName(tt.args.username)
				if err != nil {
					t.Fatal(err)
				}
				if usage.Tier != models.Plus {
					t.Fatal("failed to correctly set usage tier")
				}
			}
		})
	}
}
