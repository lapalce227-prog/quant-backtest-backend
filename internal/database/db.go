package database

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	Charset         string
	ParseTime       bool
	Loc             string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
}

func OpenMySQLFromEnv() (*gorm.DB, error) {
	cfg := Config{
		Host:            getEnv("MYSQL_HOST", "127.0.0.1"),
		Port:            getEnv("MYSQL_PORT", "3306"),
		User:            getEnv("MYSQL_USER", "root"),
		Password:        getEnv("MYSQL_PASSWORD", "root"),
		Name:            getEnv("MYSQL_DB", "tradingsystem"),
		Charset:         getEnv("MYSQL_CHARSET", "utf8mb4"),
		ParseTime:       getEnvBool("MYSQL_PARSE_TIME", true),
		Loc:             getEnv("MYSQL_LOC", "Local"),
		MaxIdleConns:    getEnvInt("MYSQL_MAX_IDLE_CONNS", 10),
		MaxOpenConns:    getEnvInt("MYSQL_MAX_OPEN_CONNS", 20),
		ConnMaxLifetime: time.Duration(getEnvInt("MYSQL_CONN_MAX_LIFETIME_MINUTES", 30)) * time.Minute,
	}

	return OpenMySQL(cfg)
}

func OpenMySQL(cfg Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%t&loc=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.Charset,
		cfg.ParseTime,
		cfg.Loc,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	return db, nil
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func getEnvInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}

func getEnvBool(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}

	return parsed
}
