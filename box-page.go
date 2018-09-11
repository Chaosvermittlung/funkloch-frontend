package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func showBoxlist(w http.ResponseWriter, token string) {
	var blp BoxListPage
	tp := "templates/box/list.html"
	blp.Default.Sidebar = BuildSidebar(BoxesActive)
	blp.Default.Pagename = "Box List"

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

func saveNewBox(w http.ResponseWriter, r *http.Request, token string) {
	var bap BoxAddPage
	tp := "templates/obx/add.html"
	var bnp BoxNewPage
	tp2 := "templates/box/new.html"
	bap.Default.Sidebar = BuildSidebar(BoxesActive)
	bap.Default.Pagename = "Add Box"
	d := r.FormValue("description")
	s := r.FormValue("store")
	sid, err := strconv.Atoi(s)
	if err != nil {
		bap.Default.Message = BuildMessage(errormessage, "Error converting Store ID"+err.Error())
		showtemplate(w, tp, bap)
		return
	}

	var bo Box
	bo.Description = d
	bo.StoreID = sid

	b := new(bytes.Buffer)
	encoder := json.NewEncoder(b)
	encoder.Encode(bo)

	bnp.Default.Pagename = "New Box"
	bnp.Default.Sidebar = BuildSidebar(BoxesActive)
	err = sendauthorizedHTTPRequest("POST", "box/", token, b, &bo)
	if err != nil {
		bnp.Default.Message = BuildMessage(errormessage, "Error sending Box request: "+err.Error())
		showtemplate(w, tp2, bnp)
		return
	}

	var bx boxResponse
	err = sendauthorizedHTTPRequest("GET", "box/"+strconv.Itoa(bo.BoxID), token, nil, &bx)
	if err != nil {
		bnp.Default.Message = BuildMessage(errormessage, "Error recieving Box request: "+err.Error())
		showtemplate(w, tp2, bnp)
		return
	}

	bnp.ID = strconv.Itoa(bx.Box.Code)
	bnp.Description = bx.Box.Description
	bnp.Store = bx.Store.Name

	showtemplate(w, tp2, bnp)
}

func boxLabel(w http.ResponseWriter, r *http.Request, token string) {

	s := r.FormValue("store")
	e := r.FormValue("EAN")
	w.Header().Set("Content-type", "application/pdf")
	err := createlabel(e, s, w)
	if err != nil {
		fmt.Println(err)
	}
}

func contentLabel(w http.ResponseWriter, r *http.Request, token string) {
	id := r.FormValue("boxid")
	var irs []itemResponse
	err := sendauthorizedHTTPRequest("GET", "box/"+id+"/items", token, nil, &irs)
	if err != nil {
		fmt.Fprintln(w, "Error getting items request: "+err.Error())
		fmt.Println("Error getting items request: " + err.Error())
		return
	}
	w.Header().Set("Content-type", "application/pdf")
	err = createContentlabel(irs, w)
	if err != nil {
		fmt.Println(err)
	}
}

func viewBox(w http.ResponseWriter, r *http.Request, token string) {
	var bvp BoxViewPage
	tp := "templates/box/view.html"
	bvp.Default.Sidebar = BuildSidebar(BoxesActive)
	bvp.Default.Pagename = "View Box"
	id := r.FormValue("id")
	err := sendauthorizedHTTPRequest("GET", "box/"+id, token, nil, &bvp.Box)
	if err != nil {
		bvp.Default.Message = BuildMessage(errormessage, "Error getting box request: "+err.Error())
		showtemplate(w, tp, bvp)
		return
	}
	err = sendauthorizedHTTPRequest("GET", "box/"+id+"/items", token, nil, &bvp.Items)
	if err != nil {
		bvp.Default.Message = BuildMessage(errormessage, "Error getting items request: "+err.Error())
		showtemplate(w, tp, bvp)
		return
	}
	err = sendauthorizedHTTPRequest("GET", "item/storeless", token, nil, &bvp.Storeless)
	if err != nil {
		bvp.Default.Message = BuildMessage(errormessage, "Error getting storeless items request: "+err.Error())
		showtemplate(w, tp, bvp)
		return
	}
	showtemplate(w, tp, bvp)
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
		saveNewBox(w, r, token)
	case "edit":
		showItemEditForm(w, r, token)
	case "patch":
		patchItem(w, r, token)
	case "delete":
		deleteItem(w, r, token)
	case "view":
		viewBox(w, r, token)
	case "label":
		boxLabel(w, r, token)
	case "content-label":
		contentLabel(w, r, token)
	case "add-fault":
		addFault(w, r, token)
	default:
		showBoxlist(w, token)
	}

}
