package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"ai-developer-workbench/internal/config"
	"ai-developer-workbench/internal/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect establishes a database connection based on the configured driver.
func Connect(cfg *config.DatabaseConfig, isDev bool) (*gorm.DB, error) {
	var logLevel logger.LogLevel
	if isDev {
		logLevel = logger.Info
	} else {
		logLevel = logger.Warn
	}

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	}

	dsn := cfg.DSN()
	slog.Info("Connecting to MySQL", "host", cfg.Host, "port", cfg.Port, "database", cfg.Name)
	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Configure connection pool.
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Minute)

	return db, nil
}

// Ping checks if the database connection is alive.
func Ping(ctx context.Context, db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}

// Close closes the database connection.
func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// RunMigrations executes database migrations.
func RunMigrations(db *gorm.DB, autoMigrate bool) error {
	if autoMigrate {
		if err := db.AutoMigrate(
			&model.Report{},
			&model.GeneratedFile{},
			&model.ReportAsset{},
		); err != nil {
			return fmt.Errorf("auto-migrate failed: %w", err)
		}
	}
	return nil
}