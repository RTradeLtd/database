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
			if usr, err := om.RegisterOrgUser(
				tt.args.name, tt.args.owner, "password123", "password123@example.org",
			); (err != nil) != tt.wantErr {
				t.Fatalf("RegisterOrgUser err %v, wantErr %v", err, tt.wantErr)
			} else if usr != nil {
				defer om.DB.Unscoped().Delete(usr)
			}
			if model, err := om.FindByName(tt.args.name); (err != nil) != tt.wantErr {
				t.Fatalf("FindByName() err %v, wantErr %v", err, tt.wantErr)
			} else if model != nil {
				defer om.DB.Unscoped().Delete(model)
			}
		})
	}
}
