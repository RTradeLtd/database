package models

import (
	"testing"
	"time"

	"github.com/k0kubun/pp"
)

func TestOrganizationManager_Full(t *testing.T) {
	var om = NewOrgManager(newTestDB(t, &Organization{}))
	om.DB.AutoMigrate(User{})
	om.DB.AutoMigrate(Upload{})
	type args struct {
		name, owner string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Pass", args{"testorg", "mytestuserownerr"}, false},
		{"Fail", args{"testorg22", "mytestuserownerr"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := om.NewOrganization(
				tt.args.name,
				tt.args.owner,
			); (err != nil) != tt.wantErr {
				t.Fatalf("NewOrganization() err %v, wantErr %v", err, tt.wantErr)
			}
			if _, err := om.RegisterOrgUser(
				tt.args.name, tt.args.owner, "password123", "password123@example.org",
			); (err != nil) != tt.wantErr {
				t.Fatalf("RegisterOrgUser err %v, wantErr %v", err, tt.wantErr)
			}
			if _, err := om.FindByName(tt.args.name); (err != nil) != tt.wantErr {
				t.Fatalf("FindByName() err %v, wantErr %v", err, tt.wantErr)
			}
			if usrs, err := om.GetOrgUsers(tt.args.name); (err != nil) != tt.wantErr {
				t.Fatalf("GetOrgUsers() err %v, wantErr %v", err, tt.wantErr)
			} else if usrs != nil {
				var found bool
				for _, usr := range usrs {
					if usr == tt.args.owner {
						found = true
						break
					}
				}
				if !found {
					t.Fatal("failed to find user")
				}
			}
		})
	}
	// do a cleanup
	org, err := om.FindByName("testorg")
	if err != nil {
		t.Fatal(err)
	}
	om.DB.Unscoped().Delete(org)
}

func Test_BillingReport(t *testing.T) {
	var om = NewOrgManager(newTestDB(t, &Organization{}))
	om.DB.AutoMigrate(User{})
	om.DB.AutoMigrate(Upload{})
	om.DB.AutoMigrate(Usage{})
	// create the organization
	// create the organization
	if _, err := om.NewOrganization("testorg", "testorg-owner"); err != nil {
		t.Fatal(err)
	}
	org, err := om.FindByName("testorg")
	if err != nil {
		t.Fatal(err)
	}
	defer om.DB.Unscoped().Delete(org)
	// create an org user
	usr1, err := om.RegisterOrgUser(
		"testorg",
		"testorg-user1",
		"password123",
		"testorg-user1@example.org",
	)
	if err != nil {
		t.Fatal(err)
	}
	defer om.DB.Unscoped().Delete(usr1)
	if usr1.Organization != "testorg" {
		t.Fatal("bad organization set")
	}
	usage, err := NewUsageManager(om.DB).FindByUserName("testorg-user1")
	if err != nil {
		t.Fatal(err)
	}
	defer om.DB.Unscoped().Delete(usage)
	// create an upload
	upload, err := NewUploadManager(om.DB).NewUpload(
		"testhash", "upload", UploadOptions{
			NetworkName: "public",
			Username:    "testorg-user1",
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	defer om.DB.Unscoped().Delete(upload)
	report, err := om.GenerateBillingReport("testorg", time.Now().AddDate(0, 0, -30), time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Items) == 0 {
		t.Fatal("items length should be non 0")
	}
	if report.Time == 0 {
		t.Fatal("time should be non 0")
	}
	pp.Println(report)
}

func Test_AccountBalance(t *testing.T) {
	var om = NewOrgManager(newTestDB(t, &Organization{}))
	om.DB.AutoMigrate(User{})
	om.DB.AutoMigrate(Upload{})
	om.DB.AutoMigrate(Usage{})
	// create the organization
	// create the organization
	if _, err := om.NewOrganization("testorg", "testorg-owner"); err != nil {
		t.Fatal(err)
	}
	org, err := om.FindByName("testorg")
	if err != nil {
		t.Fatal(err)
	}
	defer om.DB.Unscoped().Delete(org)
	if err := om.IncreaseAmountOwed("testorg", 100); err != nil {
		t.Fatal(err)
	}
	org, err = om.FindByName("testorg")
	if org.AmountOwed != 100 {
		t.Fatal("bad account balance")
	}
	if err := om.DecreaseAmountOwed("testorg", 60.5); err != nil {
		t.Fatal(err)
	}
	org, err = om.FindByName("testorg")
	if err != nil {
		t.Fatal(err)
	}
	if org.AmountOwed != 39.5 {
		t.Fatal("bad account balance")
	}
	// now register an org user to test RemoveCredits updating balance
	usr, err := om.RegisterOrgUser(
		"testorg",
		"testorg-user33",
		"password123",
		"testorg-user33@example.org",
	)
	if err != nil {
		t.Fatal(err)
	}
	defer om.DB.Unscoped().Delete(usr)
	usg, err := NewUsageManager(om.DB).FindByUserName("testorg-user33")
	if err != nil {
		t.Fatal(err)
	}
	defer om.DB.Unscoped().Delete(usg)
	if _, err := NewUserManager(om.DB).RemoveCredits("testorg-user33", 0.5); err != nil {
		t.Fatal(err)
	}
	org, err = om.FindByName("testorg")
	if err != nil {
		t.Fatal(err)
	}
	if org.AmountOwed != 40 {
		t.Fatal("bad account balance")
	}
}

func Test_TotalStorageUsed(t *testing.T) {
	var om = NewOrgManager(newTestDB(t, &Organization{}))
	om.DB.AutoMigrate(User{})
	om.DB.AutoMigrate(Upload{})
	om.DB.AutoMigrate(Usage{})
	// create the organization
	if _, err := om.NewOrganization("testorg", "testorg-owner"); err != nil {
		t.Fatal(err)
	}
	org, err := om.FindByName("testorg")
	if err != nil {
		t.Fatal(err)
	}
	defer om.DB.Unscoped().Delete(org)
	usr, err := om.RegisterOrgUser(
		"testorg",
		"testorg-user",
		"password123",
		"testorg-user@example.org",
	)
	if err != nil {
		t.Fatal(err)
	}
	defer om.DB.Unscoped().Delete(usr)
	usg, err := NewUsageManager(om.DB).
		FindByUserName("testorg-user")
	if err != nil {
		t.Fatal(err)
	}
	defer om.DB.Unscoped().Delete(usg)
	if err := NewUsageManager(om.DB).UpdateDataUsage(
		"testorg-user",
		100,
	); err != nil {
		t.Fatal(err)
	}
	amount, err := om.GetTotalStorageUsed("testorg")
	if err != nil {
		t.Fatal(err)
	}
	if amount != 100 {
		t.Fatal("bad amount returned")
	}
}
