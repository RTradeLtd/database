package models_test

import (
	"testing"

	"github.com/RTradeLtd/config"
	"github.com/RTradeLtd/database/models"
	"github.com/RTradeLtd/database/utils"
)

var (
	testNetwork = "test_network"
	testKeyName = "test_key_name"
	testKeyID   = "test_key_id"
	testCredits = float64(99999999)
)

type args struct {
	userName string
	email    string
	password string
}

func TestMigration_User(t *testing.T) {
	cfg, err := config.LoadConfig(testCfgPath)
	if err != nil {
		t.Fatal(err)
	}
	db, err := openDatabaseConnection(t, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if check := db.AutoMigrate(&models.User{}); check.Error != nil {
		t.Fatal(err)
	}
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
		randUtils = utils.GenerateRandomUtils()
		username  = randUtils.GenerateString(10, utils.LetterBytes)
		email     = randUtils.GenerateString(10, utils.LetterBytes)
	)

	tests := []struct {
		name string
		args args
	}{
		{"Success", args{username, email, "password123"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := um.NewUserAccount(tt.args.userName, tt.args.password, tt.args.email)
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
		randUtils = utils.GenerateRandomUtils()
		username  = randUtils.GenerateString(10, utils.LetterBytes)
		email     = randUtils.GenerateString(10, utils.LetterBytes)
	)

	tests := []struct {
		name string
		args args
	}{
		{"Success", args{username, email, "password123"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := um.NewUserAccount(tt.args.userName, tt.args.password, tt.args.email)
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

func TestUserManager_AddandRemoveIPFSNetworkForUSer(t *testing.T) {
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
		randUtils = utils.GenerateRandomUtils()
		username  = randUtils.GenerateString(10, utils.LetterBytes)
		email     = randUtils.GenerateString(10, utils.LetterBytes)
	)

	tests := []struct {
		name string
		args args
	}{
		{"Success", args{username, email, "password123"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := um.NewUserAccount(tt.args.userName, tt.args.password, tt.args.email)
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
			if err = um.RemoveIPFSNetworkForUser(
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
		randUtils = utils.GenerateRandomUtils()
		username  = randUtils.GenerateString(10, utils.LetterBytes)
		email     = randUtils.GenerateString(10, utils.LetterBytes)
	)

	tests := []struct {
		name string
		args args
	}{
		{"Success", args{username, email, "password123"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := um.NewUserAccount(tt.args.userName, tt.args.password, tt.args.email)
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
		randUtils = utils.GenerateRandomUtils()
		username  = randUtils.GenerateString(10, utils.LetterBytes)
		email     = randUtils.GenerateString(10, utils.LetterBytes)
	)

	tests := []struct {
		name string
		args args
	}{
		{"Success", args{username, email, "password123"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := um.NewUserAccount(tt.args.userName, tt.args.password, tt.args.email)
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
		randUtils = utils.GenerateRandomUtils()
		username  = randUtils.GenerateString(10, utils.LetterBytes)
		email     = randUtils.GenerateString(10, utils.LetterBytes)
	)

	tests := []struct {
		name string
		args args
	}{
		{"Success", args{username, email, "password123"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := um.NewUserAccount(tt.args.userName, tt.args.password, tt.args.email)
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
		randUtils = utils.GenerateRandomUtils()
		username  = randUtils.GenerateString(10, utils.LetterBytes)
		email     = randUtils.GenerateString(10, utils.LetterBytes)
	)

	tests := []struct {
		name string
		args args
	}{
		{"Success", args{username, email, "password123"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := um.NewUserAccount(tt.args.userName, tt.args.password, tt.args.email)
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
		randUtils = utils.GenerateRandomUtils()
		username  = randUtils.GenerateString(10, utils.LetterBytes)
		email     = randUtils.GenerateString(10, utils.LetterBytes)
	)

	tests := []struct {
		name string
		args args
	}{
		{"Success", args{username, email, "password123"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := um.NewUserAccount(tt.args.userName, tt.args.password, tt.args.email)
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
		randUtils = utils.GenerateRandomUtils()
		username  = randUtils.GenerateString(10, utils.LetterBytes)
		email     = randUtils.GenerateString(10, utils.LetterBytes)
	)

	tests := []struct {
		name string
		args args
	}{
		{"Success", args{username, email, "password123"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := um.NewUserAccount(tt.args.userName, tt.args.password, tt.args.email)
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
		randUtils = utils.GenerateRandomUtils()
		username  = randUtils.GenerateString(10, utils.LetterBytes)
		email     = randUtils.GenerateString(10, utils.LetterBytes)
	)

	tests := []struct {
		name string
		args args
	}{
		{"Success", args{username, email, "password123"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := um.NewUserAccount(tt.args.userName, tt.args.password, tt.args.email)
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
		randUtils = utils.GenerateRandomUtils()
		username  = randUtils.GenerateString(10, utils.LetterBytes)
		email     = randUtils.GenerateString(10, utils.LetterBytes)
	)

	tests := []struct {
		name string
		args args
	}{
		{"Success", args{username, email, "password123"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := um.NewUserAccount(tt.args.userName, tt.args.password, tt.args.email)
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
		randUtils = utils.GenerateRandomUtils()
		username  = randUtils.GenerateString(10, utils.LetterBytes)
		email     = randUtils.GenerateString(10, utils.LetterBytes)
	)

	tests := []struct {
		name string
		args args
	}{
		{"Success", args{username, email, "password123"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := um.NewUserAccount(tt.args.userName, tt.args.password, tt.args.email)
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
		randUtils = utils.GenerateRandomUtils()
		username  = randUtils.GenerateString(10, utils.LetterBytes)
		email     = randUtils.GenerateString(10, utils.LetterBytes)
	)

	tests := []struct {
		name string
		args args
	}{
		{"Success", args{username, email, "password123"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := um.NewUserAccount(tt.args.userName, tt.args.password, tt.args.email)
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
		randUtils = utils.GenerateRandomUtils()
		username  = randUtils.GenerateString(10, utils.LetterBytes)
		email     = randUtils.GenerateString(10, utils.LetterBytes)
	)

	tests := []struct {
		name string
		args args
	}{
		{"Success", args{username, email, "password123"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := um.NewUserAccount(tt.args.userName, tt.args.password, tt.args.email)
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
		randUtils = utils.GenerateRandomUtils()
		username  = randUtils.GenerateString(10, utils.LetterBytes)
		email     = randUtils.GenerateString(10, utils.LetterBytes)
	)

	tests := []struct {
		name string
		args args
	}{
		{"Success", args{username, email, "password123"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := um.NewUserAccount(tt.args.userName, tt.args.password, tt.args.email)
			if err != nil {
				t.Fatal(err)
			}
			defer um.DB.Delete(user)
			userCopy, err := um.AddCredits(
				tt.args.userName,
				float64(1),
			)
			if err != nil {
				t.Fatal(err)
			}
			if userCopy.Credits != testCredits+1 {
				t.Fatal("failed to add credits")
			}
			credits, err := um.GetCreditsForUser(
				tt.args.userName,
			)
			if err != nil {
				t.Fatal(err)
			}
			if credits != testCredits+1 {
				t.Fatal("failed to get credits")
			}
			userCopy, err = um.RemoveCredits(
				tt.args.userName,
				testCredits,
			)
			if err != nil {
				t.Fatal(err)
			}
			if userCopy.Credits != 1 {
				t.Fatal("failed to remove credits")
			}
		})
	}
}
func TestUserManager_ResetPassword(t *testing.T) {
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
		randUtils = utils.GenerateRandomUtils()
		username  = randUtils.GenerateString(10, utils.LetterBytes)
		email     = randUtils.GenerateString(10, utils.LetterBytes)
	)

	tests := []struct {
		name string
		args args
	}{
		{"Success", args{username, email, "password123"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := um.NewUserAccount(tt.args.userName, tt.args.password, tt.args.email)
			if err != nil {
				t.Fatal(err)
			}
			defer um.DB.Delete(user)
			newPass, err := um.ResetPassword(tt.args.userName)
			if err != nil {
				t.Fatal(err)
			}
			if newPass == "" {
				t.Fatal("failed to reset password")
			}
		})
	}
}
