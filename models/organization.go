package models

import (
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
	um := NewUserManager(om.DB)
	user, err := um.NewUserAccount(username, password, email)
	if err != nil {
		return nil, err
	}
	// TODO(postables): update user model as being part of organization
	org, err := om.FindByName(name)
	if err != nil {
		return nil, err
	}
	// update registered users
	org.RegisteredUsers = append(org.RegisteredUsers, username)
	// update org model model
	if err := om.DB.Model(org).Update(
		"regisered_users",
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
