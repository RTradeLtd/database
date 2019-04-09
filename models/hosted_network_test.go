package models

import (
	"testing"
	"time"
)

func TestHostedNetworkManager_Access(t *testing.T) {
	var hm = NewHostedNetworkManager(newTestDB(t, &HostedNetwork{}))
	defer hm.DB.Close()
	if network, err := hm.CreateHostedPrivateNetwork(
		"myveryrandomnetworkname",
		"such swarm much protec",
		nil,
		NetworkAccessOptions{
			Owner: "testuserguy1",
			Users: []string{"testuserguy1", "testuserguy2"},
		},
	); err != nil {
		t.Fatal(err)
	} else if network.Owner != "testuserguy1" {
		t.Fatal("failed to correctly set network owner")
	}
}

func TestHostedNetworkManager_GetOfflineNetworks(t *testing.T) {
	var hm = NewHostedNetworkManager(newTestDB(t, &HostedNetwork{}))
	defer hm.DB.Close()

	hm.SaveNetwork(&HostedNetwork{
		Name: "online",
		Activated: func() *time.Time {
			var t = time.Now()
			return &t
		}(),
	})
	defer hm.Delete("online")

	var tnDisabled = &HostedNetwork{
		Name:     "disabled",
		Disabled: true,
	}
	hm.SaveNetwork(tnDisabled)
	defer hm.Delete(tnDisabled.Name)

	var tnEnabled = &HostedNetwork{
		Name:     "enabled",
		Disabled: false,
	}
	hm.SaveNetwork(tnEnabled)
	defer hm.Delete(tnEnabled.Name)

	type args struct {
		disabled bool
	}
	tests := []struct {
		name    string
		args    args
		want    []*HostedNetwork
		wantErr bool
	}{
		{"find disabled", args{true}, []*HostedNetwork{tnDisabled}, false},
		{"find enabled", args{false}, []*HostedNetwork{tnEnabled}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := hm.GetOfflineNetworks(tt.args.disabled)
			if (err != nil) != tt.wantErr {
				t.Errorf("HostedNetworkManager.GetOfflineNetworks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("expected %d entries, got %d", len(tt.want), len(got))
				return
			}
			for i, want := range tt.want {
				if got[i].Name != want.Name {
					t.Errorf("HostedNetworkManager.GetOfflineNetworks() = %s, want %s", got[i].Name, want.Name)
				}
			}
		})
	}
}
