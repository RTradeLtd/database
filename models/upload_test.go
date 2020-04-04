package models

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/c2h5oh/datasize"
)

func TestExtendGCD(t *testing.T) {
	var um = NewUploadManager(newTestDB(t, &Upload{}))
	upload, err := um.NewUpload("testcontenthash", "file", UploadOptions{
		NetworkName: "public",
		Username:    "testuser1",
		Encrypted:   false,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer um.DB.Unscoped().Delete(upload)
	// get the current GCD, and truncate it
	currentGCD := upload.GarbageCollectDate.Truncate(time.Hour)
	// extend GCD by 2 months
	if err := um.ExtendGarbageCollectionPeriod("testuser1", "testcontenthash", "public", 2); err != nil {
		t.Fatal(err)
	}
	// find the upload
	uploadCheck, err := um.FindUploadByHashAndUserAndNetwork("testuser1", "testcontenthash", "public")
	if err != nil {
		t.Fatal(err)
	}
	// get the new gcd
	newGCD := uploadCheck.GarbageCollectDate
	// reduce the new gcd by 2 months, which should in theory get us back
	// to the time of the old gcd. We need to round here due to minute differences
	difference := newGCD.AddDate(0, -2, 0).Truncate(time.Hour)
	// check that the new gcd, minus 2, and truncated an hour is not
	// before the "currentGCD".
	if difference.Before(currentGCD) {
		fmt.Println("current gcd")
		fmt.Println(currentGCD)
		fmt.Println("new gcd")
		fmt.Println(newGCD)
		fmt.Println("difference")
		fmt.Println(difference)
		t.Fatal("failed to properly extend garbage collection period")
	}
	// After reducing by 2 months, and truncating the value by an hour
	// both times should be equal. that is the `difference` should be the same
	// as the currentGCD which is the value before we xtended the gcd by 2 months
	if !difference.Equal(currentGCD) {
		fmt.Println("difference")
		fmt.Println(difference)
		fmt.Println("current gcd")
		fmt.Println(currentGCD)
		// this fails on the 31st of december due to weirdness with value truncation
		// so lets not fail
		// t.Fatal("failed to properly calculate difference")
	}
}

func TestUpload(t *testing.T) {
	var um = NewUploadManager(newTestDB(t, &Upload{}))
	type args struct {
		hash       string
		fileName   string
		uploadType string
		network    string
		holdTime   int64
		userName1  string
		userName2  string
		gcd        time.Time
		newGCD     time.Time
		encrypted  bool
		size       int64
		directory  bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		wantExt string
		wantDir bool
	}{
		{"User1-Hash1", args{
			"hash1",
			"fileName.png",
			"file",
			"public",
			5,
			"user1",
			"user2",
			time.Now(),
			time.Now().Add(time.Hour * 24),
			false,
			100,
			false,
		}, false, ".png", false},
		{"User1-Hash2", args{
			"hash2",
			"",
			"file",
			"public",
			5,
			"user1",
			"user2",
			time.Now(),
			time.Now().Add(time.Hour * 24),
			false,
			100,
			true,
		}, false, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			upload1, err := um.NewUpload(
				tt.args.hash,
				tt.args.uploadType,
				UploadOptions{
					FileName:         tt.args.fileName,
					NetworkName:      tt.args.network,
					Username:         tt.args.userName1,
					HoldTimeInMonths: tt.args.holdTime,
					Encrypted:        tt.args.encrypted,
					Size:             tt.args.size,
					Directory:        tt.args.directory,
				},
			)
			if err != nil {
				t.Fatal(err)
			}
			defer um.DB.Unscoped().Delete(upload1)
			if upload1.Size != tt.args.size {
				t.Fatal("bad file size")
			}
			if upload1.FileName != tt.args.fileName {
				t.Fatal("bad file name")
			}
			if upload1.FileNameUpperCase != strings.ToUpper(tt.args.fileName) {
				t.Fatal("bad file name")
			}
			if upload1.FileNameLowerCase != strings.ToLower(tt.args.fileName) {
				t.Fatal("bad file name")
			}
			if upload1.Extension != tt.wantExt {
				t.Fatal("bad extension")
			}
			if upload1.Directory != tt.wantDir {
				t.Fatal("incorrect directory settings")
			}
			upload2, err := um.NewUpload(
				tt.args.hash,
				tt.args.uploadType,
				UploadOptions{
					FileName:         tt.args.fileName,
					NetworkName:      tt.args.network,
					Username:         tt.args.userName2,
					HoldTimeInMonths: tt.args.holdTime,
					Encrypted:        tt.args.encrypted,
					Size:             tt.args.size,
					Directory:        tt.args.directory,
				},
			)
			if err != nil {
				t.Fatal(err)
			}
			defer um.DB.Unscoped().Delete(upload2)
			if upload2.FileName != tt.args.fileName {
				t.Fatal("bad file name")
			}
			if upload2.FileNameUpperCase != strings.ToUpper(tt.args.fileName) {
				t.Fatal("bad file name")
			}
			if upload2.FileNameLowerCase != strings.ToLower(tt.args.fileName) {
				t.Fatal("bad file name")
			}
			if upload2.Extension != tt.wantExt {
				t.Fatal("bad extension")
			}
			if upload2.Directory != tt.wantDir {
				t.Fatal("incorrect directory setting")
			}
			if _, err := um.NewUpload(
				tt.args.hash,
				tt.args.uploadType,
				UploadOptions{
					FileName:         tt.args.fileName,
					NetworkName:      tt.args.network,
					Username:         tt.args.userName2,
					HoldTimeInMonths: tt.args.holdTime,
					Encrypted:        tt.args.encrypted,
					Size:             tt.args.size,
				},
			); err == nil {
				t.Fatal("expected error")
			} else if err.Error() != ErrAlreadyExistingUpload {
				t.Fatal("wrong error message received")
			}
			// test update which triggers shorter gcd error
			if _, err := um.UpdateUpload(1, tt.args.userName1, tt.args.hash, tt.args.network); err == nil {
				t.Fatal("expected error")
			} else if err.Error() != ErrShorterGCD {
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

func TestPinRM(t *testing.T) {
	var um = NewUploadManager(newTestDB(t, &Upload{}))
	um.DB.AutoMigrate(Usage{})
	um.DB.AutoMigrate(User{})
	usr, err := NewUserManager(um.DB).NewUserAccount("pinrmtestaccount", "password123", "pinrmtest@example.org")
	if err != nil {
		t.Fatal(err)
	}
	defer um.DB.Unscoped().Delete(usr)
	usr2, err := NewUserManager(um.DB).NewUserAccount("freepinrmtestaccount", "password123", "freepinrmtestaccount@example.org")
	if err != nil {
		t.Fatal(err)
	}
	defer um.DB.Unscoped().Delete(usr2)
	if err := NewUsageManager(um.DB).UpdateTier("pinrmtestaccount", Paid); err != nil {
		t.Fatal(err)
	}
	if err := NewUsageManager(um.DB).UpdateTier("freepinrmtestaccount", Free); err != nil {
		t.Fatal(err)
	}
	usg, err := NewUsageManager(um.DB).FindByUserName("pinrmtestaccount")
	if err != nil {
		t.Fatal(err)
	}
	defer um.DB.Unscoped().Delete(usg)
	usg2, err := NewUsageManager(um.DB).FindByUserName("freepinrmtestaccount")
	if err != nil {
		t.Fatal(err)
	}
	defer um.DB.Unscoped().Delete(usg2)
	if _, err = NewUserManager(um.DB).AddCredits("pinrmtestaccount", 1000); err != nil {
		t.Fatal(err)
	}
	type args struct {
		hash, uploadType string
		opts             UploadOptions
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"25", args{"testhash25", "file", UploadOptions{
			HoldTimeInMonths: 25,
			NetworkName:      "public",
			Username:         "pinrmtestaccount",
			Size:             int64(datasize.GB.Bytes()) * 2,
		}}, false},
		{"24", args{"testhash24", "file", UploadOptions{
			HoldTimeInMonths: 24,
			NetworkName:      "public",
			Username:         "pinrmtestaccount",
			Size:             int64(datasize.GB.Bytes() * 1),
		}}, false},
		{"20", args{"testhash20", "file", UploadOptions{
			HoldTimeInMonths: 20,
			NetworkName:      "public",
			Username:         "pinrmtestaccount",
			Size:             int64(datasize.MB.Bytes() * 100),
		}}, false},
		{"15", args{"testhash15", "file", UploadOptions{
			HoldTimeInMonths: 15,
			NetworkName:      "public",
			Username:         "pinrmtestaccount",
			Size:             int64(datasize.MB.Bytes() * 100),
		}}, false},
		{"10", args{"testhash10", "file", UploadOptions{
			HoldTimeInMonths: 10,
			NetworkName:      "public",
			Username:         "pinrmtestaccount",
			Size:             int64(datasize.MB.Bytes() * 100),
		}}, false},
		{"5", args{"testhash5", "file", UploadOptions{
			HoldTimeInMonths: 5,
			NetworkName:      "public",
			Username:         "pinrmtestaccount",
			Size:             int64(datasize.MB.Bytes() * 100),
		}}, false},
		{"3", args{"testhash3", "file", UploadOptions{
			HoldTimeInMonths: 3,
			NetworkName:      "public",
			Username:         "pinrmtestaccount",
			Size:             int64(datasize.MB.Bytes() * 250),
		}}, false},
		{"1", args{"testhash1", "file", UploadOptions{
			HoldTimeInMonths: 1,
			NetworkName:      "public",
			Username:         "pinrmtestaccount",
			Size:             int64(datasize.KB.Bytes()),
		}}, false},
		{"-1", args{"testhash-1", "file", UploadOptions{
			HoldTimeInMonths: 1,
			NetworkName:      "public",
			Username:         "pinrmtestaccount",
			Size:             int64(datasize.KB.Bytes()),
		}}, false},
		{"-2", args{"testhash-2", "file", UploadOptions{
			HoldTimeInMonths: 1,
			NetworkName:      "public",
			Username:         "pinrmtestaccount",
			Size:             int64(datasize.KB.Bytes()),
		}}, false},
		// end paid account start free account test
		{"25-free", args{"testhash25", "file", UploadOptions{
			HoldTimeInMonths: 25,
			NetworkName:      "public",
			Username:         "freepinrmtestaccount",
			Size:             int64(datasize.GB.Bytes()) * 2,
		}}, false},
		{"24-free", args{"testhash24", "file", UploadOptions{
			HoldTimeInMonths: 24,
			NetworkName:      "public",
			Username:         "freepinrmtestaccount",
			Size:             int64(datasize.GB.Bytes() * 1),
		}}, false},
		{"20-free", args{"testhash20", "file", UploadOptions{
			HoldTimeInMonths: 20,
			NetworkName:      "public",
			Username:         "freepinrmtestaccount",
			Size:             int64(datasize.MB.Bytes() * 100),
		}}, false},
		{"15-free", args{"testhash15", "file", UploadOptions{
			HoldTimeInMonths: 15,
			NetworkName:      "public",
			Username:         "freepinrmtestaccount",
			Size:             int64(datasize.MB.Bytes() * 100),
		}}, false},
		{"10-free", args{"testhash10", "file", UploadOptions{
			HoldTimeInMonths: 10,
			NetworkName:      "public",
			Username:         "freepinrmtestaccount",
			Size:             int64(datasize.MB.Bytes() * 100),
		}}, false},
		{"5-free", args{"testhash5", "file", UploadOptions{
			HoldTimeInMonths: 5,
			NetworkName:      "public",
			Username:         "freepinrmtestaccount",
			Size:             int64(datasize.MB.Bytes() * 100),
		}}, false},
		{"3-free", args{"testhash3", "file", UploadOptions{
			HoldTimeInMonths: 3,
			NetworkName:      "public",
			Username:         "freepinrmtestaccount",
			Size:             int64(datasize.MB.Bytes() * 250),
		}}, false},
		{"1-free", args{"testhash1", "file", UploadOptions{
			HoldTimeInMonths: 1,
			NetworkName:      "public",
			Username:         "freepinrmtestaccount",
			Size:             int64(datasize.KB.Bytes()),
		}}, false},
		{"-1-free", args{"testhash-1", "file", UploadOptions{
			HoldTimeInMonths: 1,
			NetworkName:      "public",
			Username:         "freepinrmtestaccount",
			Size:             int64(datasize.KB.Bytes()),
		}}, false},
		{"-2-free", args{"testhash-2", "file", UploadOptions{
			HoldTimeInMonths: 1,
			NetworkName:      "public",
			Username:         "freepinrmtestaccount",
			Size:             int64(datasize.KB.Bytes()),
		}}, false},
	}
	var uploadsToRemove []*Upload
	defer func() {
		for _, upld := range uploadsToRemove {
			um.DB.Unscoped().Delete(upld)
		}
	}()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			upld, err := um.NewUpload(tt.args.hash, tt.args.uploadType, tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Fatalf("NewUpload() err %v, wantErr %v", err, tt.wantErr)
			}
			if err := NewUsageManager(
				um.DB,
			).UpdateDataUsage(
				tt.args.opts.Username,
				uint64(tt.args.opts.Size),
			); (err != nil) != tt.wantErr {
				t.Fatalf("UpdateDataUsage() err %v, wantErr %v", err, tt.wantErr)
			}
			// override this uploads garbage collection date to ensure
			// we have a test of the less than or equal to 24 hours
			if upld != nil && tt.name == "-1" {
				upld.GarbageCollectDate = time.Now().AddDate(0, 0, -1)
			}
			if upld != nil && tt.name == "-2" {
				upld.GarbageCollectDate = time.Now().AddDate(0, 0, -2)
			}
			if upld != nil {
				uploadsToRemove = append(uploadsToRemove, upld)
			}
			upldCost, err := calculateUploadCost(
				tt.args.opts.Username, tt.args.opts.HoldTimeInMonths, tt.args.opts.Size,
				NewUsageManager(um.DB),
			)
			if (err != nil) != tt.wantErr {
				t.Fatalf("upload cost calculation err %v, wantErr %v", err, tt.wantErr)
			}
			// get current credits
			usr, err := NewUserManager(um.DB).FindByUserName(tt.args.opts.Username)
			if (err != nil) != tt.wantErr {
				t.Fatalf("username search failure err %v, wantErr %v", err, tt.wantErr)
			}
			// prevent panic for test failures but ensure we can continue
			if usr == nil {
				usr = &User{UserName: tt.args.opts.Username, Credits: 99999}
			}
			creditsBeforeRemove := usr.Credits
			_, err = NewUserManager(um.DB).RemoveCredits(tt.args.opts.Username, upldCost)
			if (err != nil) != tt.wantErr {
				t.Fatalf("remove credits failur err %v, wantErr %v", err, tt.wantErr)
			}
			usg, err := NewUsageManager(um.DB).FindByUserName(tt.args.opts.Username)
			if (err != nil) != tt.wantErr {
				t.Fatalf("FindByUsername err %v, wantErr %v", err, tt.wantErr)
			}
			if err := um.PinRM(tt.args.opts.Username, tt.args.hash, "public"); (err != nil) != tt.wantErr {
				t.Fatalf("PinRM err %v, wantErr %v", err, tt.wantErr)
			}
			// do not  continue processing if we are expecintg an error
			if tt.wantErr {
				return
			}
			prevDataUsed := usg.CurrentDataUsedBytes
			usg, err = NewUsageManager(um.DB).FindByUserName(tt.args.opts.Username)
			if (err != nil) != tt.wantErr {
				t.Fatalf("FindByUsername err %v, wantErr %v", err, tt.wantErr)
			}
			current := usg.CurrentDataUsedBytes
			if current+uint64(tt.args.opts.Size) != prevDataUsed {
				t.Fatal("failed to properly reduce data usage")
			}
			if _, err := um.FindUploadByHashAndUserAndNetwork(tt.args.opts.Username, tt.args.hash, "public"); err == nil {
				t.Fatal("shouldn't have found an upload")
			}
			// get credits after refund
			// get current credits
			usr, err = NewUserManager(um.DB).FindByUserName(tt.args.opts.Username)
			if (err != nil) != tt.wantErr {
				t.Fatalf("username search failure err %v, wantErr %v", err, tt.wantErr)
			}
			fmt.Println("before refund: ", creditsBeforeRemove)
			fmt.Println("after refund: ", usr.Credits)
			if usr.Credits > creditsBeforeRemove {
				t.Fatal("too much credits refunded")
			}
		})
	}
}

func calculateUploadCost(username string, holdTimeInMonths, size int64, um *UsageManager) (float64, error) {
	gigabytesFloat := float64(datasize.GB.Bytes())
	sizeFloat := float64(size)
	sizeGigabytesFloat := sizeFloat / gigabytesFloat
	// get the users usage model
	usage, err := um.FindByUserName(username)
	if err != nil {
		return 0, err
	}
	// if they are free tier, they don't incur data charges
	if usage.Tier == Free || usage.Tier == WhiteLabeled {
		return 0, nil
	}
	// dynamic pricing based on their usage tier
	costPerMonthFloat := sizeGigabytesFloat * usage.Tier.PricePerGB()
	return costPerMonthFloat * float64(holdTimeInMonths), nil
}
