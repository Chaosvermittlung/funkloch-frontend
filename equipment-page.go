package main

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func showEquipmentlist(w http.ResponseWriter, token string) {
	var elp EquipmentListPage
	tp := "templates/equipment/list.html"
	elp.Default.Sidebar = BuildSidebar(EquipmentActive)
	elp.Default.Pagename = "Equipment List"

	err := sendauthorizedHTTPRequest("GET", "equipment/list", token, nil, &elp.Equipments)
	if err != nil {
		elp.Default.Message = BuildMessage(errormessage, "Error creating equipment/list request: "+err.Error())
		showtemplate(w, tp, elp)
		return
	}

	showtemplate(w, tp, elp)
}

func showEquipmentAddForm(w http.ResponseWriter, token string) {
	var eap EquipmentAddPage
	eap.Default.Sidebar = BuildSidebar(EquipmentActive)
	eap.Default.Pagename = "Add Equipment"
	tp := "templates/equipment/add.html"
	showtemplate(w, tp, eap)
}

func saveNewEquipment(w http.ResponseWriter, r *http.Request, token string) {
	var elp EquipmentListPage
	tp := "templates/equipment/list.html"
	elp.Default.Sidebar = BuildSidebar(EquipmentActive)
	elp.Default.Pagename = "Equipment List"
	n := r.FormValue("equipmentname")
	var e Equipment
	e.Name = n
	b := new(bytes.Buffer)
	encoder := json.NewEncoder(b)
	encoder.Encode(e)

	err := sendauthorizedHTTPRequest("POST", "equipment/", token, b, &e)
	if err != nil {
		elp.Default.Message = BuildMessage(errormessage, "Error posting equipment request: "+err.Error())
		showtemplate(w, tp, elp)
		return
	}
	http.Redirect(w, r, "/equipment", http.StatusSeeOther)
}

func showEquipmentEditForm(w http.ResponseWriter, r *http.Request, token string) {
	var eep EquipmentEditPage
	tp := "templates/equipment/edit.html"
	eep.Default.Sidebar = BuildSidebar(EquipmentActive)
	eep.Default.Pagename = "Edit Equipment"
	id := r.FormValue("id")
	err := sendauthorizedHTTPRequest("GET", "equipment/"+id, token, nil, &eep.Equip)
	if err != nil {
		eep.Default.Message = BuildMessage(errormessage, "Error sending Equipment request: "+err.Error())
		showtemplate(w, tp, eep)
		return
	}

	showtemplate(w, tp, eep)
}

func deleteEquipment(w http.ResponseWriter, r *http.Request, token string) {
	var eep EquipmentEditPage
	tp := "templates/equipment/edit.html"
	eep.Default.Sidebar = BuildSidebar(EquipmentActive)
	eep.Default.Pagename = "Edit Equipment"
	id := r.FormValue("id")

	err := sendauthorizedHTTPRequest("DELETE", "equipment/"+id, token, nil, nil)
	if err != nil {
		eep.Default.Message = BuildMessage(errormessage, "Error sending Equipment request: "+err.Error())
		showtemplate(w, tp, eep)
		return
	}

	http.Redirect(w, r, "/equipment", http.StatusSeeOther)

}

func patchEquipment(w http.ResponseWriter, r *http.Request, token string) {
	var eep EquipmentEditPage
	tp := "templates/equipment/edit.html"
	eep.Default.Sidebar = BuildSidebar(EquipmentActive)
	eep.Default.Pagename = "Edit Equipment"
	id := r.FormValue("id")
	n := r.FormValue("equipmentname")
	var e Equipment
	e.Name = n
	b := new(bytes.Buffer)
	encoder := json.NewEncoder(b)
	encoder.Encode(e)

	err := sendauthorizedHTTPRequest("PATCH", "equipment/"+id, token, b, nil)
	if err != nil {
		eep.Default.Message = BuildMessage(errormessage, "Error sending Equipment request: "+err.Error())
		showtemplate(w, tp, eep)
		return
	}

	http.Redirect(w, r, "/equipment", http.StatusSeeOther)
}

func equipmentHandler(w http.ResponseWriter, r *http.Request) {

	token, err := GetCookie(r, "token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	a := r.FormValue("action")
	switch a {
	case "add":
		showEquipmentAddForm(w, token)
	case "save":
		saveNewEquipment(w, r, token)
	case "edit":
		showEquipmentEditForm(w, r, token)
	case "patch":
		patchEquipment(w, r, token)
	case "delete":
		deleteEquipment(w, r, token)
	default:
		showEquipmentlist(w, token)
	}

}
