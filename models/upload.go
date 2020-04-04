package models

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/RTradeLtd/database/v2/utils"
	"github.com/c2h5oh/datasize"
	"github.com/jinzhu/gorm"
)

const (
	// ErrShorterGCD is an error triggered when updating to update an upload for a user
	// with a hold time that would result in a shorter garbage collection date
	ErrShorterGCD = "upload would not extend garbage collection date so there is no need to process"
	// ErrAlreadyExistingUpload is an error triggered when attempting to insert  a new row into the database
	// for a content that already exists in the database for a user. This means you should be using the UpdateUpload
	// function to allow for updating garbage collection dates.
	ErrAlreadyExistingUpload = "the content you are inserting into the database already exists, please use the UpdateUpload function"
)

// Upload is a file or pin based upload to temporal
type Upload struct {
	gorm.Model
	Hash               string `gorm:"type:varchar(255);not null;"`
	Type               string `gorm:"type:varchar(255);not null;"` //  file, pin
	NetworkName        string `gorm:"type:varchar(255)"`
	HoldTimeInMonths   int64  `gorm:"type:integer;not null;"`
	UserName           string `gorm:"type:varchar(255);not null;"`
	GarbageCollectDate time.Time
	Encrypted          bool   `gorm:"type:bool"`
	FileName           string `gorm:"type:varchar(255)"`
	FileNameLowerCase  string `gorm:"type:varchar(255)"`
	FileNameUpperCase  string `gorm:"type:varchar(255)"`
	Extension          string `gorm:"type:varchar(255)"`
	Size               int64  `gorm:"type:bigint"` // upload size in bytes
	Directory          bool   `gorm:"type:bool;default:false"`
}

// UploadManager is used to manipulate upload objects in the database
type UploadManager struct {
	DB *gorm.DB
}

// NewUploadManager is used to generate an upload manager interface
func NewUploadManager(db *gorm.DB) *UploadManager {
	return &UploadManager{DB: db}
}

// UploadOptions is used to configure an upload
type UploadOptions struct {
	NetworkName      string
	Username         string
	FileName         string
	HoldTimeInMonths int64
	Size             int64
	Encrypted        bool
	Directory        bool
}

// NewUpload is used to create a new upload in the database
func (um *UploadManager) NewUpload(contentHash, uploadType string, opts UploadOptions) (*Upload, error) {
	_, err := um.FindUploadByHashAndUserAndNetwork(opts.Username, contentHash, opts.NetworkName)
	if err == nil {
		// this means that there is already an upload in hte database matching this content hash and network name, so we will skip
		return nil, errors.New(ErrAlreadyExistingUpload)
	}
	holdInt, err := strconv.Atoi(fmt.Sprintf("%+v", opts.HoldTimeInMonths))
	if err != nil {
		return nil, err
	}
	upload := Upload{
		Hash:               contentHash,
		Type:               uploadType,
		NetworkName:        opts.NetworkName,
		HoldTimeInMonths:   opts.HoldTimeInMonths,
		UserName:           opts.Username,
		GarbageCollectDate: utils.CalculateGarbageCollectDate(holdInt),
		Encrypted:          opts.Encrypted,
		FileName:           opts.FileName,
		FileNameLowerCase:  strings.ToLower(opts.FileName),
		FileNameUpperCase:  strings.ToUpper(opts.FileName),
		Extension:          filepath.Ext(opts.FileName),
		Size:               opts.Size,
		Directory:          opts.Directory,
	}
	if check := um.DB.Create(&upload); check.Error != nil {
		return nil, check.Error
	}
	return &upload, nil
}

// UpdateUpload is used to update the garbage collection time for an already existing upload
func (um *UploadManager) UpdateUpload(holdTimeInMonths int64, username, contentHash, networkName string) (*Upload, error) {
	upload, err := um.FindUploadByHashAndUserAndNetwork(username, contentHash, networkName)
	if err != nil {
		return nil, err
	}
	oldGcd := upload.GarbageCollectDate
	newGcd := utils.CalculateGarbageCollectDate(int(holdTimeInMonths))
	if newGcd.Unix() < oldGcd.Unix() {
		return nil, errors.New(ErrShorterGCD)
	}
	upload.HoldTimeInMonths = holdTimeInMonths
	upload.GarbageCollectDate = newGcd
	if check := um.DB.Save(upload); check.Error != nil {
		return nil, err
	}
	return upload, nil
}

// FindUploadsByNetwork is used to find all uploads corresponding to a given network
func (um *UploadManager) FindUploadsByNetwork(networkName string) ([]Upload, error) {
	uploads := []Upload{}
	if check := um.DB.Where("network_name = ?", networkName).Find(&uploads); check.Error != nil {
		return nil, check.Error
	}
	return uploads, nil
}

// FindUploadByHashAndNetwork is used to search for an upload by its hash, and the network it was stored on
func (um *UploadManager) FindUploadByHashAndNetwork(hash, networkName string) (*Upload, error) {
	upload := &Upload{}
	if check := um.DB.Where("hash = ? AND network_name = ?", hash, networkName).First(upload); check.Error != nil {
		return nil, check.Error
	}
	return upload, nil
}

// FindUploadsByHash is used to return all instances of uploads matching the given hash
func (um *UploadManager) FindUploadsByHash(hash string) ([]Upload, error) {
	uploads := []Upload{}
	if err := um.DB.Where("hash = ?", hash).Find(&uploads).Error; err != nil {
		return nil, err
	}
	return uploads, nil
}

// FindUploadByHashAndUserAndNetwork is used to look for an upload based off its hash, user, and network
func (um *UploadManager) FindUploadByHashAndUserAndNetwork(username, hash, networkName string) (*Upload, error) {
	upload := &Upload{}
	if err := um.DB.Where("user_name = ? AND hash = ? AND network_name = ?", username, hash, networkName).First(upload).Error; err != nil {
		return nil, err
	}
	return upload, nil
}

// GetUploadByHashForUser is used to retrieve the last (most recent) upload for a user
func (um *UploadManager) GetUploadByHashForUser(hash string, username string) ([]Upload, error) {
	uploads := []Upload{}
	if err := um.DB.Where("hash = ? AND user_name = ?", hash, username).Find(&uploads).Error; err != nil {
		return nil, err
	}
	return uploads, nil
}

// GetUploads is used to return all  uploads
func (um *UploadManager) GetUploads() ([]Upload, error) {
	uploads := []Upload{}
	if check := um.DB.Find(&uploads); check.Error != nil {
		return nil, check.Error
	}
	return uploads, nil
}

// GetUploadsForUser is used to retrieve all uploads by a user name
func (um *UploadManager) GetUploadsForUser(username string) ([]Upload, error) {
	uploads := []Upload{}
	if check := um.DB.Where("user_name = ?", username).Find(&uploads); check.Error != nil {
		return nil, check.Error
	}
	return uploads, nil
}

// ExtendGarbageCollectionPeriod is used to extend the garbage collection period for a particular upload
func (um *UploadManager) ExtendGarbageCollectionPeriod(username, hash, network string, holdTimeInMonths int) error {
	upload, err := um.FindUploadByHashAndUserAndNetwork(username, hash, network)
	if err != nil {
		return err
	}
	// update garbage collection period
	upload.GarbageCollectDate = upload.GarbageCollectDate.AddDate(0, holdTimeInMonths, 0)
	// save the updated model
	return um.DB.Model(upload).Update("garbage_collect_date", upload.GarbageCollectDate).Error
}

// PinRM allows removing a pin and refunding extra data costs
func (um *UploadManager) PinRM(username, hash, network string) error {
	upload, err := um.FindUploadByHashAndUserAndNetwork(username, hash, network)
	if err != nil {
		return err
	}
	// get the amount to refund the user before removing the upload
	refundAmt, err := um.CalculateRefundCost(upload)
	if err != nil {
		return err
	}
	// remove upload returning if this fails
	if err := um.DB.Delete(upload).Error; err != nil {
		return err
	}
	// will be greater than 0 if they are not free
	// as only non-free users will need to have their credits refunded
	if refundAmt > 0 {
		// add credits to the user's balance
		_, err = NewUserManager(um.DB).AddCredits(username, refundAmt)
		if err != nil {
			return err
		}
	}
	// reduce user's storage consumption
	return NewUsageManager(um.DB).ReduceDataUsage(username, uint64(upload.Size))
}

// CalculateRefundCost returns the amount of credits to refund the user
// when they invoke pinRM
func (um *UploadManager) CalculateRefundCost(upload *Upload) (float64, error) {
	// do this first to not waste time processing needlessly
	usg, err := NewUsageManager(um.DB).FindByUserName(upload.UserName)
	if err != nil {
		return 0, err
	}
	if usg.Tier == Free {
		return 0, nil
	}
	var (
		startDate  time.Time
		removeDate = upload.GarbageCollectDate
	)
	if upload.UpdatedAt == nilTime {
		startDate = upload.CreatedAt
	} else {
		startDate = upload.UpdatedAt
	}
	now := time.Now()
	// indicates the number of days we have stored this object for
	daysStored := now.Sub(startDate)
	// get the number of hours remaining so we can calculate a refund
	// shave of 72 hour buffer from refund amount
	daysRemaining := removeDate.AddDate(0, 0, -3).Sub(now).Truncate(time.Hour)
	// total number of hours to refund minus an additional 72 hour buffer
	// helps to ensure that on all edge cases we dont refund the user extra
	// but they will be refunded slightly less, however this is deemed
	// acceptable, and is intentional for a few reasons:
	// 	* Refunds aren't a required feature of the platform as specified in Terms And Service
	//  * The pin wont actually be removed from the underlying IPFS node for up to 4 weeks
	//  * Prevent exploits to gain perpetual free credits due to rounding errors
	//  * Time spent processing the data as creating many pins and then removing them isn't cheap
	//	* Unpinning data requires bringing nodes offline and needs to be scheduled, thus this increases maintenance duties
	// Whenever removing a pin, there is a 72 hour buffer, which means even
	// if you pin data, and remove it immediately, you will still be charged 72 hours worth of data storage
	// this helps mitigate abuse of the system by having to have our nodes be under sustained GC load as removing
	// data from the system isn't a cheap process due to extreme inefficiencies with go-ipfs
	var refundHours float64
	// if less than or equal to 72 hours, don't refund anything
	if (daysRemaining - daysStored).Hours() <= (time.Hour.Hours() * 72) {
		refundHours = 0
	} else {
		refundHours = (daysRemaining - daysStored).Hours() - (time.Hour.Hours() * 72)
	}
	// calculates a refund based on the size of the object
	return calculateSizeRefund(refundHours, upload.Size, usg)
}

func calculateSizeRefund(refundHours float64, size int64, usage *Usage) (float64, error) {
	gigabytesFloat := float64(datasize.GB.Bytes())
	sizeFloat := float64(size)
	// returns how many gigabytes this upload is
	sizeGigabytesFloat := sizeFloat / gigabytesFloat
	// if they are free tier, they don't incur data charges
	if usage.Tier == Free || usage.Tier == WhiteLabeled {
		return 0, nil
	}
	// return the cost of the refund calculated by:
	// * number of hours remaining * gigabytes per hour = size cost multipliier
	// * size of data multiplied by size cost multiplier
	return sizeGigabytesFloat * (usage.Tier.PricePerGBPerHour() * refundHours), nil
}
