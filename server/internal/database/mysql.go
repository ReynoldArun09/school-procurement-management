package database

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/ReynoldArun09/school-procurement-management/internal/config"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func NewMySQL(ctx context.Context, cfg config.MysqlConfig, logger *slog.Logger) (*gorm.DB, error) {
	mysqlConfig := mysqlDriver.Config{
		User:      cfg.User,
		Passwd:    cfg.Password,
		Net:       "tcp",
		Addr:      net.JoinHostPort(cfg.Host, cfg.Port),
		DBName:    cfg.Database,
		ParseTime: true,
		Loc:       time.UTC,
	}
	dsn := mysqlConfig.FormatDSN()

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableAutomaticPing: true,
		Logger: gormlogger.Default.LogMode(
			gormlogger.Error,
		),
	})

	if err != nil {
		return nil, fmt.Errorf("initialize mysql connection: %w", err)
	}

	sqlDB, err := db.DB()

	if err != nil {
		return nil, fmt.Errorf("get underlying sql connection pool: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(pingCtx); err != nil {
		if closeErr := sqlDB.Close(); closeErr != nil {
			logger.Error("failed to close Mysql pool after connection failure", "error", closeErr)
		}
		return nil, fmt.Errorf("connect to Mysql: %w", err)
	}
	logger.Info("Mysql connection established", "host", cfg.Host, "database", cfg.Database)
	return db, nil

}
