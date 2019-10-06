package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // mysql driver
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"main.go/api"
	"main.go/models"
)

// Database struct
type Database struct {
	*gorm.DB
}

func migrateTable(db *Database) {
	logrus.Infof("Migrating tables...")
	db.AutoMigrate(
		&models.Role{},
		&models.User{},
		&models.Schedule{},
		&models.Attendance{},
	)
	logrus.Infof("All tables migrated.")
}

// New connection
func New() (*Database, error) {
	var DBUri = viper.GetString("database.user") + ":" + viper.GetString("database.password") + "@/" + viper.GetString("database.database_name") + "?charset=utf8&parseTime=True&loc=Local"
	if DBUri == "" {
		logrus.Warnf("DatabaseURI must be set")
	}

	db, err := gorm.Open("mysql", DBUri)
	api.DB = db
	
	if err != nil {
		return nil, errors.Wrap(err, "Unable to connect to database")
	}
	
	migrateTable(&Database{db})

	return &Database{db}, nil
}
