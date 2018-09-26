package models_test

import (
	"fmt"
	"testing"

	"github.com/RTradeLtd/Temporal/config"
	"github.com/RTradeLtd/Temporal/models"
	"github.com/RTradeLtd/Temporal/utils"
	"github.com/jinzhu/gorm"
)

var (
	testCfgPath = "../test/config.json"
)

type args struct {
	ethAddress        string
	userName          string
	email             string
	password          string
	enterpriseEnabled bool
}

func TestUserManager_ChangeEthereumAddress(t *testing.T) {
	cfg, err := config.LoadConfig(testCfgPath)
	if err != nil {
		t.Fatal(err)
	}

	db, err := openDatabaseConnection(t, cfg)
	if err != nil {
		t.Fatal(err)
	}
	um := models.NewUserManager(db)

	var (
		randUtils  = utils.GenerateRandomUtils()
		username   = randUtils.GenerateString(10, utils.LetterBytes)
		ethAddress = randUtils.GenerateString(10, utils.LetterBytes)
		email      = randUtils.GenerateString(10, utils.LetterBytes)
	)

	tests := []struct {
		name string
		args args
	}{
		{"ChangeEthereumAddress", args{ethAddress, username, email, "password123", false}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := um.NewUserAccount(tt.args.ethAddress, tt.args.userName, tt.args.password, tt.args.email, tt.args.enterpriseEnabled); err != nil {
				t.Fatal(err)
			}
			if _, err := um.ChangeEthereumAddress(tt.args.userName, tt.args.ethAddress); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestUserManager_ChangePassword(t *testing.T) {
	cfg, err := config.LoadConfig(testCfgPath)
	if err != nil {
		t.Fatal(err)
	}

	db, err := openDatabaseConnection(t, cfg)
	if err != nil {
		t.Fatal(err)
	}
	um := models.NewUserManager(db)

	var (
		randUtils  = utils.GenerateRandomUtils()
		username   = randUtils.GenerateString(10, utils.LetterBytes)
		ethAddress = randUtils.GenerateString(10, utils.LetterBytes)
		email      = randUtils.GenerateString(10, utils.LetterBytes)
	)

	tests := []struct {
		name string
		args args
	}{
		{"ChangePassword", args{ethAddress, username, email, "password123", false}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := um.NewUserAccount(tt.args.ethAddress, tt.args.userName, tt.args.password, tt.args.email, tt.args.enterpriseEnabled); err != nil {
				t.Fatal(err)
			}
			changed, err := um.ChangePassword(tt.args.userName, tt.args.password, "newpassword")
			if err != nil {
				t.Fatal(err)
			}
			if !changed {
				t.Error("password changed failed, but no error occured")
			}
		})
	}
}

func TestUserManager_NewAccount(t *testing.T) {
	cfg, err := config.LoadConfig(testCfgPath)
	if err != nil {
		t.Fatal(err)
	}

	db, err := openDatabaseConnection(t, cfg)
	if err != nil {
		t.Fatal(err)
	}
	um := models.NewUserManager(db)

	var (
		randUtils  = utils.GenerateRandomUtils()
		username   = randUtils.GenerateString(10, utils.LetterBytes)
		ethAddress = randUtils.GenerateString(10, utils.LetterBytes)
		email      = randUtils.GenerateString(10, utils.LetterBytes)
	)

	tests := []struct {
		name string
		args args
	}{
		{"AccountCreation", args{ethAddress, username, email, "password123", false}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := um.NewUserAccount(tt.args.ethAddress, tt.args.userName, tt.args.password, tt.args.email, tt.args.enterpriseEnabled); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestUserManager_SignIn(t *testing.T) {
	cfg, err := config.LoadConfig(testCfgPath)
	if err != nil {
		t.Fatal(err)
	}

	db, err := openDatabaseConnection(t, cfg)
	if err != nil {
		t.Fatal(err)
	}
	um := models.NewUserManager(db)

	var (
		randUtils  = utils.GenerateRandomUtils()
		username   = randUtils.GenerateString(10, utils.LetterBytes)
		ethAddress = randUtils.GenerateString(10, utils.LetterBytes)
		email      = randUtils.GenerateString(10, utils.LetterBytes)
	)

	tests := []struct {
		name string
		args args
	}{
		{"AccountCreation", args{ethAddress, username, email, "password123", false}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := um.NewUserAccount(tt.args.ethAddress, tt.args.userName, tt.args.password, tt.args.email, tt.args.enterpriseEnabled); err != nil {
				t.Fatal(err)
			}
			success, err := um.SignIn(tt.args.userName, tt.args.password)
			if err != nil {
				t.Fatal(err)
			}
			if !success {
				t.Error("sign in failed")
			}
		})
	}
}

func openDatabaseConnection(t *testing.T, cfg *config.TemporalConfig) (*gorm.DB, error) {
	dbConnURL := fmt.Sprintf("host=127.0.0.1 port=5433 user=postgres dbname=temporal password=%s sslmode=disable",
		cfg.Database.Password)

	db, err := gorm.Open("postgres", dbConnURL)
	if err != nil {
		t.Fatal(err)
	}
	return db, nil
}
