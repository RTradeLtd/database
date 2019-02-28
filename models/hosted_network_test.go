package models

import (
	"reflect"
	"testing"
	"time"
)

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
		want    []HostedNetwork
		wantErr bool
	}{
		{"find disabled", args{true}, []HostedNetwork{*tnDisabled}, false},
		{"find enabled", args{false}, []HostedNetwork{*tnEnabled}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := hm.GetOfflineNetworks(tt.args.disabled)
			if (err != nil) != tt.wantErr {
				t.Errorf("HostedNetworkManager.GetAllOfflineNetworks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HostedNetworkManager.GetAllOfflineNetworks() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
