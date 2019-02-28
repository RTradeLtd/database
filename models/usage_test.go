package models

import (
	"testing"

	"github.com/c2h5oh/datasize"
)

func TestUsage(t *testing.T) {
	var bm = NewUsageManager(newTestDB(t, &Usage{}))
	type args struct {
		username       string
		tier           DataUsageTier
		testUploadSize uint64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Free", args{"free", Free, datasize.GB.Bytes()}, false},
		{"Partner", args{"partner", Partner, datasize.GB.Bytes() * 10}, false},
		{"Light", args{"light", Light, datasize.GB.Bytes() * 100}, false},
		{"Plus", args{"plus", Plus, datasize.GB.Bytes() * 10}, false},
		{"Fail", args{"fail", Free, 1}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				usage *Usage
				err   error
			)
			if !tt.wantErr {
				// test create usage
				usage, err = bm.NewUsageEntry(tt.args.username, tt.args.tier)
				if (err != nil) != tt.wantErr {
					t.Fatalf("NewUsage() err = %v, wantErr %v", err, tt.wantErr)
				}
				defer bm.DB.Delete(usage)
			}
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
			if err := bm.CanPublishIPNS(tt.args.username); (err != nil) != tt.wantErr {
				t.Fatalf("CanPublishIPNS() err = %v, wantErr %v", err, tt.wantErr)
			}
			// test ipns publish check
			if err := bm.CanPublishPubSub(tt.args.username); (err != nil) != tt.wantErr {
				t.Fatalf("CanPublishPubSub() err = %v, wantErr %v", err, tt.wantErr)
			}
			// test update data usage
			if err := bm.UpdateDataUsage(tt.args.username, tt.args.testUploadSize); (err != nil) != tt.wantErr {
				t.Fatalf("UpdateDataUsage() err = %v, wantErr %v", err, tt.wantErr)
			}
			// test update tiers for all tier types
			// an account may never enter free status once exiting
			tiers := []DataUsageTier{Partner, Light, Plus}
			for _, tier := range tiers {
				if err := bm.UpdateTier(tt.args.username, tier); (err != nil) != tt.wantErr {
					t.Fatalf("UpdateTier() err = %v, wantErr %v", err, tt.wantErr)
				}
			}
			// test that the light tier was upgraded
			if tt.name == "Light" && !tt.wantErr {
				// validate that the tier was upgraded
				usage, err = bm.FindByUserName(tt.args.username)
				if err != nil {
					t.Fatal(err)
				}
				if usage.Tier != Plus {
					t.Fatal("failed to correctly set usage tier")
				}
			}
			// test pubsub increment
			if err := bm.IncrementPubSubUsage(tt.args.username, 5); (err != nil) != tt.wantErr {
				t.Fatalf("IncrementPubSubUsage() err = %v, wantErr %v", err, tt.wantErr)
			}
			// if no error is expected, validate the pubsub count
			if !tt.wantErr {
				usage, err := bm.FindByUserName(tt.args.username)
				if err != nil {
					t.Fatal(err)
				}
				if usage.PubSubMessagesSent != 5 {
					t.Fatal("failed to count pubsub usage")
				}
			}
			// test ipns increment
			if err := bm.IncrementIPNSUsage(tt.args.username, 5); (err != nil) != tt.wantErr {
				t.Fatalf("IncrementIPNSUsage() err = %v, wantErr %v", err, tt.wantErr)
			}
			// if no error is expected, validate the ipns count
			if !tt.wantErr {
				usage, err := bm.FindByUserName(tt.args.username)
				if err != nil {
					t.Fatal(err)
				}
				if usage.IPNSRecordsPublished != 5 {
					t.Fatal("failed to count ipns usage")
				}
			}
		})
	}
}
