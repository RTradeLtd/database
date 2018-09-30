package models_test

import (
	"fmt"
	"testing"

	"github.com/RTradeLtd/config"
	"github.com/RTradeLtd/database/models"
	"github.com/RTradeLtd/database/utils"
	"github.com/jinzhu/gorm"
)

var (
	testCfgPath = "../test/config.json"
	testNetwork = "test_network"
	testKeyName = "test_key_name"
	testKeyID   = "test_key_id"
	testCredits = 10.5
)

type args struct {
	ethAddress        string
	userName          string
	email             string
	password          string
	enterpriseEnabled bool
}

func TestUserManager_GetPrivateIPFSNetworksForUSer(t *testing.T) {
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
		{"Test1", args{ethAddress, username, email, "password123", false}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := um.NewUserAccount(tt.args.ethAddress, tt.args.userName, tt.args.password, tt.args.email, tt.args.enterpriseEnabled)
			if err != nil {
				t.Fatal(err)
			}
			defer um.DB.Delete(user)
			if _, err := um.GetPrivateIPFSNetworksForUser(tt.args.userName); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestUserManager_CheckIfUserHasAccessToNetwork(t *testing.T) {
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
		{"Test1", args{ethAddress, username, email, "password123", false}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := um.NewUserAccount(tt.args.ethAddress, tt.args.userName, tt.args.password, tt.args.email, tt.args.enterpriseEnabled)
			if err != nil {
				t.Fatal(err)
			}
			defer um.DB.Delete(user)
			access, err := um.CheckIfUserHasAccessToNetwork(
				tt.args.userName,
				testNetwork,
			)
			if err != nil {
				t.Fatal("err")
			}
			if access {
				t.Fatal("access to non existent network, this should not happen")
			}
		})
	}
}

func TestUserManager_AddIPFSNetworkForUSer(t *testing.T) {
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
		{"Test1", args{ethAddress, username, email, "password123", false}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := um.NewUserAccount(tt.args.ethAddress, tt.args.userName, tt.args.password, tt.args.email, tt.args.enterpriseEnabled)
			if err != nil {
				t.Fatal(err)
			}
			defer um.DB.Delete(user)
			if err := um.AddIPFSNetworkForUser(
				tt.args.userName,
				testNetwork,
			); err != nil {
				t.Fatal(err)
			}
			if _, err := um.CheckIfUserHasAccessToNetwork(
				tt.args.userName,
				testNetwork,
			); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestUserManager_AddIPFSKeyForUser(t *testing.T) {
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
		{"Test1", args{ethAddress, username, email, "password123", false}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := um.NewUserAccount(tt.args.ethAddress, tt.args.userName, tt.args.password, tt.args.email, tt.args.enterpriseEnabled)
			if err != nil {
				t.Fatal(err)
			}
			defer um.DB.Delete(user)
			if err := um.AddIPFSKeyForUser(
				tt.args.userName,
				testKeyName,
				testKeyID,
			); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestUserManager_GetKeysForUser(t *testing.T) {
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
		{"Test1", args{ethAddress, username, email, "password123", false}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := um.NewUserAccount(tt.args.ethAddress, tt.args.userName, tt.args.password, tt.args.email, tt.args.enterpriseEnabled)
			if err != nil {
				t.Fatal(err)
			}
			defer um.DB.Delete(user)
			if _, err := um.GetKeysForUser(tt.args.userName); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestUserManager_GetKeyIDByName(t *testing.T) {
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
		{"Test1", args{ethAddress, username, email, "password123", false}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := um.NewUserAccount(tt.args.ethAddress, tt.args.userName, tt.args.password, tt.args.email, tt.args.enterpriseEnabled)
			if err != nil {
				t.Fatal(err)
			}
			defer um.DB.Delete(user)
			if err := um.AddIPFSKeyForUser(
				tt.args.userName,
				testKeyName,
				testKeyID,
			); err != nil {
				t.Fatal(err)
			}
			if _, err := um.GetKeyIDByName(
				tt.args.userName,
				testKeyName,
			); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestUserManager_CheckIfKeyOwnedByUser(t *testing.T) {
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
		{"Test1", args{ethAddress, username, email, "password123", false}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := um.NewUserAccount(tt.args.ethAddress, tt.args.userName, tt.args.password, tt.args.email, tt.args.enterpriseEnabled)
			if err != nil {
				t.Fatal(err)
			}
			defer um.DB.Delete(user)
			if err := um.AddIPFSKeyForUser(
				tt.args.userName,
				testKeyName,
				testKeyID,
			); err != nil {
				t.Fatal(err)
			}
			if valid, err := um.CheckIfKeyOwnedByUser(
				tt.args.userName,
				testKeyName,
			); err != nil {
				t.Fatal(err)
			} else if !valid {
				t.Fatal("no error returned, but user does not own key, this is unexpected")
			}
		})
	}
}

func TestUserManager_CheckIfAccountEnabled(t *testing.T) {
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
		{"Test1", args{ethAddress, username, email, "password123", false}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := um.NewUserAccount(tt.args.ethAddress, tt.args.userName, tt.args.password, tt.args.email, tt.args.enterpriseEnabled)
			if err != nil {
				t.Fatal(err)
			}
			defer um.DB.Delete(user)
			enabled, err := um.CheckIfUserAccountEnabled(tt.args.userName)
			if err != nil {
				t.Fatal(err)
			}
			if !enabled {
				t.Fatal("user account is not enabled")
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
		{"Test1", args{ethAddress, username, email, "password123", false}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := um.NewUserAccount(tt.args.ethAddress, tt.args.userName, tt.args.password, tt.args.email, tt.args.enterpriseEnabled)
			if err != nil {
				t.Fatal(err)
			}
			defer um.DB.Delete(user)
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
		{"Test1", args{ethAddress, username, email, "password123", false}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := um.NewUserAccount(tt.args.ethAddress, tt.args.userName, tt.args.password, tt.args.email, tt.args.enterpriseEnabled)
			if err != nil {
				t.Fatal(err)
			}
			um.DB.Delete(user)
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
		{"Test1", args{ethAddress, username, email, "password123", false}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := um.NewUserAccount(tt.args.ethAddress, tt.args.userName, tt.args.password, tt.args.email, tt.args.enterpriseEnabled)
			if err != nil {
				t.Fatal(err)
			}
			defer um.DB.Delete(user)
			valid, err := um.SignIn(
				tt.args.userName,
				tt.args.password,
			)
			if err != nil {
				t.Fatal(err)
			}
			if !valid {
				t.Fatal("failed to sign user in")
			}
		})
	}
}

func TestUserManager_ComparePlaintextPasswordToHash(t *testing.T) {
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
		{"Test1", args{ethAddress, username, email, "password123", false}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := um.NewUserAccount(tt.args.ethAddress, tt.args.userName, tt.args.password, tt.args.email, tt.args.enterpriseEnabled)
			if err != nil {
				t.Fatal(err)
			}
			defer um.DB.Delete(user)
			valid, err := um.ComparePlaintextPasswordToHash(
				tt.args.userName,
				tt.args.password,
			)
			if err != nil {
				t.Fatal(err)
			}
			if !valid {
				t.Fatal("failed to compare plaintext pass to hash")
			}
		})
	}
}

func TestUserManager_FindByAddress(t *testing.T) {
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
		{"Test1", args{ethAddress, username, email, "password123", false}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := um.NewUserAccount(tt.args.ethAddress, tt.args.userName, tt.args.password, tt.args.email, tt.args.enterpriseEnabled)
			if err != nil {
				t.Fatal(err)
			}
			defer um.DB.Delete(user)
			userCopy, err := um.FindByAddress(
				tt.args.ethAddress,
			)
			if err != nil {
				t.Fatal(err)
			}
			if userCopy.UserName != user.UserName {
				t.Fatal("failed to find correct account")
			}
		})
	}
}

func TestUserManager_FindEthAddressByUserName(t *testing.T) {
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
		{"Test1", args{ethAddress, username, email, "password123", false}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := um.NewUserAccount(tt.args.ethAddress, tt.args.userName, tt.args.password, tt.args.email, tt.args.enterpriseEnabled)
			if err != nil {
				t.Fatal(err)
			}
			defer um.DB.Delete(user)
			address, err := um.FindEthAddressByUserName(
				tt.args.userName,
			)
			if err != nil {
				t.Fatal(err)
			}
			if address != user.EthAddress {
				t.Fatal("failed to find correct account")
			}
		})
	}
}

func TestUserManager_FindEmailByUserName(t *testing.T) {
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
		{"Test1", args{ethAddress, username, email, "password123", false}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := um.NewUserAccount(tt.args.ethAddress, tt.args.userName, tt.args.password, tt.args.email, tt.args.enterpriseEnabled)
			if err != nil {
				t.Fatal(err)
			}
			defer um.DB.Delete(user)
			if _, err := um.FindEmailByUserName(tt.args.userName); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestUserManager_FindUserByUserName(t *testing.T) {
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
		{"Test1", args{ethAddress, username, email, "password123", false}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := um.NewUserAccount(tt.args.ethAddress, tt.args.userName, tt.args.password, tt.args.email, tt.args.enterpriseEnabled)
			if err != nil {
				t.Fatal(err)
			}
			defer um.DB.Delete(user)
			userCopy, err := um.FindByUserName(
				tt.args.userName,
			)
			if err != nil {
				t.Fatal(err)
			}
			if userCopy.UserName != user.UserName {
				t.Fatal("failed to find correct account")
			}
		})
	}
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
		{"Test1", args{ethAddress, username, email, "password123", false}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := um.NewUserAccount(tt.args.ethAddress, tt.args.userName, tt.args.password, tt.args.email, tt.args.enterpriseEnabled)
			if err != nil {
				t.Fatal(err)
			}
			defer um.DB.Delete(user)
			if _, err := um.ChangeEthereumAddress(
				tt.args.userName,
				tt.args.ethAddress,
			); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestUserManager_Credits(t *testing.T) {
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
		{"Test1", args{ethAddress, username, email, "password123", false}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := um.NewUserAccount(tt.args.ethAddress, tt.args.userName, tt.args.password, tt.args.email, tt.args.enterpriseEnabled)
			if err != nil {
				t.Fatal(err)
			}
			defer um.DB.Delete(user)
			userCopy, err := um.AddCreditsForUser(
				tt.args.userName,
				testCredits,
			)
			if err != nil {
				t.Fatal(err)
			}
			if userCopy.Credits != testCredits {
				t.Fatal("failed to add credits")
			}
			credits, err := um.GetCreditsForUser(
				tt.args.userName,
			)
			if err != nil {
				t.Fatal(err)
			}
			if credits != testCredits {
				t.Fatal("failed to get credits")
			}
			userCopy, err = um.RemoveCredits(
				tt.args.userName,
				testCredits,
			)
			if err != nil {
				t.Fatal(err)
			}
			if userCopy.Credits != 0 {
				t.Fatal("failed to remove credits")
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
