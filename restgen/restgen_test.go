package restgen

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

func TestGenControllerCreate(t *testing.T) {
	modelName, _ := testModel()
	result := genControllerCreate(modelName)
	log.Info(result)
}

func TestGenControllerGet(t *testing.T) {
	modelName, _ := testModel()
	result := genControllerGet(modelName)
	log.Info(result)
}

func TestGenControllerDelete(t *testing.T) {
	modelName, _ := testModel()
	result := genControllerDelete(modelName)
	log.Info(result)
}

func TestGenControllerUpdate(t *testing.T) {
	modelName, _ := testModel()
	result := genControllerUpdate(modelName)
	log.Info(result)
}
func TestGenControllerActions(t *testing.T) {
	modelName, _ := testModel()
	result := genControllerRestActions(modelName)
	log.Info(result)
}

func TestGenControllerTop(t *testing.T) {
	result := genControllerTop("swanwater-go")
	log.Info(result)
}
func TestTs(t *testing.T) {
	modelName, props := testModel()
	result := genTsModel(modelName, props)
	log.Info(result)
}
func TestTsApis(t *testing.T) {
	modelName, _ := testModel()
	result := genTsApiEndpoints(modelName)
	log.Info(result)
}
func TestReduxActions(t *testing.T) {
	modelName, _ := testModel()
	result := genReduxActions(modelName)
	log.Info(result)
}
func testModel() (string, map[string]string) {
	modelName := "Paddock"
	modelPropMap := make(map[string]string)
	modelPropMap["Name"] = "string"
	modelPropMap["Acres"] = "float64"
	modelPropMap["FarmID"] = "uint"
	modelPropMap["CentreLat"] = "float64"
	modelPropMap["CentreLng"] = "float64"
	modelPropMap["Notes"] = "string"

	return modelName, modelPropMap
}
