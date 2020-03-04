package main

import (
	"testing"
	log "github.com/sirupsen/logrus"

)

func TestGenViewModel(t *testing.T) {
	
	modelName, modelPropMap := testModel()
	result := genViewModel(modelName, modelPropMap)

	log.Info(result)

}
func TestGenValidate(t *testing.T) {
	
	modelName, _ := testModel()
	result := genValidate(modelName)

	log.Info(result)
}

func TestGenCreate(t *testing.T) {
	modelName, _ := testModel()
	result := genCreate(modelName)
	log.Info(result)
}
func TestGenDelete(t *testing.T) {
	modelName, _ := testModel()
	result := genDelete(modelName)
	log.Info(result)
}
func TestGenUpdate(t *testing.T) {
	modelName, _ := testModel()
	result := genUpdate(modelName)
	log.Info(result)
}
func TestGenGetByID(t *testing.T) {
	modelName, _ := testModel()
	result := genGetByID(modelName)
	log.Info(result)
}

func TestGenModel(t *testing.T) {
	modelName, props := testModel()
	result := genModel(modelName, props)
	log.Info(result)
}


func testModel() (string, map[string]string) {
	modelName := "Paddock"
	modelPropMap := make(map[string]string)
	modelPropMap["Name"] = "string"
	modelPropMap["Acres"] = "float32"
	modelPropMap["FarmID"] = "uint"
	return modelName, modelPropMap
}