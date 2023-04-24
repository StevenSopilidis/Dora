package vault

import (
	"context"
	"fmt"
	"log"
	"time"

	e "github.com/stevensopilidis/dora/errors"
	m "github.com/stevensopilidis/dora/vault/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	maxDbOperationWaitingTime = 10 * time.Second
)

type InitializeVaultConfig struct {
	Host string
	Port int
	Db   string
	User string
	Pass string
}

type Vault struct {
	db *gorm.DB
}

// gets all services from vault
func (v *Vault) GetServices(ctx context.Context) (error, []m.ServiceModel) {
	var services []m.ServiceModel
	ctx, cancel := context.WithTimeout(ctx, maxDbOperationWaitingTime)
	defer cancel()
	result := v.db.WithContext(ctx).Find(&services)
	if result.Error == gorm.ErrRecordNotFound {
		return &e.ServiceNotFoundError{}, nil
	}
	if result.Error != nil {
		return result.Error, nil
	}
	return nil, services
}

// adds a service to vault
func (v *Vault) AddService(ctx context.Context, service *m.ServiceModel) error {
	ctx, cancel := context.WithTimeout(ctx, maxDbOperationWaitingTime)
	defer cancel()
	result := v.db.WithContext(ctx).Create(service)
	return result.Error
}

func (v *Vault) RemoveService(ctx context.Context, id uint) error {
	ctx, cancel := context.WithTimeout(ctx, maxDbOperationWaitingTime)
	defer cancel()
	result := v.db.WithContext(ctx).Delete(&m.ServiceModel{}, id)
	if result.Error == gorm.ErrRecordNotFound {
		return &e.ServiceNotFoundError{}
	}
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (v *Vault) UpdateService(ctx context.Context, service *m.ServiceModel) error {
	ctx, cancel := context.WithTimeout(ctx, maxDbOperationWaitingTime)
	defer cancel()
	result := v.db.Save(service)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func InitializeVault(config *InitializeVaultConfig) *Vault {
	// first create the database
	dsn := fmt.Sprintf("host=%s user=%s password=%s port=%d sslmode=disable ", config.Host, config.User, config.Pass, config.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panicf("Could not initialize connection to vault: {%s}", err.Error())
	}
	// check if specified database exists
	var result int64
	query := fmt.Sprintf("SELECT 1 FROM pg_database WHERE datname = '%s' LIMIT 1;", config.Db)
	db.Raw(query).Scan(&result)
	// if not create it
	if result != 1 {
		// The database does not exist, create it
		cmd := fmt.Sprintf("CREATE DATABASE %s;", config.Db)
		db.Exec(cmd)
	}
	dbSQL, err := db.DB()
	if err != nil {
		panic("Failed to get DB: " + err.Error())
	}
	dbSQL.Close()
	// connect to the newly created database
	dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", config.Host, config.User, config.Pass, config.Db, config.Port)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panicf("Could not initialize connection to vault: {%s}", err.Error())
	}
	log.Println("Connection to vault initialized successfully")
	migrateVault(db)
	return &Vault{
		db,
	}
}

func CloseVault(vault *Vault) {
	dbSQL, err := vault.db.DB()
	if err != nil {
		log.Panicf("Could not get database: {%s}", err.Error())
	}
	log.Println("Vault closed successfully")
	dbSQL.Close()
}
