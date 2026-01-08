package database

import (
	"fmt"
	"log"
	"time"

	"github.com/Rugved7/collaborative-code-review/internal/common/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitProgress(cfg *config.Config) (*gorm.DB, error) {
	var err error
	var logLevel logger.LogLevel

	if cfg.Environment == "production" {
		logLevel = logger.Error
	} else {
		logLevel = logger.Info
	}

	gormConfig := &gorm.Config{
		Logger: logger.New(log.New(nil, "", log.LstdFlags), logger.Config{LogLevel: logLevel}),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		PrepareStmt: true, // prepared statement cache for performance
	}

	// Connect to postgres database
	DB, err := gorm.Open(postgres.Open(cfg.DatabaseURL), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get Database instance: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.DBMaxConnections)
	sqlDB.SetMaxIdleConns(cfg.DBMaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.DBConnMaxLifeTime)
	sqlDB.SetConnMaxIdleTime(30 * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	log.Printf("[Database] PostgreSQL connected successfully (Pool: max=%d, idle=%d)",
		cfg.DBMaxConnections, cfg.DBMaxIdleConns)

	return DB, nil
}

// Close the DB connection
func Close() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// HealthCheck verifies database connectivity
func HealthCheck() error {
	if DB == nil {
		return fmt.Errorf("database not initialized")
	}
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}
