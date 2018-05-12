package main

import "net/http"

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
		//showItemEditForm(w, r, token)
	case "patch":
		//patchItem(w, r, token)
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
