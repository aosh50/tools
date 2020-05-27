package web

import (
	"github.com/jinzhu/gorm"
)

type BaseModel struct {
	gorm.Model
	UserID uint
}
type BaseModelList []BaseModel

func GetDB() *gorm.DB {
	// OVERWRITE
	var db *gorm.DB
	return db
}

func (Items *BaseModelList) Get(user uint) {

	err := GetDB().Where(&BaseModel{UserID: user}, user).Find(&Items).Error
	if err != nil { /* Should do stuff */
	}
}
func (Item *BaseModel) Validate() (string, bool) {

	return "Requirement passed", true
}

func (Item *BaseModel) Create() (bool, string) {
	if resp, ok := Item.Validate(); !ok {
		return ok, resp
	}
	GetDB().Create(&Item)

	if Item.ID <= 0 {
		return false, "Failed to create, connection error."
	}
	return true, "Successfylly been created"

}
func (Item *BaseModel) GetByID(ID uint) {

	err := GetDB().First(&Item, ID).Error

	if err != nil { /* Should do stuff */
	}

}
func (i *BaseModel) Put() (bool, string) {

	GetDB().Save(&i)
	return true, "Updated"
}

func (Items *BaseModelList) Search(q map[string]string, user uint) {

	tx := GetDB().Where(&BaseModel{UserID: user}, user)
	for k, v := range q {
		tx = tx.Where(k, v)
	}
	err := tx.Find(&Items).Error
	if err != nil { /* Should do stuff */
	}

}

// DeleteItem ...
func (item *BaseModel) Delete() (bool, string) {

	err := GetDB().Delete(&item).Error
	if err != nil { /* Should do stuff */
	}
	return true, "Item Deleted"
}
