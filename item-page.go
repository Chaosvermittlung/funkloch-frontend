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

	err := sendauthorizedHTTPRequest("GET", "item/list", token, nil, &ilp.Items)
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
	err := sendauthorizedHTTPRequest("GET", "equipment/list", token, nil, &iap.Equipments)
	if err != nil {
		iap.Default.Message = BuildMessage(errormessage, "Error sending Equipment request: "+err.Error())
		showtemplate(w, tp, iap)
		return
	}
	err = sendauthorizedHTTPRequest("GET", "box/list", token, nil, &iap.Boxes)
	if err != nil {
		iap.Default.Message = BuildMessage(errormessage, "Error sending Box request: "+err.Error())
		showtemplate(w, tp, iap)
		return
	}
	err = sendauthorizedHTTPRequest("GET", "store/list", token, nil, &iap.Stores)
	if err != nil {
		iap.Default.Message = BuildMessage(errormessage, "Error sending Box request: "+err.Error())
		showtemplate(w, tp, iap)
		return
	}
	showtemplate(w, tp, iap)
}

func saveNewItem(w http.ResponseWriter, r *http.Request, token string) {
	var iap ItemAddPage
	tp := "templates/item/add.html"
	var inp ItemNewPage
	tp2 := "templates/item/new.html"
	iap.Default.Sidebar = BuildSidebar(ItemsActive)
	iap.Default.Pagename = "Add Item"
	e := r.FormValue("equipment")
	bo := r.FormValue("box")
	c := r.FormValue("count")
	bid, err := strconv.Atoi(bo)
	if err != nil {
		iap.Default.Message = BuildMessage(errormessage, "Error converting Box ID"+err.Error())
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
		var si Item
		var bx boxResponse
		var ini ItemNewItem
		si.EquipmentID = eid
		si.BoxID = bid
		b := new(bytes.Buffer)
		encoder := json.NewEncoder(b)
		encoder.Encode(si)
		var res Item
		err = sendauthorizedHTTPRequest("POST", "item/", token, b, &res)
		if err != nil {
			inp.Default.Message = BuildMessage(errormessage, "Error sending StoreItem request: "+err.Error())
			showtemplate(w, tp2, inp)
			return
		}
		err = sendauthorizedHTTPRequest("GET", "box/"+bo, token, nil, &bx)
		if err != nil {
			inp.Default.Message = BuildMessage(errormessage, "Error sending Store request: "+err.Error())
			showtemplate(w, tp2, inp)
			return
		}
		ini.ID = strconv.Itoa(res.Code)
		ini.Box = strconv.Itoa(bx.Box.Code)
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
	err := sendauthorizedHTTPRequest("GET", "item/"+id, token, nil, &iep.Ite)
	if err != nil {
		iep.Default.Message = BuildMessage(errormessage, "Error sending Item request: "+err.Error())
		showtemplate(w, tp, iep)
		return
	}
	err = sendauthorizedHTTPRequest("GET", "box/list", token, nil, &iep.Boxes)
	if err != nil {
		iep.Default.Message = BuildMessage(errormessage, "Error sending Box request: "+err.Error())
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
	var iep ItemEditPage
	tp := "templates/item/edit.html"
	iep.Default.Sidebar = BuildSidebar(ItemsActive)
	iep.Default.Pagename = "Edit Item"
	id := r.FormValue("itemid")
	err := sendauthorizedHTTPRequest("DELETE", "item/"+id, token, nil, nil)
	if err != nil {
		iep.Default.Message = BuildMessage(errormessage, "Error sending Box request: "+err.Error())
		showtemplate(w, tp, iep)
		return
	}
	http.Redirect(w, r, "/item", http.StatusSeeOther)
}

func patchItem(w http.ResponseWriter, r *http.Request, token string) {
	var iep ItemEditPage
	tp := "templates/item/edit.html"
	iep.Default.Sidebar = BuildSidebar(ItemsActive)
	iep.Default.Pagename = "Edit Item"
	id := r.FormValue("itemid")
	e := r.FormValue("equipment")
	b := r.FormValue("box")
	bid, err := strconv.Atoi(b)
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
	var si Item
	si.EquipmentID = eid
	si.BoxID = bid
	by := new(bytes.Buffer)
	encoder := json.NewEncoder(by)
	encoder.Encode(si)

	err = sendauthorizedHTTPRequest("PATCH", "item/"+id, token, by, nil)
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
	err := sendauthorizedHTTPRequest("GET", "item/"+id, token, nil, &ivp.Ite)
	if err != nil {
		ivp.Default.Message = BuildMessage(errormessage, "Error getting item request: "+err.Error())
		showtemplate(w, tp, ivp)
		return
	}
	err = sendauthorizedHTTPRequest("GET", "item/"+id+"/fault", token, nil, &ivp.Faults)
	if err != nil {
		ivp.Default.Message = BuildMessage(errormessage, "Error getting fault request: "+err.Error())
		showtemplate(w, tp, ivp)
		return
	}
	showtemplate(w, tp, ivp)
}

func itemLabel(w http.ResponseWriter, r *http.Request, token string) {

	b := r.FormValue("box")
	e := r.FormValue("EAN")
	w.Header().Set("Content-type", "application/pdf")
	err := createlabel(e, b, w)
	if err != nil {
		fmt.Println(err)
	}
}

func addFault(w http.ResponseWriter, r *http.Request, token string) {
	var ivp ItemViewPage
	tp := "templates/item/view.html"
	ivp.Default.Sidebar = BuildSidebar(ItemsActive)
	ivp.Default.Pagename = "View Item"
	id := r.FormValue("itemid")
	c := r.FormValue("comment")
	iid, err := strconv.Atoi(id)
	if err != nil {
		ivp.Default.Message = BuildMessage(errormessage, "Error converting Storeitem ID"+err.Error())
		showtemplate(w, tp, ivp)
		return
	}
	var f Fault
	f.Comment = c
	f.ItemID = iid
	f.Status = FaultStatusNew
	b := new(bytes.Buffer)
	encoder := json.NewEncoder(b)
	encoder.Encode(f)
	err = sendauthorizedHTTPRequest("POST", "fault/", token, b, nil)
	if err != nil {
		ivp.Default.Message = BuildMessage(errormessage, "Error posting fault request: "+err.Error())
		showtemplate(w, tp, ivp)
		return
	}
	err = sendauthorizedHTTPRequest("GET", "item/"+id, token, nil, &ivp.Ite)
	if err != nil {
		ivp.Default.Message = BuildMessage(errormessage, "Error getting item request: "+err.Error())
		showtemplate(w, tp, ivp)
		return
	}
	err = sendauthorizedHTTPRequest("GET", "item/"+id+"/fault", token, nil, &ivp.Faults)
	if err != nil {
		ivp.Default.Message = BuildMessage(errormessage, "Error getting fault request: "+err.Error())
		showtemplate(w, tp, ivp)
		return
	}
	showtemplate(w, tp, ivp)
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
	case "add-fault":
		addFault(w, r, token)
	default:
		showItemlist(w, token)
	}

}
