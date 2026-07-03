package migrations

import (
	"fmt"
	"log"

	"gin-mysql-demo/database"
	"gin-mysql-demo/models"
)

func AutoMigrate() {
	err := database.DB.AutoMigrate(
		&models.User{},
	)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	fmt.Println("✅ Database migration completed successfully!")
}

func DropTables() {
	err := database.DB.Migrator().DropTable(
		&models.User{},
	)
	if err != nil {
		log.Fatalf("Drop tables failed: %v", err)
	}
	fmt.Println("✅ Tables dropped successfully!")
}
