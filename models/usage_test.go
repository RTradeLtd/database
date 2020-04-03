package models

import (
	"fmt"
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
		{"Paid", args{"paid", Paid, datasize.GB.Bytes() * 100}, false},
		{"WhiteLabelled", args{"whitelabelled", WhiteLabeled, datasize.GB.Bytes() * 100}, false},
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
				defer bm.DB.Unscoped().Delete(usage)
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
			if err := bm.CanCreateKey(tt.args.username); (err != nil) != tt.wantErr {
				t.Fatalf("CanCreateKey() err = %v, wantErr %v", err, tt.wantErr)
			}
			// test update data usage
			if err := bm.UpdateDataUsage(tt.args.username, tt.args.testUploadSize); (err != nil) != tt.wantErr {
				t.Fatalf("UpdateDataUsage() err = %v, wantErr %v", err, tt.wantErr)
			}
			// test update tiers for all tier types
			// an account may never enter free status once exiting
			tiers := []DataUsageTier{Paid, Partner}
			for _, tier := range tiers {
				// skip white labelled acounts
				if tt.args.tier == WhiteLabeled {
					continue
				}
				if err := bm.UpdateTier(tt.args.username, tier); (err != nil) != tt.wantErr {
					t.Fatalf("UpdateTier() err = %v, wantErr %v", err, tt.wantErr)
				}
			}
			// test that the light tier was upgraded
			if tt.name == "Paid" && !tt.wantErr {
				// validate that the tier was upgraded
				usage, err = bm.FindByUserName(tt.args.username)
				if err != nil {
					t.Fatal(err)
				}
				if usage.Tier != Partner {
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
			if !tt.wantErr {
				if err := bm.ResetCounts(tt.args.username); err != nil {
					t.Fatal(err)
				}
				usage, err := bm.FindByUserName(tt.args.username)
				if err != nil {
					t.Fatal(err)
				}
				if usage.IPNSRecordsPublished != 0 || usage.PubSubMessagesSent != 0 {
					t.Fatal("should be 0")
				}
			}
		})
	}
}

func Test_ENS(t *testing.T) {
	var bm = NewUsageManager(newTestDB(t, &Usage{}))
	type args struct {
		user string
		tier DataUsageTier
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		wantClaim bool
	}{
		{"Register-Success", args{"testuser1", Paid}, false, true},
		{"Register-Fail-Already-Claimed", args{"testuser1", Paid}, true, true},
		{"Regiser-Fail-Free-Tier", args{"testuser2", Free}, true, false},
	}
	b, err := bm.NewUsageEntry("testuser1", Paid)
	if err != nil {
		t.Fatal(err)
	}
	defer bm.DB.Unscoped().Delete(b)
	b, err = bm.NewUsageEntry("testuser2", Free)
	if err != nil {
		t.Fatal(err)
	}
	defer bm.DB.Unscoped().Delete(b)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := bm.ClaimENSName(tt.args.user); (err != nil) != tt.wantErr {
				t.Fatalf("ClaimENSName() err %v, wantErr %v", err, tt.wantErr)
			}
			b, err := bm.FindByUserName(tt.args.user)
			if err != nil {
				t.Fatal(err)
			}
			if b.ClaimedENSName != tt.wantClaim {
				t.Fatal("ens claim is incorrect")
			}
		})
	}
	if err := bm.UnclaimENSName("testuser1"); err != nil {
		t.Fatal(err)
	}
	if err := bm.UnclaimENSName("testuser1"); err == nil {
		t.Fatal("error expected")
	}
}

func Test_Tier_Upgrade(t *testing.T) {
	var bm = NewUsageManager(newTestDB(t, &Usage{}))
	b, err := bm.NewUsageEntry("testuser", Free)
	if err != nil {
		t.Fatal(err)
	}
	defer bm.DB.Unscoped().Delete(b)
	if b.Tier != Free {
		t.Fatal("bad tier set")
	}
	if b.MonthlyDataLimitBytes != FreeUploadLimit {
		t.Fatal("bad upload limit set")
	}
	if err := bm.UpdateTier("testuser", Paid); err != nil {
		t.Fatal(err)
	}
	b, err = bm.FindByUserName("testuser")
	if err != nil {
		t.Fatal(err)
	}
	if b.Tier != Paid {
		t.Fatal("bad tier set")
	}
	if b.MonthlyDataLimitBytes != NonFreeUploadLimit {
		t.Fatal("bad upload limit set")
	}
}

func Test_UpdateDataUsage_Free(t *testing.T) {
	var bm = NewUsageManager(newTestDB(t, &Usage{}))
	b, err := bm.NewUsageEntry("testuser", Free)
	if err != nil {
		t.Fatal(err)
	}
	defer bm.DB.Unscoped().Delete(b)
	if b.Tier != Free {
		t.Fatal("bad tier set")
	}
	b.CurrentDataUsedBytes = datasize.GB.Bytes() * 2
	if err := bm.DB.Save(b).Error; err != nil {
		t.Fatal(err)
	}
	b, err = bm.FindByUserName("testuser")
	if err != nil {
		t.Fatal(err)
	}
	if b.CurrentDataUsedBytes != datasize.GB.Bytes()*2 {
		t.Fatal("bad usage set")
	}
	if err := bm.UpdateDataUsage("testuser", datasize.GB.Bytes()*2); err == nil {
		t.Fatal("error expected")
	}
	if err := bm.UpdateDataUsage("testuser", datasize.MB.Bytes()*100); err != nil {
		t.Fatal(err)
	}
}

func Test_ReduceDataUsage(t *testing.T) {
	var bm = NewUsageManager(newTestDB(t, &Usage{}))
	b, err := bm.NewUsageEntry("testuser", Paid)
	if err != nil {
		t.Fatal(err)
	}
	defer bm.DB.Unscoped().Delete(b)
	if b.Tier != Paid {
		t.Fatal("bad tier set")
	}
	if err := bm.UpdateDataUsage(
		"testuser",
		datasize.GB.Bytes()+datasize.MB.Bytes()*100,
	); err != nil {
		t.Fatal(err)
	}
	b, err = bm.FindByUserName("testuser")
	if err != nil {
		t.Fatal(err)
	}
	if b.CurrentDataUsedBytes != datasize.GB.Bytes()+datasize.MB.Bytes()*100 {
		t.Fatal("bad datasize")
	}
	currentSize := b.CurrentDataUsedBytes
	expectedSize := b.CurrentDataUsedBytes - datasize.MB.Bytes()*100
	if err := bm.ReduceDataUsage("testuser", datasize.MB.Bytes()*100); err != nil {
		t.Fatal(err)
	}
	b, err = bm.FindByUserName("testuser")
	if b.CurrentDataUsedBytes != expectedSize {
		fmt.Println("current size", currentSize)
		fmt.Println("got size", b.CurrentDataUsedBytes)
		fmt.Println("expected size", expectedSize)
		t.Fatal("bad reduction in datasize")
	}
}

func Test_ReduceKeyCount(t *testing.T) {
	var bm = NewUsageManager(newTestDB(t, &Usage{}))
	b, err := bm.NewUsageEntry("testuser", Paid)
	if err != nil {
		t.Fatal(err)
	}
	defer bm.DB.Unscoped().Delete(b)
	if b.Tier != Paid {
		t.Fatal("bad tier set")
	}
	if err := bm.IncrementKeyCount("testuser", 5); err != nil {
		t.Fatal(err)
	}
	if err := bm.ReduceKeyCount("testuser", 4); err != nil {
		t.Fatal(err)
	}
	b, err = bm.FindByUserName("testuser")
	if err != nil {
		t.Fatal(err)
	}
	if b.KeysCreated != 1 {
		t.Fatal("bad key count")
	}
	if err := bm.ReduceKeyCount("testuser", 3); err != nil {
		t.Fatal(err)
	}
	b, err = bm.FindByUserName("testuser")
	if err != nil {
		t.Fatal(err)
	}
	if b.KeysCreated != 0 {
		t.Fatal("bad key count")
	}
}

func TestPricePerGB(t *testing.T) {
	tests := []struct {
		tier             DataUsageTier
		wantPriceMonthly float64
		wantPriceHourly  float64
		wantString       string
	}{
		{Paid, 0.07, 0.07 / 730, "paid"},
		{Partner, 0.05, 0.05 / 730, "partner"},
		{WhiteLabeled, 0.05, 0.05 / 730, "white-labeled"},
		{Free, 9999, 9999, "free"},
	}
	for _, tt := range tests {
		if tt.tier.PricePerGB() != tt.wantPriceMonthly {
			t.Fatal("bad PricePerGB returned")
		}
		if tt.tier.PricePerGBPerHour() != tt.wantPriceHourly {
			t.Fatal("bad hourly price returned")
		}
		if tt.tier.String() != tt.wantString {
			t.Fatal("bad string returned")
		}
	}
}
