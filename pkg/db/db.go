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

type Database interface {
	GetAllCredentials() ([]models.Credentials, error)
	AddCredential(*models.Credentials) error
	AutoMigrate() error
}

type database struct {
	connection *gorm.DB
	address    string
	port       string
	user       string
	password   string
	dbName     string
}

func NewDatabase(address, port, user, password, dbName string) Database {
	return &database{
		address:  address,
		port:     port,
		user:     user,
		password: password,
		dbName:   dbName,
	}
}

func (db *database) AutoMigrate() error {
	if err := db.connect(); err != nil {
		return err
	}
	defer db.close()
	return db.connection.AutoMigrate(&models.Credentials{}) // add more models if needed
}

func (db *database) connect() error {
	dbURL := fmt.Sprintf("%s://%s:%s@%s:%s/%s", schema, db.user, db.password, db.address, db.port, db.dbName)
	conn, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		return err
	}
	database, err := conn.DB()
	if err != nil {
		return err
	}
	if err = database.Ping(); err != nil {
		return err
	}
	db.connection = conn
	return nil
}

func (db *database) close() error {
	sqlDb, err := db.connection.DB()
	if err != nil {
		return err
	}
	return sqlDb.Close()
}

func (db *database) GetAllCredentials() ([]models.Credentials, error) {
	var credentials []models.Credentials
	if err := db.connect(); err != nil {
		return nil, err
	}
	defer db.close()
	if result := db.connection.Find(&credentials); result.Error != nil {
		return nil, result.Error
	}
	return credentials, nil
}

func (db *database) AddCredential(credential *models.Credentials) error {
	if err := db.connect(); err != nil {
		return err
	}
	defer db.close()
	result := db.connection.Create(credential)
	return result.Error
}
