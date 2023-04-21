package models

import (
	r "github.com/stevensopilidis/dora/registry"
	"gorm.io/gorm"
)

type ServiceModel struct {
	gorm.Model
	r.Service
}
