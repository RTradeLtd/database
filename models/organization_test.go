package models

import (
	"fmt"
	"testing"
)

func TestOrganizationManager_Full(t *testing.T) {
	var om = NewOrgManager(newTestDB(t, &Organization{}))

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
			if err := om.NewOrganization(
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
	// create the organization
	if err := om.NewOrganization("testorg", "testorg-owner"); err != nil {
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
	report, err := om.GenerateBillingReport("testorg")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", report)
}
