package main

import (
	"net/http"
)

func showBoxlist(w http.ResponseWriter, token string) {
	var blp BoxListPage
	tp := "templates/box/list.html"
	blp.Default.Sidebar = BuildSidebar(BoxesActive)
	blp.Default.Pagename = "Item List"

	err := sendauthorizedHTTPRequest("GET", "box/list", token, nil, &blp.Boxes)
	if err != nil {
		blp.Default.Message = BuildMessage(errormessage, "Error creating item/list request: "+err.Error())
		showtemplate(w, tp, blp)
		return
	}
	showtemplate(w, tp, blp)
}

func showBoxAddForm(w http.ResponseWriter, token string) {
	var bap BoxAddPage
	tp := "templates/box/add.html"
	bap.Default.Sidebar = BuildSidebar(BoxesActive)
	bap.Default.Pagename = "Add Box"
	err := sendauthorizedHTTPRequest("GET", "store/list", token, nil, &bap.Stores)
	if err != nil {
		bap.Default.Message = BuildMessage(errormessage, "Error sending Box request: "+err.Error())
		showtemplate(w, tp, bap)
		return
	}
	showtemplate(w, tp, bap)
}

func boxHandler(w http.ResponseWriter, r *http.Request) {

	token, err := GetCookie(r, "token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	a := r.FormValue("action")
	switch a {
	case "add":
		showBoxAddForm(w, token)
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
		showBoxlist(w, token)
	}

}
