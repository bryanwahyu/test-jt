package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"github.com/bryanwahyu/test-jt/internal/config"
	"github.com/bryanwahyu/test-jt/internal/phone"
)

var db *gorm.DB

// InitDB initializes the database connection
func InitDB(cfg config.Config) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBName)

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	// Access the underlying *sql.DB instance
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)

	// AutoMigrate ensures the table is created
	if err := db.AutoMigrate(&phone.Phone{}); err != nil {
		return err
	}

	return nil
}

func Migrate() error {
    // AutoMigrate ensures the necessary tables are created
    return db.AutoMigrate(&phone.Phone{})
}
// GetDB returns the database instance
func GetDB() *gorm.DB {
	return db
}
