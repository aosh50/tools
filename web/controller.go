package web

import (
	"encoding/json"
	"net/http"
)

var Get = func(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value("user").(uint)

	Items := BaseModelList{}
	Items.Get(user)
	resp := Message(true, "success")
	resp["data"] = Items
	Respond(w, resp)
}
var Search = func(w http.ResponseWriter, r *http.Request) {

	queryString := make(map[string]string)
	user := r.Context().Value("user").(uint)
	Items := BaseModelList{}
	Items.Search(queryString, user)
	resp := Message(true, "success")
	resp["data"] = Items
	Respond(w, resp)
}

var GetByID = func(w http.ResponseWriter, r *http.Request) {

	ID, err := IdFromUrl(r)
	if err != nil {
		Respond(w, Message(false, err.Error()))
		return
	}
	Item := BaseModel{}
	Item.GetByID(ID)
	resp := Message(true, "success")
	resp["data"] = Item
	Respond(w, resp)
}

var Create = func(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(uint)

	Item := BaseModel{}

	err := json.NewDecoder(r.Body).Decode(&Item)
	if err != nil {
		Respond(w, Message(false, "Error while decoding request body"))
		return
	}
	Item.UserID = user
	success, message := Item.Create()
	resp := Message(success, message)
	if success {
		resp["data"] = Item
	}
	Respond(w, resp)
}

var Put = func(w http.ResponseWriter, r *http.Request) {

	Item := BaseModel{}

	err := json.NewDecoder(r.Body).Decode(&Item)
	if err != nil {
		Respond(w, Message(false, "Error while decoding request body"))
		return
	}
	success, message := Item.Put()
	resp := Message(success, message)
	if success {
		resp["data"] = Item
	}
	Respond(w, resp)
}

var Delete = func(w http.ResponseWriter, r *http.Request) {

	ID, err := IdFromUrl(r)
	if err != nil {
		Respond(w, Message(false, err.Error()))
		return
	}
	Item := BaseModel{}
	Item.GetByID(ID)
	data, msg := Item.Delete()
	resp := Message(true, msg)
	resp["data"] = data
	Respond(w, resp)
}

func InvoiceActions() RestActions {
	return RestActions{
		Create:  Create,
		Get:     Get,
		GetByID: GetByID,
		Put:     Put,
		Search:  Search,
	}
}
