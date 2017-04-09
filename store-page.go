package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
)

func showStorelist(w http.ResponseWriter, token string) {
	var slp StoreListPage
	tp := "templates/store/list.html"
	slp.Default.Sidebar = BuildSidebar(StoresActive)
	slp.Default.Pagename = "Store List"

	err := sendauthorizedHTTPRequest("GET", "store/list", token, nil, &slp.Stores)
	if err != nil {
		slp.Default.Message = BuildMessage(errormessage, "Error creating store/list request: "+err.Error())
		showtemplate(w, tp, slp)
		return
	}
	showtemplate(w, tp, slp)
}

func showStoreAddForm(w http.ResponseWriter, token string) {
	var sap StoreAddPage
	tp := "templates/store/add.html"
	sap.Default.Sidebar = BuildSidebar(StoresActive)
	sap.Default.Pagename = "Add Store"
	err := sendauthorizedHTTPRequest("GET", "user/list", token, nil, &sap.Users)
	if err != nil {
		sap.Default.Message = BuildMessage(errormessage, "Error sending Users request: "+err.Error())
		showtemplate(w, tp, sap)
		return
	}
	showtemplate(w, tp, sap)
}

func saveNewStore(w http.ResponseWriter, r *http.Request, token string) {
	var sap EventEditPage
	tp := "templates/store/add.html"
	sap.Default.Sidebar = BuildSidebar(StoresActive)
	sap.Default.Pagename = "Add Store"
	n := r.FormValue("storename")
	ad := r.FormValue("adress")
	mn := r.FormValue("manager")
	mnid, err := strconv.Atoi(mn)
	var s Store
	s.Name = n
	s.Adress = ad
	if err != nil {
		sap.Default.Message = BuildMessage(errormessage, "Error converting Manager ID"+err.Error())
		showtemplate(w, tp, sap)
		return
	}
	s.Manager = mnid
	b := new(bytes.Buffer)
	encoder := json.NewEncoder(b)
	encoder.Encode(s)

	err = sendauthorizedHTTPRequest("POST", "store/", token, b, nil)
	if err != nil {
		sap.Default.Message = BuildMessage(errormessage, "Error sending Store request: "+err.Error())
		showtemplate(w, tp, sap)
		return
	}
	http.Redirect(w, r, "/store", http.StatusSeeOther)
}

func showStoreEditForm(w http.ResponseWriter, r *http.Request, token string) {
	var sep StoreEditPage
	tp := "templates/store/edit.html"
	sep.Default.Sidebar = BuildSidebar(StoresActive)
	sep.Default.Pagename = "Edit Store"
	id := r.FormValue("storeid")
	err := sendauthorizedHTTPRequest("GET", "store/"+id, token, nil, &sep.Sto)
	if err != nil {
		sep.Default.Message = BuildMessage(errormessage, "Error sending Event request: "+err.Error())
		showtemplate(w, tp, sep)
		return
	}
	err = sendauthorizedHTTPRequest("GET", "user/list", token, nil, &sep.Users)
	if err != nil {
		sep.Default.Message = BuildMessage(errormessage, "Error sending Users request: "+err.Error())
		showtemplate(w, tp, sep)
		return
	}
	showtemplate(w, tp, sep)
}

func deleteStore(w http.ResponseWriter, r *http.Request, token string) {
	var sep EventEditPage
	tp := "templates/store/edit.html"
	sep.Default.Sidebar = BuildSidebar(StoresActive)
	sep.Default.Pagename = "Edit Store"
	id := r.FormValue("storeid")

	err := sendauthorizedHTTPRequest("DELETE", "store/"+id, token, nil, nil)
	if err != nil {
		sep.Default.Message = BuildMessage(errormessage, "Error sending Store request: "+err.Error())
		showtemplate(w, tp, sep)
		return
	}

	http.Redirect(w, r, "/store", http.StatusSeeOther)
}

func patchStore(w http.ResponseWriter, r *http.Request, token string) {
	var sep EventEditPage
	tp := "templates/store/edit.html"
	sep.Default.Sidebar = BuildSidebar(StoresActive)
	sep.Default.Pagename = "Edit Store"
	id := r.FormValue("storeid")
	n := r.FormValue("storename")
	ad := r.FormValue("adress")
	mn := r.FormValue("manager")
	mnid, err := strconv.Atoi(mn)
	var s Store
	s.Name = n
	s.Adress = ad
	if err != nil {
		sep.Default.Message = BuildMessage(errormessage, "Error converting Manager ID"+err.Error())
		showtemplate(w, tp, sep)
		return
	}
	s.Manager = mnid
	b := new(bytes.Buffer)
	encoder := json.NewEncoder(b)
	encoder.Encode(s)

	err = sendauthorizedHTTPRequest("PATCH", "store/"+id, token, b, nil)
	if err != nil {
		sep.Default.Message = BuildMessage(errormessage, "Error sending Store request: "+err.Error())
		showtemplate(w, tp, sep)
		return
	}
	http.Redirect(w, r, "/store", http.StatusSeeOther)
}

func viewStore(w http.ResponseWriter, r *http.Request, token string) {
	var svp StoreViewPage
	tp := "templates/store/view.html"
	svp.Default.Sidebar = BuildSidebar(StoresActive)
	svp.Default.Pagename = "View Store"
	id := r.FormValue("id")
	err := sendauthorizedHTTPRequest("GET", "store/"+id, token, nil, &svp.Sto)
	if err != nil {
		svp.Default.Message = BuildMessage(errormessage, "Error creating event request: "+err.Error())
		showtemplate(w, tp, svp)
		return
	}
	err = sendauthorizedHTTPRequest("GET", "store/"+id+"/Manager", token, nil, &svp.Manager)
	if err != nil {
		svp.Default.Message = BuildMessage(errormessage, "Error creating Manager request: "+err.Error())
		showtemplate(w, tp, svp)
		return
	}

	showtemplate(w, tp, svp)
}

func storeHandler(w http.ResponseWriter, r *http.Request) {

	token, err := GetCookie(r, "token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	a := r.FormValue("action")
	switch a {
	case "add":
		showStoreAddForm(w, token)
	case "save":
		saveNewStore(w, r, token)
	case "edit":
		showStoreEditForm(w, r, token)
	case "patch":
		patchStore(w, r, token)
	case "delete":
		deleteStore(w, r, token)
	case "view":
		viewStore(w, r, token)
	default:
		showStorelist(w, token)
	}

}
