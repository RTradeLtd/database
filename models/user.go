package models

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// User contains a user account
type User struct {
	gorm.Model
	EthAddress        string         `gorm:"type:varchar(255);unique"`
	UserName          string         `gorm:"type:varchar(255);unique"`
	EmailAddress      string         `gorm:"type:varchar(255);unique"`
	EnterpriseEnabled bool           `gorm:"type:boolean"`
	AccountEnabled    bool           `gorm:"type:boolean"`
	APIAccess         bool           `gorm:"type:boolean"`
	EmailEnabled      bool           `gorm:"type:boolean"`
	HashedPassword    string         `gorm:"type:varchar(255)"`
	Credits           float64        `gorm:"type:float"`
	IPFSKeyNames      pq.StringArray `gorm:"type:text[];column:ipfs_key_names"`
	IPFSKeyIDs        pq.StringArray `gorm:"type:text[];column:ipfs_key_ids"`
	IPFSNetworkNames  pq.StringArray `gorm:"type:text[];column:ipfs_network_names"`
}

// UserManager is used to manage user accounts within the DB
type UserManager struct {
	DB *gorm.DB
}

// NewUserManager is a method used to generate our user manager
func NewUserManager(db *gorm.DB) *UserManager {
	um := UserManager{}
	um.DB = db
	return &um
}

// GetPrivateIPFSNetworksForUser is used to get all the authorized private networks a user has access to
func (um *UserManager) GetPrivateIPFSNetworksForUser(username string) ([]string, error) {
	u := &User{}
	if check := um.DB.Where("user_name = ?", username).First(u); check.Error != nil {
		return nil, check.Error
	}
	return u.IPFSNetworkNames, nil
}

// CheckIfUserHasAccessToNetwork is used to check if a user has access to a private IFPS network
func (um *UserManager) CheckIfUserHasAccessToNetwork(username, networkName string) (bool, error) {
	u := &User{}
	if check := um.DB.Where("user_name = ?", username).First(u); check.Error != nil {
		return false, check.Error
	}
	for _, v := range u.IPFSNetworkNames {
		if v == networkName {
			return true, nil
		}
	}
	return false, nil
}

// AddIPFSNetworkForUser is used to add an IPFS network to a users authorized networks
func (um *UserManager) AddIPFSNetworkForUser(username, networkName string) error {
	u := &User{}
	if check := um.DB.Where("user_name = ?", username).First(u); check.Error != nil {
		return check.Error
	}
	for _, v := range u.IPFSNetworkNames {
		if v == networkName {
			return errors.New("network already configured for user")
		}
	}
	u.IPFSNetworkNames = append(u.IPFSNetworkNames, networkName)
	if check := um.DB.Model(u).Update("ipfs_network_names", u.IPFSNetworkNames); check.Error != nil {
		return check.Error
	}

	return nil
}

// AddIPFSKeyForUser is used to add a key to a user account
func (um *UserManager) AddIPFSKeyForUser(username, keyName, keyID string) error {
	var user User
	if errCheck := um.DB.Where("user_name = ?", username).First(&user); errCheck.Error != nil {
		return errCheck.Error
	}

	if user.CreatedAt == nilTime {
		return errors.New("user account does not exist")
	}

	for _, v := range user.IPFSKeyNames {
		if v == keyName {
			fmt.Println("key already exists in db, skipping")
			return nil
		}
	}
	user.IPFSKeyNames = append(user.IPFSKeyNames, keyName)
	user.IPFSKeyIDs = append(user.IPFSKeyIDs, keyID)
	// The following only updates the specified column for the given model
	if errCheck := um.DB.Model(&user).Updates(map[string]interface{}{
		"ipfs_key_names": user.IPFSKeyNames,
		"ipfs_key_ids":   user.IPFSKeyIDs,
	}); errCheck.Error != nil {
		return errCheck.Error
	}
	return nil
}

// GetKeysForUser is used to fetch all of the users keys
func (um *UserManager) GetKeysForUser(username string) (map[string][]string, error) {
	var user User
	keys := make(map[string][]string)
	if errCheck := um.DB.Where("user_name = ?", username).First(&user); errCheck.Error != nil {
		return nil, errCheck.Error
	}

	if user.CreatedAt == nilTime {
		return nil, errors.New("user account does not exist")
	}

	keys["key_names"] = user.IPFSKeyNames
	keys["key_ids"] = user.IPFSKeyIDs
	return keys, nil
}

// GetKeyIDByName is used to get the ID (public key hash) of a key
func (um *UserManager) GetKeyIDByName(username, keyName string) (string, error) {
	var user User
	if errCheck := um.DB.Where("user_name = ?", username).First(&user); errCheck.Error != nil {
		return "", errCheck.Error
	}

	if user.CreatedAt == nilTime {
		return "", errors.New("user account does not exist")
	}
	for k, v := range user.IPFSKeyNames {
		if v == keyName {
			return user.IPFSKeyIDs[k], nil
		}
	}
	return "", errors.New("key not found")
}

// CheckIfKeyOwnedByUser is used to check if a given user owns a key
func (um *UserManager) CheckIfKeyOwnedByUser(username, keyName string) (bool, error) {
	var user User
	if errCheck := um.DB.Where("user_name = ?", username).First(&user); errCheck.Error != nil {
		return false, errCheck.Error
	}

	if user.CreatedAt == nilTime {
		return false, errors.New("user account does not exist")
	}

	for _, v := range user.IPFSKeyNames {
		if v == keyName {
			return true, nil
		}
	}
	return false, nil
}

// CheckIfUserAccountEnabled is used to get whether or not a user account is enabled
func (um *UserManager) CheckIfUserAccountEnabled(username string) (bool, error) {
	user := User{}
	if check := um.DB.Where("user_name = ?", username).First(&user); check.Error != nil {
		return false, check.Error
	}
	return user.AccountEnabled, nil
}

// ChangePassword is used to change a users password
func (um *UserManager) ChangePassword(username, currentPassword, newPassword string) (bool, error) {
	user := User{}
	if check := um.DB.Where("user_name = ?", username).First(&user); check.Error != nil {
		return false, check.Error
	}
	decodedPassword, err := hex.DecodeString(user.HashedPassword)
	if err != nil {
		return false, err
	}
	if err := bcrypt.CompareHashAndPassword(decodedPassword, []byte(currentPassword)); err != nil {
		return false, errors.New("invalid current password")
	}
	newHashedPass, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return false, err
	}
	encodedNewHashedPass := hex.EncodeToString(newHashedPass)
	if check := um.DB.Model(&user).Update("hashed_password", encodedNewHashedPass); check.Error != nil {
		return false, check.Error
	}
	return true, nil
}

// NewUserAccount is used to create a new user account, enabling API and account access
func (um *UserManager) NewUserAccount(ethAddress, username, password, email string, enterpriseEnabled bool) (*User, error) {
	user := User{}
	check := um.DB.Where("user_name = ?", username).First(&user)
	if check.Error != nil && check.Error != gorm.ErrRecordNotFound {
		return nil, check.Error
	}
	if user.CreatedAt != nilTime {
		return nil, errors.New("user account already created")
	}
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	if ethAddress == "" {
		ethAddress = username
	}
	user = User{
		UserName:          username,
		EthAddress:        ethAddress,
		EnterpriseEnabled: enterpriseEnabled,
		HashedPassword:    hex.EncodeToString(hashedPass),
		EmailAddress:      email,
		AccountEnabled:    true,
		APIAccess:         true,
	}
	if check := um.DB.Create(&user); check.Error != nil {
		return nil, check.Error
	}
	return &user, nil
}

// SignIn is used to authenticate a user, and check if their account is enabled.
// Returns bool on succesful login, or false with an error on failure
func (um *UserManager) SignIn(username, password string) (bool, error) {
	var user User
	if check := um.DB.Where("user_name = ?", username).First(&user); check.Error != nil {
		return false, check.Error
	}
	if !user.AccountEnabled {
		return false, errors.New("account is disabled")
	}
	validPassword, err := um.ComparePlaintextPasswordToHash(username, password)
	if err != nil {
		return false, err
	}
	if !validPassword {
		return false, errors.New("invalid password supplied")
	}
	return true, nil
}

// ComparePlaintextPasswordToHash is a helper method used to validate the provided password
func (um *UserManager) ComparePlaintextPasswordToHash(username, password string) (bool, error) {
	user := User{}
	if check := um.DB.Where("user_name = ?", username).First(&user); check.Error != nil {
		return false, check.Error
	}
	passwordBytes, err := hex.DecodeString(user.HashedPassword)
	if err != nil {
		return false, err
	}
	err = bcrypt.CompareHashAndPassword(passwordBytes, []byte(password))
	if err != nil {
		return false, err
	}
	return true, nil

}

// FindByAddress is used to find a user by searching for their eth address
func (um *UserManager) FindByAddress(address string) (*User, error) {
	u := User{}
	if check := um.DB.Where("eth_address = ?", address).First(&u); check.Error != nil {
		return nil, check.Error
	}
	return &u, nil
}

// FindEthAddressByUserName is used to find the eth address associated with a user account
func (um *UserManager) FindEthAddressByUserName(username string) (string, error) {
	u := User{}
	if check := um.DB.Where("user_name = ?", username).First(&u); check.Error != nil {
		return "", check.Error
	}
	return u.EthAddress, nil
}

// FindEmailByUserName is used to find an email address by searching for the users eth address
// the returned map contains their eth address as a key, and their email address as a value
func (um *UserManager) FindEmailByUserName(username string) (map[string]string, error) {
	u := User{}
	if check := um.DB.Where("user_name = ?", username).First(&u); check.Error != nil {
		return nil, check.Error
	}
	emails := make(map[string]string)
	emails[username] = u.EmailAddress
	return emails, nil
}

// FindByUserName is used to find a user by their username
func (um *UserManager) FindByUserName(username string) (*User, error) {
	u := User{}
	if check := um.DB.Where("user_name = ?", username).First(&u); check.Error != nil {
		return nil, check.Error
	}
	return &u, nil
}

// ChangeEthereumAddress is used to change a user's ethereum address
func (um *UserManager) ChangeEthereumAddress(username, ethAddress string) (*User, error) {
	u := User{}
	if check := um.DB.Where("user_name = ?", username).First(&u); check.Error != nil {
		return nil, check.Error
	}
	u.EthAddress = ethAddress
	if check := um.DB.Model(u).Update("eth_address", ethAddress); check.Error != nil {
		return nil, check.Error
	}
	return &u, nil
}

// AddCreditsForUser is used to add credits to a users account
func (um *UserManager) AddCreditsForUser(username string, credits float64) (*User, error) {
	u := User{}
	if check := um.DB.Where("user_name = ?", username).First(&u); check.Error != nil {
		return nil, check.Error
	}
	if u.Credits > 0 {
		u.Credits = u.Credits + credits
	} else {
		u.Credits = credits
	}
	if check := um.DB.Model(u).Update("credits", credits); check.Error != nil {
		return nil, check.Error
	}
	return &u, nil
}

// GetCreditsForUser is used to get the user's current credits
func (um *UserManager) GetCreditsForUser(username string) (float64, error) {
	u := User{}
	if check := um.DB.Where("user_name = ?", username).First(&u); check.Error != nil {
		return 0, check.Error
	}
	return u.Credits, nil
}

// RemoveCredits is used to remove credits from a users balance
func (um *UserManager) RemoveCredits(username string, credits float64) (*User, error) {
	user, err := um.FindByUserName(username)
	if err != nil {
		return nil, err
	}
	if user.Credits < credits {
		return nil, errors.New("unable to remove credits, would result in negative balance")
	}
	user.Credits = user.Credits - credits
	if check := um.DB.Model(user).Update("credits", user.Credits); check.Error != nil {
		return nil, check.Error
	}
	return user, nil
}
