package db

import (
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" // sqlite dialect
	"github.com/txross1993/superman-api/models"
)

const dbfile = "./gorm.db"

// DB is the concrete implementation of persistence
type DB struct {
	db       *gorm.DB
	filePath string
}

// InitDB creates a new DB instance provided a local db file path
// and reflects the API models to the backend storage layer
func InitDB(dbFile string) (DB, error) {
	var repo DB
	db, err := gorm.Open("sqlite3", dbFile)
	if err != nil {
		return repo, err
	}

	repo.db = db
	repo.filePath = dbFile

	if err := repo.db.AutoMigrate(&models.UserIPAccessEvent{}).Error; err != nil {
		return repo, err
	}

	return repo, nil
}

// Cleanup deletes the db storage file, deleting all data
func (d DB) Cleanup() error {
	return os.Remove(d.filePath)
}

// Close ends the database connection
func (d DB) Close() error {
	return d.db.Close()
}

// FindOrCreateUserIPAccessEvent will save the ip access event record if new
func (d DB) FindOrCreateUserIPAccessEvent(event *models.UserIPAccessEvent) error {
	return d.db.FirstOrCreate(&event).Error
}

// FindPrecedingIPAccessEvent retrieves the ip access event that occurred most
// recently before the input event if any
func (d DB) FindPrecedingIPAccessEvent(event *models.UserIPAccessEvent) (*models.UserIPAccessEvent, error) {
	var priorEvent models.UserIPAccessEvent
	err := d.db.Limit(1).Where("username = ?", event.Username).Where("unix_timestamp  <= ?", event.UnixTimestamp).Where("event_uuid != ?", event.EventUUID).Order("unix_timestamp DESC").Find(&priorEvent).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &priorEvent, nil

}

// FindSubsequentIPAccessEvent retrieves the ip access event that occurred most
// recently after the input event if any
func (d DB) FindSubsequentIPAccessEvent(event *models.UserIPAccessEvent) (*models.UserIPAccessEvent, error) {
	var subsequentEvent models.UserIPAccessEvent
	err := d.db.Limit(1).Where("username = ?", event.Username).Where("unix_timestamp  >= ?", event.UnixTimestamp).Where("event_uuid != ?", event.EventUUID).Order("unix_timestamp ASC").Find(&subsequentEvent).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &subsequentEvent, nil
}
