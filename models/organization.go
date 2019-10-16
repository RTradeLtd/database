package models

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

// Organization represents a company using white-labeled Temporal
type Organization struct {
	gorm.Model
	// the name of the organization
	Name string `gorm:"type:varchar(255);unique"`
	// the corresponding temporal user account that manages this org
	UserOwner string `gorm:"type:varchar(255);unique"`
	// the usd value owed by the organization
	AccountBalance float64 `gorm:"type:float"`
	// the user accounts who have signed up under this organization
	RegisteredUsers pq.StringArray `gorm:"type:text[];column:registered_users"`
}

// OrgManager is an organization model manager
type OrgManager struct {
	DB *gorm.DB
}

// NewOrgManager instantiates an OrgManager
func NewOrgManager(db *gorm.DB) *OrgManager {
	return &OrgManager{DB: db}
}

// NewOrganization is used to create a new organization
func (om *OrgManager) NewOrganization(name, owner string) error {
	org := &Organization{
		Name:      name,
		UserOwner: owner,
	}
	return om.DB.Create(org).Error
}

// RegisterOrgUser registers an organization user
func (om *OrgManager) RegisterOrgUser(
	name,
	username,
	password,
	email string,
) (*User, error) {
	// create the user account
	user, err := NewUserManager(om.DB).NewUserAccount(
		username,
		password,
		email,
	)
	if err != nil {
		return nil, err
	}
	// update user model associated organization
	user.Organization = name
	// save updated user model
	if err := om.DB.Model(user).Update(
		"organization", user.Organization,
	).Error; err != nil {
		return nil, err
	}
	// update their tier to white-labeled
	// which will enable organizational based billing
	if err := NewUsageManager(om.DB).UpdateTier(
		username,
		WhiteLabeled,
	); err != nil {
		return nil, err
	}
	// find organization model
	org, err := om.FindByName(name)
	if err != nil {
		return nil, err
	}
	// update organization registered users
	org.RegisteredUsers = append(org.RegisteredUsers, username)
	// save updated org model model
	if err := om.DB.Model(org).Update(
		"registered_users",
		org.RegisteredUsers,
	).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// FindByName finds an organization by name
func (om *OrgManager) FindByName(name string) (*Organization, error) {
	org := &Organization{}
	if err := om.DB.Where(
		"name = ?",
		name,
	).First(&org).Error; err != nil {
		return nil, err
	}
	return org, nil
}

// GetOrgUsers is used toreturn the registered users an organization has
func (om *OrgManager) GetOrgUsers(name string) ([]string, error) {
	org, err := om.FindByName(name)
	if err != nil {
		return nil, err
	}
	return org.RegisteredUsers, nil
}

// BillingReport contains a summary
// of an organizations entire active
// user base in the last 30 days along with
// the USD value currently owned by the account
type BillingReport struct {
	Name  string        `json:"name"`
	Items []BillingItem `json:"items"`
	// amount owed in USD
	AmountDue float64 `json:"amount_due"`
}

// BillingItem is an individual user's
// billing history
type BillingItem struct {
	User    string   `json:"user"`
	Uploads []Upload `json:"uploads"`
}

// GenerateBillingReport is used to generate a billing report object for an
// organization's entire user base in the last 30 days. Care must be taken so that
// only the organization owner may interact with this function, and is it returns sensitive information
func (om *OrgManager) GenerateBillingReport(name string) (*BillingReport, error) {
	org, err := om.FindByName(name)
	if err != nil {
		return nil, err
	}
	report := &BillingReport{Name: name, AmountDue: org.AccountBalance}
	for _, usr := range org.RegisteredUsers {
		// sanity check that the user exists
		if _, err := NewUserManager(
			om.DB,
		).FindByUserName(usr); err != nil {
			// dont fail and return, just continue onto the next user
			continue
		}
		var uploads []Upload
		// find all uploads from the user that were
		// updated in the last 30 days, we dont check
		// create at since it is possible for uploads to be extended
		if err := om.DB.Model(Upload{}).Where(
			"user_name = ? AND updated_at BETWEEN ? AND ?",
			usr, time.Now().AddDate(0, 0, -30), time.Now(),
		).Find(&uploads).Error; err != nil {
			// dont fail and return, just continue onto the next user
			continue
		}
		if len(uploads) == 0 {
			continue
		}
		report.Items = append(report.Items, BillingItem{
			User:    usr,
			Uploads: uploads,
		})
	}
	return report, nil
}
