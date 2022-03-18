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
	UserByName(string) (*models.User, error)
	AddUser(*models.User) error
	GetAllUsers() ([]models.User, error)
	AutoMigrate() error
	Init() error
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
	if err := db.connection.AutoMigrate(&models.User{}); err != nil {
		return err
	}
	return db.connection.AutoMigrate(&models.User{})
}

func (db *database) Init() error {
	if err := db.connect(); err != nil {
		return err
	}
	defer db.close()
	// ToDo: move this to env var!
	// Create admin user
	result := db.connection.Create(&models.User{Username: "floge77", Password: "secure", Role: models.Admin})
	return result.Error
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

func (db *database) UserByName(username string) (*models.User, error) {
	if err := db.connect(); err != nil {
		return nil, err
	}
	defer db.close()
	var user *models.User
	result := db.connection.First(&user, "username = ?", username)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (db *database) AddUser(user *models.User) error {
	if err := db.connect(); err != nil {
		return err
	}
	defer db.close()
	result := db.connection.Create(user)
	return result.Error
}

func (db *database) GetAllUsers() ([]models.User, error) {
	var users []models.User
	if err := db.connect(); err != nil {
		return nil, err
	}
	defer db.close()
	if result := db.connection.Find(&users); result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}
