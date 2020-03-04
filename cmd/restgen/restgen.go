package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"unicode"
)

var VALID_TYPES = []string{"string", "bool", "uint", "int", "int32", "int64", "float32", "float64", "time.Time"}

func main() {
	newModel()
}

func newModel() {
	appName := readString("App name")
	name := readString("Name of model:")
	props := make(map[string]string)
	cont := true
	for cont {
		input := readString("Add prop: PROPERTYNAME PROPERYTYPE or type 'exit' to continue")
		if input == "exit" {
			cont = false
			continue
		}
		propAndType := strings.Split(input, " ")
		if len(propAndType) != 2 {
			fmt.Println("Incorrect format, try again")
			continue
		}
		propType := propAndType[1]
		prop := propAndType[0]

		if !validateInputType(propType) {
			fmt.Printf("Incorrect type: %s, try again\n", propType)
			continue
		}
		props[prop] = propType

	}
	fileContents := genModel(name, props)
	controllerContents := genController(name, appName)
	err := ensureBaseDir(fmt.Sprintf("gen/%s", appName))
	if err != nil {
		log.Fatalln("Ensure error")
		log.Fatalln(err.Error())
	}
	err = WriteToFile(fmt.Sprintf("gen/%s/%s.go", appName, name), fileContents)
	if err != nil {
		log.Fatalln(err.Error())
	}
	err = WriteToFile(fmt.Sprintf("gen/%s/%sController.go", appName, name), controllerContents)
	if err != nil {
		log.Fatalln(err.Error())
	}

}

func genModel(name string, props map[string]string) string {
	struc := genStruct(name, props)
	viewModel := genViewModel(name, props)
	validate := genValidate(name)
	create := genCreate(name)
	delete := genDelete(name)
	update := genUpdate(name)
	getByID := genGetByID(name)
	top := fmt.Sprintf("package models\n\nimport \"github.com/jinzhu/gorm\"\n\n")
	full := fmt.Sprintf("%s%s%s%s%s%s%s%s", top, struc, viewModel, validate, create, delete, update, getByID)
	return full
}

func validateInputType(input string) bool {
	valid := true
	firstChar := rune(input[0])
	if !unicode.IsUpper(firstChar) {
		valid = stringInSlice(input, VALID_TYPES)
	}
	return valid
}
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
func genStruct(model string, props map[string]string) string {
	res := fmt.Sprintf("type %s struct {\n\tgorm.Model\n", model)
	for k, v := range props {
		line := fmt.Sprintf("\t%s %s\n", k, v) // `json:\"%s\"`
		res = fmt.Sprintf("%s%s", res, line)
	}
	res = fmt.Sprintf("%s}\n", res)
	return res
}
func genViewModel(model string, props map[string]string) string {
	res := fmt.Sprintf("type %sViewModel struct {\n\tID uint\n", model)
	for k, v := range props {
		line := fmt.Sprintf("\t%s %s,\n", k, v)
		res = fmt.Sprintf("%s%s", res, line)
	}
	res = fmt.Sprintf("%s}\n", res)

	res = fmt.Sprintf("%s\nfunc (m *%s) ViewModel() %sViewModel {\n\tvm := %sViewModel{\n", res, model, model, model)
	for k, _ := range props {
		line := fmt.Sprintf("\t\t%s: m.%s\n", k, k)
		res = fmt.Sprintf("%s%s", res, line)
	}
	res = fmt.Sprintf("%s\t}\n\treturn vm\n}\n", res)

	return res
}

func genValidate(model string) string {
	return fmt.Sprintf("func (m *%s) Validate() (bool, string) {\n\treturn true, \"Valid\"\n}\n", model)
}
func genCreate(model string) string {
	res := fmt.Sprintf("func (m *%s) Create() (bool, string) {\n\tif resp, ok := p.Validate(); !ok {\n\t\treturn ok, resp\n\t}\n", model)
	res = fmt.Sprintf("%s\n\terr := GetDB().Create(m).Error\n\tif err != nil {\n\t\treturn false, err.Error()\n\t}\n", res)
	res = fmt.Sprintf("%s\n\tif p.ID <= 0 {\n\t\treturn false, \"Failed to create %s, connection error.\"\n\t}\n", res, model)
	res = fmt.Sprintf("%s\n\treturn true, \"%s has been created\"\n}\n", res, model)
	return res

}
func genDelete(model string) string {
	res := fmt.Sprintf("func (m *%s) Delete() (bool, string) {\n", model)
	res = fmt.Sprintf("%s\n\terr := GetDB().Delete(m).Error\n\tif err != nil {\n\t\treturn false, err.Error()\n\t}\n", res)
	res = fmt.Sprintf("%s\n\treturn true, \"%s has been deleted\"\n}\n", res, model)
	return res

}

func genUpdate(model string) string {
	res := fmt.Sprintf("func (m *%s) Update() (bool, string) {\n", model)
	res = fmt.Sprintf("%s\n\terr := GetDB().Update(m).Error\n\tif err != nil {\n\t\treturn false, err.Error()\n\t}\n", res)
	res = fmt.Sprintf("%s\n\treturn true, \"%s has been updated\"\n}\n", res, model)
	return res

}
func genGetByID(model string) string {
	lower := strings.ToLower(model)
	res := fmt.Sprintf("func Get%sByID(ID uint) *%s {\n", model, model)
	res = fmt.Sprintf("%s\n\t%s := %s{}\n", res, lower, model)
	res = fmt.Sprintf("%s\n\terr := GetDB().First(&%s, ID).Error\n\tif err != nil {\n\t\treturn nil\n\t}\n", res, lower)
	res = fmt.Sprintf("%s\n\treturn &%s\n}\n", res, lower)
	return res

}

func genControllerCreate(model string) string {
	lower := strings.ToLower(model)

	res := fmt.Sprintf("var Create%s = func(w http.ResponseWriter, r *http.Request) {\n", model)
	res = fmt.Sprintf("%s\t%s := models.%s\n", res, lower, model)
	res = fmt.Sprintf("%s\terr := json.NewDecoder(r.Body).Decode(&%s)\n\tif err != nil {\n\t\tu.Respond(w, u.Message(false, \"Error while decoding request body\"))\n\t\treturn\n\t}\n", res, lower)
	res = fmt.Sprintf("%s\tsuccess, message := %s.Create()\n", res, lower)
	res = fmt.Sprintf("%s\tresp := u.Message(success, message)\n\tif success {\n\t\tresp[\"data\"] = %s.ViewModel()\n\t}\n\tu.Respond(w, resp)\n}\n", res, lower)
	return res
}

func genControllerGet(model string) string {
	lower := strings.ToLower(model)

	res := fmt.Sprintf("var Get%s = func(w http.ResponseWriter, r *http.Request) {\n", model)
	res = fmt.Sprintf("%s\tkey, ok := r.URL.Query()[\"id\"]\n", res)
	res = fmt.Sprintf("%s\tif !ok || len(key[0]) != 1 {\n\t\tu.Respond(w, u.Message(false, \"Error while decoding request body\"))\n\t\treturn\n\t}\n", res)
	res = fmt.Sprintf("%s\tID, err := strconv.ParseUint(key[0], 10, 32)\n\tif err != nil {\n\t\tu.Respond(w, u.Message(false, \"Error while decoding request body\"))\n\t\treturn\n\t}\n", res)
	res = fmt.Sprintf("%s\tresp := u.Message(true, message)\n", res)
	res = fmt.Sprintf("%s\t%ss := models.Get%sByID(uint(ID))\n", res, lower, model)
	res = fmt.Sprintf("%s\tresp[\"data\"] = %ss\n\tu.Respond(w, resp)\n}\n", res, lower)

	return res
}
func genControllerDelete(model string) string {
	lower := strings.ToLower(model)

	res := fmt.Sprintf("var Delete%s = func(w http.ResponseWriter, r *http.Request) {\n", model)
	res = fmt.Sprintf("%s\tkey, ok := r.URL.Query()[\"id\"]\n", res)
	res = fmt.Sprintf("%s\tif !ok || len(key[0]) != 1 {\n\t\tu.Respond(w, u.Message(false, \"Error while decoding request body\"))\n\t\treturn\n\t}\n", res)
	res = fmt.Sprintf("%s\tID, err := strconv.ParseUint(key[0], 10, 32)))\n\tif err != nil {\n\t\tu.Respond(w, u.Message(false, \"Error while decoding request body\"))\n\t\treturn\n\t}\n", res)
	res = fmt.Sprintf("%s\t%s := models.Get%sByID(uint(ID))\n", res, lower, model)
	res = fmt.Sprintf("%s\tsuccess, message := %s.Delete()\n", res, lower)
	res = fmt.Sprintf("%s\tresp := u.Message(success, message)\n\tu.Respond(w, resp)\n}\n", res)

	return res
}
func genControllerUpdate(model string) string {
	lower := strings.ToLower(model)

	res := fmt.Sprintf("var Update%s = func(w http.ResponseWriter, r *http.Request) {\n", model)
	res = fmt.Sprintf("%s\t%s := models.%s\n", res, lower, model)
	res = fmt.Sprintf("%s\terr := json.NewDecoder(r.Body).Decode(&%s)\n\tif err != nil {\n\t\tu.Respond(w, u.Message(false, \"Error while decoding request body\"))\n\t\treturn\n\t}\n", res, lower)
	res = fmt.Sprintf("%s\tsuccess, message := %s.Update()\n", res, lower)
	res = fmt.Sprintf("%s\tresp := u.Message(success, message)\n\tif success {\n\t\tresp[\"data\"] = %s.ViewModel()\n\t}\n\tu.Respond(w, resp)\n}\n", res, lower)

	return res
}

func genControllerRestActions(model string) string {
	res := fmt.Sprintf("func %sActions() RestActions {\n\treturn RestActions{\n", model)
	res = fmt.Sprintf("%s\t\tGet: Get%s,\n\t\tCreate: Create%s,\n\t\tEdit: Update%s,\n\t\tDelete: Delete%s,\n\t}\n}", res, model, model, model, model)
	return res
}

func genControllerTop(appName string) string {
	res := fmt.Sprintf("package controllers\n\nimport (\n\t\"encoding/json\"\n\t\"net/http\"\n\t\"strconv\"\n\t\"%s/models\"\n\tu \"%s/utils\"\n)\n", appName, appName)
	return res
}

func genController(model string, appName string) string {
	top := genControllerTop(appName)
	create := genControllerCreate(model)
	edit := genControllerUpdate(model)
	get := genControllerGet(model)
	delete := genControllerDelete(model)
	actions := genControllerRestActions(model)
	res := fmt.Sprintf("%s%s%s%s%s%s", top, create, edit, get, delete, actions)
	return res
}

func readString(message string) string {
	//reading a string
	reader := bufio.NewReader(os.Stdin)
	var name string
	fmt.Println(message)
	name, _ = reader.ReadString('\n')
	return strings.TrimSpace(name)
}
func WriteToFile(filename string, data string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, data)
	if err != nil {
		return err
	}
	return file.Sync()
}

func ensureBaseDir(fpath string) error {
	// baseDir := path.Dir(fpath)
	info, err := os.Stat(fpath)
	if err == nil && info.IsDir() {
		return nil
	}
	return os.MkdirAll(fpath, 0755)
}
