package web

import (
	"github.com/google/martian/log"
	"github.com/jinzhu/gorm"
)

type BaseModel struct {
	gorm.Model
}

func (base *BaseModel) Create() (*BaseModel, bool, string) {
	log.Info("Inserting Base model")
	return base, true, "Inserted"
}
