package vault

import (
	"log"

	m "github.com/stevensopilidis/dora/vault/models"
	"gorm.io/gorm"
)

// func for migrating the database
func migrateVault(db *gorm.DB) {
	err := db.AutoMigrate(&m.ServiceModel{})
	if err != nil {
		log.Panicf("Could not migrate vault: {%s}", err.Error())
	}
	log.Println("Vault migrated successfully")
}
