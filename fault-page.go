package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
)

func showFaultlist(w http.ResponseWriter, token string) {
	var flp FaultListPage
	tp := "templates/fault/list.html"
	flp.Default.Sidebar = BuildSidebar(FaultsActive)
	flp.Default.Pagename = "Fault List"

	err := sendauthorizedHTTPRequest("GET", "fault/list", token, nil, &flp.Faults)
	if err != nil {
		flp.Default.Message = BuildMessage(errormessage, "Error creating fault/list request: "+err.Error())
		showtemplate(w, tp, flp)
		return
	}
	showtemplate(w, tp, flp)
}

func showFaultEditForm(w http.ResponseWriter, r *http.Request, token string) {
	var fep FaultEditPage
	tp := "templates/fault/edit.html"
	fep.Default.Sidebar = BuildSidebar(FaultsActive)
	fep.Default.Pagename = "Edit Fault"
	id := r.FormValue("id")
	err := sendauthorizedHTTPRequest("GET", "fault/"+id, token, nil, &fep.Fault.Fault)
	if err != nil {
		fep.Default.Message = BuildMessage(errormessage, "Error sending Fault request: "+err.Error())
		showtemplate(w, tp, fep)
		return
	}
	sid := strconv.Itoa(fep.Fault.Fault.StoreItemID)
	var sir storeItemResponse
	err = sendauthorizedHTTPRequest("GET", "storeitem/"+sid, token, nil, &sir)
	if err != nil {
		fep.Default.Message = BuildMessage(errormessage, "Error sending Storitem request: "+err.Error())
		showtemplate(w, tp, fep)
		return
	}
	fep.Fault.Code = sir.StoreItem.Code
	fep.Fault.Name = sir.Equipment.Name
	fep.States = getAllFaultstates()
	showtemplate(w, tp, fep)
}

func patchFault(w http.ResponseWriter, r *http.Request, token string) {
	var fep FaultEditPage
	tp := "templates/fault/edit.html"
	fep.Default.Sidebar = BuildSidebar(FaultsActive)
	fep.Default.Pagename = "Fault Item"
	id := r.FormValue("faultid")
	c := r.FormValue("comment")
	s := r.FormValue("state")
	sid, err := strconv.Atoi(s)
	if err != nil {
		fep.Default.Message = BuildMessage(errormessage, "Error converting Faultstate"+err.Error())
		showtemplate(w, tp, fep)
		return
	}
	var f Fault
	f.Comment = c
	f.Status = FaultStatus(sid)
	b := new(bytes.Buffer)
	encoder := json.NewEncoder(b)
	encoder.Encode(f)

	err = sendauthorizedHTTPRequest("PATCH", "fault/"+id, token, b, nil)
	if err != nil {
		fep.Default.Message = BuildMessage(errormessage, "Error sending Fault request: "+err.Error())
		showtemplate(w, tp, fep)
		return
	}
	http.Redirect(w, r, "/fault", http.StatusSeeOther)
}

func faultHandler(w http.ResponseWriter, r *http.Request) {

	token, err := GetCookie(r, "token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	a := r.FormValue("action")
	switch a {
	case "add":
		//showItemAddForm(w, token)
	case "save":
		//saveNewItem(w, r, token)
	case "edit":
		showFaultEditForm(w, r, token)
	case "patch":
		patchFault(w, r, token)
	case "delete":
		//deleteItem(w, r, token)
	case "view":
		//viewItem(w, r, token)
	case "label":
		//itemLabel(w, r, token)
	case "add-fault":
		//addFault(w, r, token)
	default:
		showFaultlist(w, token)
	}

}