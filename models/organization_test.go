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
		{"Pass", args{"testorg", "mytestuser"}, false},
		{"Fail", args{"testorg", "mytestuser"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := om.NewOrganization(
				tt.args.name,
				tt.args.owner,
			); (err != nil) != tt.wantErr {
				t.Fatalf("NewOrganization() err %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
