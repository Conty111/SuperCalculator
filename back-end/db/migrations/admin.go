package main

import (
	"database/sql"
	"fmt"
	"github.com/spf13/cobra"
	_ "gorm.io/driver/postgres"
	_ "gorm.io/driver/sqlite"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "migration",
	Short: "A simple database migration tool",
	Run: func(cmd *cobra.Command, args []string) {
		databaseType, _ := cmd.Flags().GetString("database")
		connectionString, _ := cmd.Flags().GetString("conn")
		userEmail, _ := cmd.Flags().GetString("email")

		var db *sql.DB
		var err error

		switch databaseType {
		case "postgres":
			db, err = sql.Open("postgres", connectionString)
		case "sqlite":
			db, err = sql.Open("sqlite3", connectionString)
		default:
			fmt.Println("Unsupported database type. Please provide a valid database type.")
			os.Exit(1)
		}

		if err != nil {
			panic(err.Error())
		}
		defer db.Close()

		// Подготовка SQL запроса с параметризированным значением email
		query := "UPDATE users SET role = 'admin' WHERE email = $1"

		// Выполнение SQL запроса с передачей email в качестве аргумента
		_, err = db.Exec(query, userEmail)
		if err != nil {
			panic(err.Error())
		}

		fmt.Println("Миграция выполнена успешно для пользователя с email:", userEmail)
	},
}

func main() {
	rootCmd.Flags().String("database", "postgres", "Database type (e.g. postgres)")
	rootCmd.Flags().String("conn", "", "Database connection string")
	rootCmd.Flags().String("email", "", "User's email for migration")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
