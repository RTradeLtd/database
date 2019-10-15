package models

import "testing"

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
			if usrs, err := om.GetOrgUsers(tt.args.name) (err != nil) != tt.wantErr {
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
