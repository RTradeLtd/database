package models

import (
	"errors"

	"github.com/RTradeLtd/gorm"
	"github.com/lib/pq"
)

// Zone is a TNS zone
type Zone struct {
	gorm.Model
	Name                 string         `gorm:"type:varchar(255)"`
	ManagerPublicKeyName string         `gorm:"type:varchar(255)"`
	ZonePublicKeyName    string         `gorm:"type:varchar(255)"`
	LatestIPFSHash       string         `gorm:"type:varchar(255)"`
	RecordNames          pq.StringArray `gorm:"type:text[]"`
}

// ZoneManager is used to manipulate zone entries in the database
type ZoneManager struct {
	DB *gorm.DB
}

// NewZoneManager is used to generate our zone manager helper to interact with the db
func NewZoneManager(db *gorm.DB) *ZoneManager {
	return &ZoneManager{DB: db}
}

// NewZone is used to create a new zone in the database
func (zm *ZoneManager) NewZone(name, managerPK, zonePK, latestIPFSHash string) (*Zone, error) {
	zone, err := zm.FindZoneByName(name)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == nil {
		return nil, errors.New("zone already exists for user")
	}
	zone = &Zone{
		Name:                 name,
		ManagerPublicKeyName: managerPK,
		ZonePublicKeyName:    zonePK,
		LatestIPFSHash:       latestIPFSHash,
	}
	if check := zm.DB.Create(zone); check.Error != nil {
		return nil, check.Error
	}
	return zone, nil
}

// FindZoneByName is used to lookup a zone by name
func (zm *ZoneManager) FindZoneByName(name string) (*Zone, error) {
	z := Zone{}
	if check := zm.DB.Where("name = ?", name).First(&z); check.Error != nil {
		return nil, check.Error
	}
	return &z, nil
}

// UpdateLatestIPFSHashForZone is used to update the latest IPFS hash for a zone file
func (zm *ZoneManager) UpdateLatestIPFSHashForZone(name, hash string) (*Zone, error) {
	z, err := zm.FindZoneByName(name)
	if err != nil {
		return nil, err
	}
	z.LatestIPFSHash = hash
	if check := zm.DB.Model(&z).Update("latest_ip_fs_hash", z.LatestIPFSHash); check.Error != nil {
		return nil, check.Error
	}
	return z, nil
}

// AddRecordForZone is used to add a record to a zone
func (zm *ZoneManager) AddRecordForZone(zoneName, recordName string) (*Zone, error) {
	z, err := zm.FindZoneByName(zoneName)
	if err != nil {
		return nil, err
	}
	present, err := zm.CheckIfRecordExistsInZone(zoneName, recordName)
	if err != nil {
		return nil, err
	}
	if present {
		return nil, errors.New("record already exists in zone")
	}
	z.RecordNames = append(z.RecordNames, recordName)
	if check := zm.DB.Model(z).Update("record_names", z.RecordNames); check.Error != nil {
		return nil, check.Error
	}
	return z, nil
}

// CheckIfRecordExistsInZone is used to check if a record exists in a particular zone
func (zm *ZoneManager) CheckIfRecordExistsInZone(zoneName, recordName string) (bool, error) {
	z, err := zm.FindZoneByName(zoneName)
	if err != nil {
		return false, err
	}
	if len(z.RecordNames) == 0 {
		return false, nil
	}
	for _, v := range z.RecordNames {
		if v == recordName {
			return true, nil
		}
	}
	return false, nil
}
