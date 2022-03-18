package db

import (
	"fmt"

	"github.com/floge77/go-authenticator/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	schema = "postgres"
)

var (
	connectToDb = connect
)

type DatabaseConnection interface {
	GetAllCredentials() ([]models.Credentials, error)
	AddCredential(*models.Credentials) error
}

type databaseConnection struct {
	db *gorm.DB
}

func NewDatabaseConnection(address, port, user, password, dbName string) (DatabaseConnection, error) {
	db, err := connectToDb(address, port, user, password, dbName)
	if err != nil {
		return nil, err
	}
	return &databaseConnection{db: db}, nil
}

func connect(address, port, user, password, dbName string) (*gorm.DB, error) {
	dbURL := fmt.Sprintf("%s://%s:%s@%s:%s/%s", schema, user, password, address, port, dbName)
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})

	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&models.Credentials{})
	return db, nil
}

func (d *databaseConnection) GetAllCredentials() ([]models.Credentials, error) {
	var credentials []models.Credentials
	if result := d.db.Find(&credentials); result.Error != nil {
		return nil, result.Error
	}
	return credentials, nil
}

func (d *databaseConnection) AddCredential(credential *models.Credentials) error {
	result := d.db.Create(credential)
	return result.Error
}
