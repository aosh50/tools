package main

import (
	"fmt"
	"os"
	"bufio"
	log "github.com/sirupsen/logrus"
	"strings"
	"unicode"
)

var VALID_TYPES = []string{"string", "bool", "uint", "int", "int32", "int64", "float32", "float64", "time.Time"}

func main() {
	newModel()
}

func newModel() {
	log.Info("Creating new model")
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

	f, _ := os.Create(fmt.Sprintf("/gen/%s.go", name))
	defer f.Close()
	f.WriteString(fileContents)

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
	fmt.Println(full)
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
		line := fmt.Sprintf("\t%s %s\n", k, v)
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
func readString(message string) string {
	//reading a string
	reader := bufio.NewReader(os.Stdin)
	var name string
	fmt.Println(message)
	name, _ = reader.ReadString('\n')
	return strings.TrimSpace(name)
}
