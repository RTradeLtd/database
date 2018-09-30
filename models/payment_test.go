package models_test

import (
	"testing"

	"github.com/RTradeLtd/config"
	"github.com/RTradeLtd/database/models"
)

func TestPaymentManager_NewPayment(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	cfg, err := config.LoadConfig(testCfgPath)
	if err != nil {
		t.Fatal(err)
	}

	db, err := openDatabaseConnection(t, cfg)
	if err != nil {
		t.Fatal(err)
	}
	pm := models.NewPaymentManager(db)
	type args struct {
		depositAddress string
		txHash         string
		usdValue       float64
		blockchain     string
		paymentType    string
		username       string
	}
	tests := []struct {
		name string
		args args
	}{
		{"Payment1", args{"depositAddress", "txHash", 0.124, "blockchain", "paymentType", "userName"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payment, err := pm.NewPayment(
				tt.args.depositAddress,
				tt.args.txHash,
				tt.args.usdValue,
				tt.args.blockchain,
				tt.args.paymentType,
				tt.args.username,
			)
			if err != nil {
				t.Fatal(err)
			}
			pm.DB.Delete(payment)
		})
	}
}

func TestPaymentManager_ConfirmPayment(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	cfg, err := config.LoadConfig(testCfgPath)
	if err != nil {
		t.Fatal(err)
	}

	db, err := openDatabaseConnection(t, cfg)
	if err != nil {
		t.Fatal(err)
	}
	pm := models.NewPaymentManager(db)
	type args struct {
		depositAddress string
		txHash         string
		usdValue       float64
		blockchain     string
		paymentType    string
		username       string
	}
	tests := []struct {
		name string
		args args
	}{
		{"Payment1", args{"depositAddress", "txHash", 0.124, "blockchain", "paymentType", "userName"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payment, err := pm.NewPayment(
				tt.args.depositAddress,
				tt.args.txHash,
				tt.args.usdValue,
				tt.args.blockchain,
				tt.args.paymentType,
				tt.args.username,
			)
			if err != nil {
				t.Fatal(err)
			}
			defer pm.DB.Delete(payment)
			paymentCopy, err := pm.ConfirmPayment(payment.TxHash)
			if err != nil {
				t.Fatal(err)
			}
			if paymentCopy.ID != payment.ID {
				t.Fatal("bad payment recovered")
			}
		})
	}
}

func TestPaymentManager_MarkAsProcessing(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	cfg, err := config.LoadConfig(testCfgPath)
	if err != nil {
		t.Fatal(err)
	}

	db, err := openDatabaseConnection(t, cfg)
	if err != nil {
		t.Fatal(err)
	}
	pm := models.NewPaymentManager(db)
	type args struct {
		depositAddress string
		txHash         string
		usdValue       float64
		blockchain     string
		paymentType    string
		username       string
	}
	tests := []struct {
		name string
		args args
	}{
		{"Payment1", args{"depositAddress", "txHash", 0.124, "blockchain", "paymentType", "userName"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payment, err := pm.NewPayment(
				tt.args.depositAddress,
				tt.args.txHash,
				tt.args.usdValue,
				tt.args.blockchain,
				tt.args.paymentType,
				tt.args.username,
			)
			if err != nil {
				t.Fatal(err)
			}
			defer pm.DB.Delete(payment)
			paymentCopy, err := pm.MarkPaymentAsProcessing(payment.TxHash)
			if err != nil {
				t.Fatal(err)
			}
			if paymentCopy.ID != payment.ID {
				t.Fatal("bad payment recovered")
			}
		})
	}
}
