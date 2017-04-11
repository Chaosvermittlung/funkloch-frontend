package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func showItemlist(w http.ResponseWriter, token string) {
	var ilp ItemListPage
	tp := "templates/item/list.html"
	ilp.Default.Sidebar = BuildSidebar(ItemsActive)
	ilp.Default.Pagename = "Item List"

	err := sendauthorizedHTTPRequest("GET", "storeitem/list", token, nil, &ilp.Items)
	if err != nil {
		ilp.Default.Message = BuildMessage(errormessage, "Error creating item/list request: "+err.Error())
		showtemplate(w, tp, ilp)
		return
	}
	showtemplate(w, tp, ilp)
}

func showItemAddForm(w http.ResponseWriter, token string) {
	var iap ItemAddPage
	tp := "templates/item/add.html"
	iap.Default.Sidebar = BuildSidebar(ItemsActive)
	iap.Default.Pagename = "Add Item"
	err := sendauthorizedHTTPRequest("GET", "store/list", token, nil, &iap.Stores)
	if err != nil {
		iap.Default.Message = BuildMessage(errormessage, "Error sending Stores request: "+err.Error())
		showtemplate(w, tp, iap)
		return
	}
	err = sendauthorizedHTTPRequest("GET", "equipment/list", token, nil, &iap.Equipments)
	if err != nil {
		iap.Default.Message = BuildMessage(errormessage, "Error sending Equipment request: "+err.Error())
		showtemplate(w, tp, iap)
		return
	}
	showtemplate(w, tp, iap)
}

func saveNewItem(w http.ResponseWriter, r *http.Request, token string) {
	var iap ItemAddPage
	tp := "templates/item/add.html"
	var inp ItemNewPAge
	tp2 := "templates/item/new.html"
	iap.Default.Sidebar = BuildSidebar(ItemsActive)
	iap.Default.Pagename = "Add Item"
	e := r.FormValue("equipment")
	s := r.FormValue("store")
	c := r.FormValue("count")
	sid, err := strconv.Atoi(s)
	if err != nil {
		iap.Default.Message = BuildMessage(errormessage, "Error converting Store ID"+err.Error())
		showtemplate(w, tp, iap)
		return
	}
	eid, err := strconv.Atoi(e)
	if err != nil {
		iap.Default.Message = BuildMessage(errormessage, "Error converting Equipment ID"+err.Error())
		showtemplate(w, tp, iap)
		return
	}
	co, err := strconv.Atoi(c)
	if err != nil {
		iap.Default.Message = BuildMessage(errormessage, "Error converting Count"+err.Error())
		showtemplate(w, tp, iap)
		return
	}

	inp.Default.Sidebar = BuildSidebar(ItemsActive)
	for i := 0; i < co; i++ {
		var si StoreItem
		var so Store
		var ini ItemNewItem
		si.EquipmentID = eid
		si.StoreID = sid
		b := new(bytes.Buffer)
		encoder := json.NewEncoder(b)
		encoder.Encode(si)
		var res StoreItem
		err = sendauthorizedHTTPRequest("POST", "storeitem/", token, b, &res)
		if err != nil {
			inp.Default.Message = BuildMessage(errormessage, "Error sending StoreItem request: "+err.Error())
			showtemplate(w, tp2, inp)
			return
		}
		err = sendauthorizedHTTPRequest("GET", "store/"+s, token, nil, &so)
		if err != nil {
			inp.Default.Message = BuildMessage(errormessage, "Error sending Store request: "+err.Error())
			showtemplate(w, tp2, inp)
			return
		}
		id := createEAN13(res.StoreItemID)
		ini.ID = id
		ini.Store = so.Name
		inp.IDs = append(inp.IDs, ini)
	}
	showtemplate(w, tp2, inp)
}

func showItemEditForm(w http.ResponseWriter, r *http.Request, token string) {
	var iep ItemEditPage
	tp := "templates/item/edit.html"
	iep.Default.Sidebar = BuildSidebar(ItemsActive)
	iep.Default.Pagename = "Edit Item"
	id := r.FormValue("itemid")
	err := sendauthorizedHTTPRequest("GET", "storeitem/"+id, token, nil, &iep.Ite)
	if err != nil {
		iep.Default.Message = BuildMessage(errormessage, "Error sending Event request: "+err.Error())
		showtemplate(w, tp, iep)
		return
	}
	err = sendauthorizedHTTPRequest("GET", "store/list", token, nil, &iep.Stores)
	if err != nil {
		iep.Default.Message = BuildMessage(errormessage, "Error sending Stores request: "+err.Error())
		showtemplate(w, tp, iep)
		return
	}
	err = sendauthorizedHTTPRequest("GET", "equipment/list", token, nil, &iep.Equipments)
	if err != nil {
		iep.Default.Message = BuildMessage(errormessage, "Error sending Equipment request: "+err.Error())
		showtemplate(w, tp, iep)
		return
	}
	showtemplate(w, tp, iep)
}

func deleteItem(w http.ResponseWriter, r *http.Request, token string) {

}

func patchItem(w http.ResponseWriter, r *http.Request, token string) {
	var iep ItemEditPage
	tp := "templates/item/edit.html"
	iep.Default.Sidebar = BuildSidebar(ItemsActive)
	iep.Default.Pagename = "Edit Item"
	id := r.FormValue("itemid")
	e := r.FormValue("equipment")
	s := r.FormValue("store")
	sid, err := strconv.Atoi(s)
	if err != nil {
		iep.Default.Message = BuildMessage(errormessage, "Error converting Store ID"+err.Error())
		showtemplate(w, tp, iep)
		return
	}
	eid, err := strconv.Atoi(e)
	if err != nil {
		iep.Default.Message = BuildMessage(errormessage, "Error converting Equipment ID"+err.Error())
		showtemplate(w, tp, iep)
		return
	}
	var si StoreItem
	si.EquipmentID = eid
	si.StoreID = sid
	b := new(bytes.Buffer)
	encoder := json.NewEncoder(b)
	encoder.Encode(si)

	err = sendauthorizedHTTPRequest("PATCH", "storeitem/"+id, token, b, nil)
	if err != nil {
		iep.Default.Message = BuildMessage(errormessage, "Error sending StoreItem request: "+err.Error())
		showtemplate(w, tp, iep)
		return
	}
	http.Redirect(w, r, "/item", http.StatusSeeOther)
}

func viewItem(w http.ResponseWriter, r *http.Request, token string) {
	var ivp ItemViewPage
	tp := "templates/item/view.html"
	ivp.Default.Sidebar = BuildSidebar(ItemsActive)
	ivp.Default.Pagename = "View Item"
	id := r.FormValue("id")
	err := sendauthorizedHTTPRequest("GET", "storeitem/"+id, token, nil, &ivp.Ite)
	if err != nil {
		ivp.Default.Message = BuildMessage(errormessage, "Error creating event request: "+err.Error())
		showtemplate(w, tp, ivp)
		return
	}
	ivp.EAN = createEAN13(ivp.Ite.StoreItem.StoreItemID)
	showtemplate(w, tp, ivp)
}

func itemLabel(w http.ResponseWriter, r *http.Request, tonken string) {

	s := r.FormValue("store")
	e := r.FormValue("EAN")
	w.Header().Set("Content-type", "application/pdf")
	err := createlabel(e, s, w)
	if err != nil {
		fmt.Println(err)
	}
}

func itemHandler(w http.ResponseWriter, r *http.Request) {

	token, err := GetCookie(r, "token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	a := r.FormValue("action")
	switch a {
	case "add":
		showItemAddForm(w, token)
	case "save":
		saveNewItem(w, r, token)
	case "edit":
		showItemEditForm(w, r, token)
	case "patch":
		patchItem(w, r, token)
	case "delete":
		deleteItem(w, r, token)
	case "view":
		viewItem(w, r, token)
	case "label":
		itemLabel(w, r, token)
	default:
		showItemlist(w, token)
	}

}
