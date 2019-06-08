package models

import (
	"testing"
)

var (
	testNetwork = "test_network"
	testKeyName = "test_key_name"
	testKeyID   = "test_key_id"
	testCredits = float64(0)
	username    = "muchuserverywow"
	email       = "muchemailverysmtp@gmail.com"
	password    = "password123"
)

type args struct {
	userName string
	email    string
	password string
}

func TestUserManager_NewAccount(t *testing.T) {
	var um = NewUserManager(newTestDB(t, &User{}))
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Success", args{username, email, "password123"}, false},
		{"Failure-UsedEmail", args{"randomuserbro", email, "password123"}, true},
		{"Failure-UsedUserName", args{username, email, "password123"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := um.NewUserAccount(tt.args.userName, tt.args.password, tt.args.email); (err != nil) != tt.wantErr {
				t.Fatalf("NewUserAccount err = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserManager_GetPrivateIPFSNetworksForUSer(t *testing.T) {
	var um = NewUserManager(newTestDB(t, &User{}))
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Success", args{username, email, "password123"}, false},
		{"Failure", args{"notarealuser", "notarealemail", "password123"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// add a private network for testing purposes
			if err := um.AddIPFSNetworkForUser(
				tt.args.userName,
				"thisisdefinitelynotgoingtobarealnamedude",
			); (err != nil) != tt.wantErr {
				t.Fatalf("AddIPFSNetworkForUser err = %v, wantErr %v", err, tt.wantErr)
			}
			if _, err := um.GetPrivateIPFSNetworksForUser(tt.args.userName); (err != nil) != tt.wantErr {
				t.Fatalf("GetPrivateIPFSNetworksForUser() wantErr = %v, error = %v", tt.wantErr, err.Error())
			}
		})
	}
}

func TestUserManager_CheckIfUserHasAccessToNetwork(t *testing.T) {
	var um = NewUserManager(newTestDB(t, &User{}))
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Success", args{username, email, "password123"}, false},
		{"Failure", args{"notarealuser", "notarealemail", "password123"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := um.CheckIfUserHasAccessToNetwork(tt.args.userName, testNetwork); (err != nil) != tt.wantErr {
				t.Fatalf("CheckIfUserHasAccessToNetwork() wantErr = %v, error = %v", tt.wantErr, err.Error())
			}
		})
	}
}

func TestUserManager_AddandRemoveIPFSNetworkForUSer(t *testing.T) {
	var um = NewUserManager(newTestDB(t, &User{}))
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Success", args{username, email, "password123"}, false},
		{"Failure", args{"notarealuser", "notarealuser", "notarealuser"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := um.AddIPFSNetworkForUser(
				tt.args.userName,
				testNetwork,
			); (err != nil) != tt.wantErr {
				t.Fatalf("AddIPFSNetworkForUser err = %v, want %v", err, tt.wantErr)
			}
			if _, err := um.CheckIfUserHasAccessToNetwork(
				tt.args.userName,
				testNetwork,
			); (err != nil) != tt.wantErr {
				t.Fatalf("CheckIfUserHasAccessToNetwork err = %v, want %v", err, tt.wantErr)
			}
			if err := um.RemoveIPFSNetworkForUser(
				tt.args.userName,
				testNetwork,
			); (err != nil) != tt.wantErr {
				t.Fatalf("RemoveIPFSNetworkForUser err = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserManager_AddIPFSKeyForUser(t *testing.T) {
	var um = NewUserManager(newTestDB(t, &User{}))
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Success", args{username, email, "password123"}, false},
		{"Failure", args{"notarealuser", "notarealemail", "password123"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := um.AddIPFSKeyForUser(
				tt.args.userName,
				testKeyName,
				testKeyID,
			); (err != nil) != tt.wantErr {
				t.Fatalf("AddIPFSKeyForUser err = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserManager_GetKeysForUser(t *testing.T) {
	var um = NewUserManager(newTestDB(t, &User{}))
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Success", args{username, email, "password123"}, false},
		{"Failure", args{"notarealuser", "notarealuser", "password123"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := um.GetKeysForUser(tt.args.userName); (err != nil) != tt.wantErr {
				t.Fatalf("GetKeysForUser() err = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserManager_GetKeyIDByName(t *testing.T) {
	var um = NewUserManager(newTestDB(t, &User{}))
	type args struct {
		userName string
		email    string
		password string
		keyName  string
		keyID    string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Success", args{username, email, "password123", "randomkeya", "randomkeyaid"}, false},
		{"Failure", args{"notarealuser", "notarealuser", "password123", "blah", "blah"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := um.AddIPFSKeyForUser(
				tt.args.userName,
				tt.args.keyName,
				tt.args.keyID,
			); (err != nil) != tt.wantErr {
				t.Fatalf("AddIPFSKeyForUser err = %v, wantErr %v", err, tt.wantErr)
			}
			if _, err := um.GetKeyIDByName(
				tt.args.userName,
				tt.args.keyName,
			); (err != nil) != tt.wantErr {
				t.Fatalf("GetKeyIDByName err = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserManager_CheckIfKeyOwnedByUser(t *testing.T) {
	var um = NewUserManager(newTestDB(t, &User{}))
	type args struct {
		userName string
		email    string
		password string
		keyName  string
		keyID    string
	}
	tests := []struct {
		name      string
		args      args
		wantValid bool
		wantErr   bool
	}{
		{"Success", args{username, email, "password123", "randomkeyb", "randomkeybid"}, true, false},
		{"Failure", args{"notarealuser", "notarealuser", "password123", "blah", "blah"}, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := um.AddIPFSKeyForUser(
				tt.args.userName,
				tt.args.keyName,
				tt.args.keyID,
			); (err != nil) != tt.wantErr {
				t.Fatalf("AddIPFSKeyForUser err = %v, wantErr %v", err, tt.wantErr)
			}
			if valid, err := um.CheckIfKeyOwnedByUser(
				tt.args.userName,
				tt.args.keyName,
			); (err != nil) != tt.wantErr {
				t.Fatalf("CheckIfKeyOwnedByUser err = %v, wantErr %v", err, tt.wantErr)
			} else if valid != tt.wantValid {
				t.Fatalf("CheckIfKeyOwnedByUser valid = %v, wantValid %v", valid, tt.wantValid)
			}
		})
	}
}

func TestUserManager_CheckIfAccountEnabled(t *testing.T) {
	var um = NewUserManager(newTestDB(t, &User{}))
	tests := []struct {
		name        string
		args        args
		wantEnabled bool
		wantErr     bool
	}{
		{"Success", args{username, email, "password123"}, true, false},
		{"Failure", args{"notarealuser", "notarealuser", "password123"}, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enabled, err := um.CheckIfUserAccountEnabled(tt.args.userName)
			if (err != nil) != tt.wantErr {
				t.Fatalf("CheckIfUserAccountEnabled err = %v, wantErr %v", err, tt.wantErr)
			}
			if enabled != tt.wantEnabled {
				t.Fatalf("CheckIfUserAccountEnabled enabled = %v, wantEnabled %v", enabled, tt.wantEnabled)
			}
		})
	}
}

func TestUserManager_ChangePassword(t *testing.T) {
	var um = NewUserManager(newTestDB(t, &User{}))
	tests := []struct {
		name        string
		args        args
		wantChanged bool
		wantErr     bool
	}{
		{"Success", args{username, email, "password123"}, true, false},
		{"Failure", args{"notarealuser", "notarealuser", "password123"}, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			changed, err := um.ChangePassword(tt.args.userName, tt.args.password, "password123")
			if (err != nil) != tt.wantErr {
				t.Fatalf("ChangePassword err = %v, wantErr %v", err, tt.wantErr)
			}
			if changed != tt.wantChanged {
				t.Errorf("ChangePassword change = %v, wantChange %v", changed, tt.wantChanged)
			}
		})
	}
}

func TestUserManager_SignIn(t *testing.T) {
	var um = NewUserManager(newTestDB(t, &User{}))
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		wantValid bool
	}{
		{"Success", args{username, email, "password123"}, false, true},
		{"Failure", args{"notarealuser", "notarealemail", "password123"}, true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if valid, err := um.SignIn(
				tt.args.userName,
				tt.args.password,
			); (err != nil) != tt.wantErr {
				t.Fatalf("SignIn() wantErr = %v, error = %v", tt.wantErr, err.Error())
			} else if valid != tt.wantValid {
				t.Fatalf("SignIn() wantValid = %v, valid = %v", tt.wantValid, valid)
			}

			if valid, err := um.SignIn(
				tt.args.email,
				tt.args.password,
			); (err != nil) != tt.wantErr {
				t.Fatalf("SignIn() wantErr = %v, error = %v", tt.wantErr, err.Error())
			} else if valid != tt.wantValid {
				t.Fatalf("SignIn() wantValid = %v, valid = %v", tt.wantValid, valid)
			}
		})
	}
}

func TestUserManager_ComparePlaintextPasswordToHash(t *testing.T) {
	var um = NewUserManager(newTestDB(t, &User{}))
	tests := []struct {
		name      string
		args      args
		wantValid bool
		wantErr   bool
	}{
		{"Success", args{username, email, "password123"}, true, false},
		{"Failure", args{"notarealuser", "NotarealEmail", "password123"}, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, err := um.ComparePlaintextPasswordToHash(
				tt.args.userName,
				tt.args.password,
			)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ComparePlaintextPasswordToHash err = %v, wantErr %v", err, tt.wantErr)
			}
			if valid != tt.wantValid {
				t.Fatalf("ComparePlaintextPasswordToHash valid = %v, wantValid %v", valid, tt.wantValid)
			}
		})
	}
}

func TestUserManager_FindUserByUserName(t *testing.T) {
	var um = NewUserManager(newTestDB(t, &User{}))
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Success", args{username, email, "password123"}, false},
		{"Failure", args{"notarealuser", "notarealuser", "password123"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if user, err := um.FindByUserName(
				tt.args.userName,
			); (err != nil) != tt.wantErr {
				t.Fatalf("FindByUserName err = %v, wantErr %v", err, tt.wantErr)
			} else if (err != nil) == false && user.UserName != tt.args.userName {
				t.Fatal("failed to find correct username")
			}
		})
	}
}

func TestUserManager_Credits(t *testing.T) {
	var um = NewUserManager(newTestDB(t, &User{}))
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Success", args{username, email, "password123"}, false},
		{"Failure", args{"notarealuser", "notarealemail", "password123"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userCopy, err := um.AddCredits(
				tt.args.userName,
				float64(1),
			)
			if (err != nil) != tt.wantErr {
				t.Fatalf("AddCredits err = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && userCopy.Credits != testCredits+1 {
				t.Fatal("failed to add credits")
			}
			credits, err := um.GetCreditsForUser(
				tt.args.userName,
			)
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetCreditsForUser err = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && credits != testCredits+1 {
				t.Fatal("failed to get credits")
			}
			userCopy, err = um.RemoveCredits(
				tt.args.userName,
				testCredits,
			)
			if (err != nil) != tt.wantErr {
				t.Fatalf("RemoveCredits err = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && userCopy.Credits != 1 {
				t.Fatal("failed to remove credits")
			}
		})
	}
}
func TestUserManager_ResetPassword(t *testing.T) {
	var um = NewUserManager(newTestDB(t, &User{}))
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Success", args{username, email, "password123"}, false},
		{"Failure", args{"notarealuser", "notarealemail", "password123"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newPass, err := um.ResetPassword(tt.args.userName)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ResetPassword err = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && newPass == "" {
				t.Fatal("failed to reset password")
			}
		})
	}
}

func TestUserManager_Customer_Hash(t *testing.T) {
	var um = NewUserManager(newTestDB(t, &User{}))
	type newArgs struct {
		args
		firstHash  string
		secondHash string
	}
	tests := []struct {
		name    string
		args    newArgs
		wantErr bool
	}{
		{"Success", newArgs{args{username, email, "password123"}, "firsthash", EmptyCustomerObjectHash}, false},
		{"Failure", newArgs{args{"notarealusername", "notarealemail", "password123"}, "firsthash", "secondhash"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// test getting the default customer object hash
			if hash, err := um.GetCustomerObjectHash(tt.args.userName); (err != nil) != tt.wantErr {
				t.Fatalf("GetCustomerObjectHash err = %v, wantErr %v", err, tt.wantErr)
			} else if !tt.wantErr && hash != EmptyCustomerObjectHash {
				t.Fatal("failed to get correct customer object hash")
			}
			// test updating it
			if err := um.UpdateCustomerObjectHash(tt.args.userName, tt.args.firstHash); (err != nil) != tt.wantErr {
				t.Fatalf("UpdateCustomerObjectHash err = %v, wantErr %v", err, tt.wantErr)
			}
			// test getting the new updated hash
			if hash, err := um.GetCustomerObjectHash(tt.args.userName); (err != nil) != tt.wantErr {
				t.Fatalf("GetCustomerOBjectHash err = %v, wantErr %v", err, tt.wantErr)
			} else if !tt.wantErr && hash != tt.args.firstHash {
				t.Fatal("failed to get correct hash")
			}
			// reset the customer hash to the default one
			if err := um.UpdateCustomerObjectHash(tt.args.userName, tt.args.secondHash); (err != nil) != tt.wantErr {
				t.Fatalf("UpdateCustomerObjectHash err = %v, wantErr %v", err, tt.wantErr)
			}
			if hash, err := um.GetCustomerObjectHash(tt.args.userName); (err != nil) != tt.wantErr {
				t.Fatalf("GetCustomerOBjectHash err = %v, wantErr %v", err, tt.wantErr)
			} else if !tt.wantErr && hash != tt.args.secondHash {
				t.Fatal("failed to get correct hash")
			}
		})
	}
}

func TestUserManager_RemoveIPFSKeys(t *testing.T) {
	var um = NewUserManager(newTestDB(t, &User{}))
	type args struct {
		userName string
		email    string
		password string
		keyName  string
		keyID    string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Success", args{username, email, "password123", "keynamedelete", "keyidelete"}, false},
		{"Failure", args{"notarealuser", "notarealemail", "password123", "", ""}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := um.AddIPFSKeyForUser(
				tt.args.userName,
				tt.args.keyName,
				tt.args.keyID,
			); (err != nil) != tt.wantErr {
				t.Fatalf("AddIPFSKeyForUser err = %v, wantErr %v", err, tt.wantErr)
			}
			if err := um.RemoveIPFSKeyForUser(
				tt.args.userName,
				tt.args.keyName,
				tt.args.keyID,
			); (err != nil) != tt.wantErr {
				t.Fatalf("RemoveIPFSKeyForUser err = %v, wantErr %v", err, tt.wantErr)
			}
			// do not do any more processing for expected error tests
			if tt.wantErr {
				return
			}
			user, err := um.FindByUserName(tt.args.userName)
			if err != nil {
				t.Fatal(err)
			}
			for _, name := range user.IPFSKeyNames {
				if name == tt.args.keyName {
					t.Fatal("failed to corectly delete key name")
				}
			}
			for _, id := range user.IPFSKeyIDs {
				if id == tt.args.keyID {
					t.Fatal("failed to correctly delete key id")
				}
			}
		})
	}
}
