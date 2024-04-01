package initializers

import (
	"fmt"
	"github.com/Conty111/SuperCalculator/back-end/models"
	"github.com/Conty111/SuperCalculator/back-end/orkestrator/internal/config"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitMigrations(db *gorm.DB) {
	var (
		tasks models.TasksModel
	)
	err := db.AutoMigrate(&tasks)
	if err != nil {
		log.Panic().Err(err).Msg("Cannot run auto migrations")
	}
}

func InitializeDatabase(dbDSN string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(dbDSN), &gorm.Config{})
	if err != nil {
		log.Panic().Err(err).Str("path", dbDSN).Msg("Cannot connect to the database")
	}

	InitMigrations(db)
	return db
}

func GetDbDSN(dbConfig config.DatabaseConfig) string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		dbConfig.Host, dbConfig.User, dbConfig.Password,
		dbConfig.DBName, dbConfig.Port, dbConfig.SSLMode,
	)
}
