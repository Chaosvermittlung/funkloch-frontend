package main

import (
	"net/http"
)

func showPackinglistList(w http.ResponseWriter, token string) {

	var plp PackinglistListPage
	tp := "templates/packinglist/list.html"
	plp.Default.Sidebar = BuildSidebar(PackinglistActive)
	plp.Default.Pagename = "Packinglist List"

	err := sendauthorizedHTTPRequest("GET", "packinglist/list", token, nil, &plp.Packinglists)
	if err != nil {
		plp.Default.Message = BuildMessage(errormessage, "Error sending packinglist/list request: "+err.Error())
		showtemplate(w, tp, plp)
		return
	}
	showtemplate(w, tp, plp)

}

func showPackinglistAddForm(w http.ResponseWriter, token string) {
}

func saveNewPackinglist(w http.ResponseWriter, r *http.Request, token string) {
}

func showPackinglistEditForm(w http.ResponseWriter, r *http.Request, token string) {
}

func patchPackinglist(w http.ResponseWriter, r *http.Request, token string) {
}

func deletePackinglist(w http.ResponseWriter, r *http.Request, token string) {
}

func viewPackinglist(w http.ResponseWriter, r *http.Request, token string) {
}

func addPackinglistItem(w http.ResponseWriter, r *http.Request, token string) {
}

func removePackinglistItem(w http.ResponseWriter, r *http.Request, token string) {
}

func packinglistHandler(w http.ResponseWriter, r *http.Request) {

	token, err := GetCookie(r, "token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	a := r.FormValue("action")
	switch a {
	case "add":
		showPackinglistAddForm(w, token)
	case "save":
		saveNewPackinglist(w, r, token)
	case "edit":
		showPackinglistEditForm(w, r, token)
	case "patch":
		patchPackinglist(w, r, token)
	case "delete":
		deletePackinglist(w, r, token)
	case "view":
		viewPackinglist(w, r, token)
	case "add-item":
		addPackinglistItem(w, r, token)
	case "remove-item":
		removePackinglistItem(w, r, token)
	default:
		showPackinglistList(w, token)
	}

}
