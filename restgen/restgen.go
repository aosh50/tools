package restgen

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"unicode"
)

var ValidTypes = []string{"string", "bool", "uint", "int", "int32", "int64", "float32", "float64", "time.Time"}
var NumberTypes = []string{"uint", "int", "int32", "int64", "float32", "float64"}
var StringTypes = []string{"string"}
var DateTypes = []string{"time.Time"}
var BoolTypes = []string{"bool"}

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
	tsContents := genTs(name, props)
	lower := strings.ToLower(name)

	err := ensureBaseDir(fmt.Sprintf("gen/%s", appName))
	if err != nil {
		log.Fatalln("Ensure error")
		log.Fatalln(err.Error())
	}
	err = WriteToFile(fmt.Sprintf("gen/%s/%s.go", appName, lower), fileContents)
	if err != nil {
		log.Fatalln(err.Error())
	}
	err = WriteToFile(fmt.Sprintf("gen/%s/%sController.go", appName, lower), controllerContents)
	if err != nil {
		log.Fatalln(err.Error())
	}
	err = WriteToFile(fmt.Sprintf("gen/%s/%s.ts", appName, name), tsContents)
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
		valid = stringInSlice(input, ValidTypes)
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
		line := fmt.Sprintf("\t\t%s: m.%s,\n", k, k)
		res = fmt.Sprintf("%s%s", res, line)
	}
	res = fmt.Sprintf("%s\t}\n\treturn vm\n}\n", res)

	return res
}

func genValidate(model string) string {
	return fmt.Sprintf("func (m *%s) Validate() (bool, string) {\n\treturn true, \"Valid\"\n}\n", model)
}
func genCreate(model string) string {
	res := fmt.Sprintf("func (m *%s) Create() (bool, string) {\n\tif ok, resp := m.Validate(); !ok {\n\t\treturn ok, resp\n\t}\n", model)
	res = fmt.Sprintf("%s\n\terr := GetDB().Create(m).Error\n\tif err != nil {\n\t\treturn false, err.Error()\n\t}\n", res)
	res = fmt.Sprintf("%s\n\tif m.ID <= 0 {\n\t\treturn false, \"Failed to create %s, connection error.\"\n\t}\n", res, model)
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
	res = fmt.Sprintf("%s\n\terr := GetDB().Save(&m).Error\n\tif err != nil {\n\t\treturn false, err.Error()\n\t}\n", res)
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
	res = fmt.Sprintf("%s\t%s := models.%s{}\n", res, lower, model)
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
	res = fmt.Sprintf("%s\tresp := u.Message(true, \"success\")\n", res)
	res = fmt.Sprintf("%s\t%ss := models.Get%sByID(uint(ID))\n", res, lower, model)
	res = fmt.Sprintf("%s\tresp[\"data\"] = %ss\n\tu.Respond(w, resp)\n}\n", res, lower)

	return res
}
func genControllerDelete(model string) string {
	lower := strings.ToLower(model)

	res := fmt.Sprintf("var Delete%s = func(w http.ResponseWriter, r *http.Request) {\n", model)
	res = fmt.Sprintf("%s\tkey, ok := r.URL.Query()[\"id\"]\n", res)
	res = fmt.Sprintf("%s\tif !ok || len(key[0]) != 1 {\n\t\tu.Respond(w, u.Message(false, \"Error while decoding request body\"))\n\t\treturn\n\t}\n", res)
	res = fmt.Sprintf("%s\tID, err := strconv.ParseUint(key[0], 10, 32)\n\tif err != nil {\n\t\tu.Respond(w, u.Message(false, \"Error while decoding request body\"))\n\t\treturn\n\t}\n", res)
	res = fmt.Sprintf("%s\t%s := models.Get%sByID(uint(ID))\n", res, lower, model)
	res = fmt.Sprintf("%s\tsuccess, message := %s.Delete()\n", res, lower)
	res = fmt.Sprintf("%s\tresp := u.Message(success, message)\n\tu.Respond(w, resp)\n}\n", res)

	return res
}
func genControllerUpdate(model string) string {
	lower := strings.ToLower(model)

	res := fmt.Sprintf("var Update%s = func(w http.ResponseWriter, r *http.Request) {\n", model)
	res = fmt.Sprintf("%s\t%s := models.%s{}\n", res, lower, model)
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

func contains(v string, a []string) bool {
	for _, i := range a {
		if i == v {
			return true
		}
	}
	return false
}

func goTypeToTs(gotype string) string {
	if contains(gotype, NumberTypes) {
		return "number"
	}
	if contains(gotype, StringTypes) {
		return "string"
	}
	if contains(gotype, DateTypes) {
		return "Date"
	}
	if contains(gotype, BoolTypes) {
		return "boolean"
	}
	return gotype
}

func genTs(model string, props map[string]string) string {
	res := genTsModel(model, props)
	res = fmt.Sprintf("%s%s", res, genTsApiEndpoints(model))
	res = fmt.Sprintf("%s%s", res, genReduxActions(model))
	return res
}

func genTsModel(model string, props map[string]string) string {
	res := fmt.Sprintf("export interface %s {\n", model)
	for k, v := range props {
		line := fmt.Sprintf("\t%s: %s;\n", k, goTypeToTs(v))
		res = fmt.Sprintf("%s%s", res, line)
	}
	res = fmt.Sprintf("%s}\n", res)
	return res
}

func genTsApiEndpoints(model string) string {
	lower := strings.ToLower(model)

	res := fmt.Sprintf("\tget%s: (id: number): AxiosPromise<M.ApiResponse<M.%s>> => base().get(`/%s?id=${id}`),\n", model, model, lower)
	res = fmt.Sprintf("%s\tcreate%s: (m: %s): AxiosPromise<M.ApiResponse<M.%s>> => base().post(`/%s`, m),\n", res, model, model, model, lower)
	res = fmt.Sprintf("%s\tupdate%s: (m: %s): AxiosPromise<M.ApiResponse<M.%s>> => base().patch(`/%s`, m),\n", res, model, model, model, lower)
	res = fmt.Sprintf("%s\tdelete%s: (id: number): AxiosPromise<M.ApiResponse<M.%s>> => base().delete(`/%s?id=${id}`),\n", res, model, model, lower)

	return res
}

func genReduxActions(model string) string {
	res := fmt.Sprintf("const get%sCreator = ac.async<number, M.ApiResponse<M.%s>, M.ApiResponse<string>>('Get%s');\n", model, model, model)
	res = fmt.Sprintf("%sconst create%sCreator = ac.async<M.%s, M.ApiResponse<M.%s>, M.ApiResponse<string>>('Create%s');\n", res, model, model, model, model)
	res = fmt.Sprintf("%sconst update%sCreator = ac.async<M.%s, M.ApiResponse<M.%s>, M.ApiResponse<string>>('Update%s');\n", res, model, model, model, model)
	res = fmt.Sprintf("%sconst delete%sCreator = ac.async<number, M.ApiResponse<string>, M.ApiResponse<string>>('Delete%s');\n", res, model, model)

	res = fmt.Sprintf("%s\nget%s: A.wrapAsyncWorker(\n\tget%sCreator,\n\t(params, dispatch) => A.api.get%s(params).then(resp => resp.data)),\n", res, model, model, model)
	res = fmt.Sprintf("%sdelete%s: A.wrapAsyncWorker(\n\tdelete%sCreator,\n\t(params, dispatch) => A.api.delete%s(params).then(resp => resp.data)),\n", res, model, model, model)
	res = fmt.Sprintf("%screate%s: A.wrapAsyncWorker(\n\tcreate%sCreator,\n\t(params, dispatch) => A.api.create%s(params).then(resp => resp.data)),\n", res, model, model, model)
	res = fmt.Sprintf("%supdate%s: A.wrapAsyncWorker(\n\tupdate%sCreator,\n\t(params, dispatch) => A.api.update%s(params).then(resp => resp.data)),\n", res, model, model, model)

	res = fmt.Sprintf("%s\n\n%s", res, reduxAction("get", model))
	res = fmt.Sprintf("%s%s", res, reduxAction("create", model))
	res = fmt.Sprintf("%s%s", res, reduxAction("delete", model))
	res = fmt.Sprintf("%s%s\n", res, reduxAction("update", model))

	return res

}

func reduxAction(action string, model string) string {
	res := fmt.Sprintf(".case(%s%sCreator.started, (state, payload) => { return state; })\n", action, model)
	res = fmt.Sprintf("%s.case(%s%sCreator.done, (state, payload) => { return state; })\n", res, action, model)
	res = fmt.Sprintf("%s.case(%s%sCreator.failed, (state, payload) => { return state; })\n", res, action, model)
	return res
}
