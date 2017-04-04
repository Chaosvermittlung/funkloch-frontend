package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

func showEventlist(w http.ResponseWriter, token string) {
	var elp EventListPage
	tp := "templates/event/list.html"
	elp.Default.Sidebar = BuildSidebar(EventsActive)
	var ee []Event
	err := sendauthorizedHTTPRequest("GET", "event/list", token, nil, &ee)
	if err != nil {
		elp.Default.Message = BuildMessage(errormessage, "Error creating event/list request: "+err.Error())
		showtemplate(w, tp, elp)
		return
	}

	for _, e := range ee {
		var le ListEvent
		var eprs []eventParticipiantsResponse
		err = sendauthorizedHTTPRequest("GET", "event/"+strconv.Itoa(e.EventID)+"/Participants", token, nil, &eprs)
		if err != nil {
			elp.Default.Message = BuildMessage(errormessage, "Error creating event/Participiants request: "+err.Error())
			showtemplate(w, tp, elp)
			return
		}
		le.Participiants = len(eprs)
		le.Eve = e
		elp.Events = append(elp.Events, le)
	}

	showtemplate(w, tp, elp)
}

func viewEvent(w http.ResponseWriter, r *http.Request, token string) {
	var evp EventViewPage
	tp := "templates/event/view.html"
	evp.Default.Sidebar = BuildSidebar(EventsActive)
	id := r.FormValue("id")
	err := sendauthorizedHTTPRequest("GET", "event/"+id, token, nil, &evp.Eve)
	if err != nil {
		evp.Default.Message = BuildMessage(errormessage, "Error creating event request: "+err.Error())
		showtemplate(w, tp, evp)
		return
	}
	err = sendauthorizedHTTPRequest("GET", "event/"+id+"/Participants", token, nil, &evp.Participiants)
	if err != nil {
		evp.Default.Message = BuildMessage(errormessage, "Error creating Participiants request: "+err.Error())
		showtemplate(w, tp, evp)
		return
	}
	err = sendauthorizedHTTPRequest("GET", "user/", token, nil, &evp.Us)
	if err != nil {
		evp.Default.Message = BuildMessage(errormessage, "Error creating User request: "+err.Error())
		showtemplate(w, tp, evp)
		return
	}
	for _, p := range evp.Participiants {
		if p.User.UserID == evp.Us.UserID {
			evp.IsPart = true
			break
		}
	}
	err = sendauthorizedHTTPRequest("GET", "event/"+id+"/Packinglist", token, nil, &evp.Lists)
	if err != nil {
		evp.Default.Message = BuildMessage(errormessage, "Error creating Packinglist request: "+err.Error())
		showtemplate(w, tp, evp)
		return
	}
	showtemplate(w, tp, evp)
}

func showEventEditForm(w http.ResponseWriter, r *http.Request, token string) {
	var eep EventEditPage
	tp := "templates/event/edit.html"
	eep.Default.Sidebar = BuildSidebar(EventsActive)
	id := r.FormValue("eventid")
	err := sendauthorizedHTTPRequest("GET", "event/"+id, token, nil, &eep.Eve)
	if err != nil {
		eep.Default.Message = BuildMessage(errormessage, "Error sending Event request: "+err.Error())
		showtemplate(w, tp, eep)
		return
	}

	showtemplate(w, tp, eep)
}

func patchEvent(w http.ResponseWriter, r *http.Request, token string) {
	var eep EventEditPage
	tp := "templates/equipment/edit.html"
	eep.Default.Sidebar = BuildSidebar(EventsActive)
	id := r.FormValue("eventid")
	n := r.FormValue("eventname")
	sd := r.FormValue("startdate")
	ed := r.FormValue("enddate")
	ad := r.FormValue("adress")
	var err error
	var e Event
	e.Name = n
	e.Adress = ad
	e.Start, err = time.Parse("2006-01-02", sd)
	if err != nil {
		eep.Default.Message = BuildMessage(errormessage, "Error sending parsing Startdate: "+err.Error())
		showtemplate(w, tp, eep)
		return
	}
	e.End, err = time.Parse("2006-01-02", ed)
	if err != nil {
		eep.Default.Message = BuildMessage(errormessage, "Error sending parsing Enddate: "+err.Error())
		showtemplate(w, tp, eep)
		return
	}
	b := new(bytes.Buffer)
	encoder := json.NewEncoder(b)
	encoder.Encode(e)

	err = sendauthorizedHTTPRequest("PATCH", "event/"+id, token, b, nil)
	if err != nil {
		eep.Default.Message = BuildMessage(errormessage, "Error sending Event request: "+err.Error())
		showtemplate(w, tp, eep)
		return
	}

	http.Redirect(w, r, "/event", http.StatusSeeOther)
}

func deleteEvent(w http.ResponseWriter, r *http.Request, token string) {
	var eep EventEditPage
	tp := "templates/event/edit.html"
	eep.Default.Sidebar = BuildSidebar(EventsActive)
	id := r.FormValue("eventid")

	err := sendauthorizedHTTPRequest("DELETE", "event/"+id, token, nil, nil)
	if err != nil {
		eep.Default.Message = BuildMessage(errormessage, "Error sending Event request: "+err.Error())
		showtemplate(w, tp, eep)
		return
	}

	http.Redirect(w, r, "/event", http.StatusSeeOther)
}

func showEventAddForm(w http.ResponseWriter, token string) {
	var eap EventAddPage
	tp := "templates/event/add.html"
	eap.Default.Sidebar = BuildSidebar(EventsActive)
	showtemplate(w, tp, nil)
}

func saveNewEvent(w http.ResponseWriter, r *http.Request, token string) {
	var elp EventListPage
	tp := "templates/event/list.html"
	elp.Default.Sidebar = BuildSidebar(EventsActive)
	n := r.FormValue("eventname")
	sd := r.FormValue("startdate")
	ed := r.FormValue("enddate")
	ad := r.FormValue("adress")
	var err error
	var e Event
	e.Name = n
	e.Adress = ad
	e.Start, err = time.Parse("2006-01-02", sd)
	if err != nil {
		elp.Default.Message = BuildMessage(errormessage, "Error sending parsing Startdate: "+err.Error())
		showtemplate(w, tp, elp)
		return
	}
	e.End, err = time.Parse("2006-01-02", ed)
	if err != nil {
		elp.Default.Message = BuildMessage(errormessage, "Error sending parsing Enddate: "+err.Error())
		showtemplate(w, tp, elp)
		return
	}
	b := new(bytes.Buffer)
	encoder := json.NewEncoder(b)
	encoder.Encode(e)
	err = sendauthorizedHTTPRequest("POST", "event/", token, b, nil)
	if err != nil {
		elp.Default.Message = BuildMessage(errormessage, "Error sending Event request: "+err.Error())
		showtemplate(w, tp, elp)
		return
	}

	http.Redirect(w, r, "/event", http.StatusSeeOther)
}

func addEventParticipant(w http.ResponseWriter, r *http.Request, token string) {
	var evp EventViewPage
	tp := "templates/event/view.html"
	evp.Default.Sidebar = BuildSidebar(EventsActive)
	id := r.FormValue("eventid")
	uid := r.FormValue("userid")
	sd := r.FormValue("startdate")
	ed := r.FormValue("enddate")
	var p Participant
	var err error
	p.Arrival, err = time.Parse("2006-01-02", sd)
	if err != nil {
		evp.Default.Message = BuildMessage(errormessage, "Error sending parsing Startdate: "+err.Error())
		showtemplate(w, tp, evp)
		return
	}
	p.Departure, err = time.Parse("2006-01-02", ed)
	if err != nil {
		evp.Default.Message = BuildMessage(errormessage, "Error sending parsing Enddate: "+err.Error())
		showtemplate(w, tp, evp)
		return
	}
	p.UserID, err = strconv.Atoi(uid)
	if err != nil {
		evp.Default.Message = BuildMessage(errormessage, "Error sending converting UserID: "+err.Error())
		showtemplate(w, tp, evp)
		return
	}
	b := new(bytes.Buffer)
	encoder := json.NewEncoder(b)
	encoder.Encode(p)
	err = sendauthorizedHTTPRequest("POST", "event/"+id+"/Participants", token, b, nil)
	if err != nil {
		evp.Default.Message = BuildMessage(errormessage, "Error sending Event request: "+err.Error())
		showtemplate(w, tp, evp)
		return
	}
	http.Redirect(w, r, "/event?action=view&id="+id, http.StatusSeeOther)
}

func removeEventParticipant(w http.ResponseWriter, r *http.Request, token string) {
	var evp EventViewPage
	tp := "templates/event/view.html"
	evp.Default.Sidebar = BuildSidebar(EventsActive)
	id := r.FormValue("eventid")
	uid := r.FormValue("userid")
	var p Participant
	var err error
	p.UserID, err = strconv.Atoi(uid)
	if err != nil {
		evp.Default.Message = BuildMessage(errormessage, "Error sending converting UserID: "+err.Error())
		showtemplate(w, tp, evp)
		return
	}
	b := new(bytes.Buffer)
	encoder := json.NewEncoder(b)
	encoder.Encode(p)
	err = sendauthorizedHTTPRequest("DELETE", "event/"+id+"/Participants", token, b, nil)
	if err != nil {
		evp.Default.Message = BuildMessage(errormessage, "Error sending Event request: "+err.Error())
		showtemplate(w, tp, evp)
		return
	}
	http.Redirect(w, r, "/event?action=view&id="+id, http.StatusSeeOther)
}

func eventHandler(w http.ResponseWriter, r *http.Request) {

	token, err := GetCookie(r, "token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	a := r.FormValue("action")
	switch a {
	case "add":
		showEventAddForm(w, token)
	case "save":
		saveNewEvent(w, r, token)
	case "edit":
		showEventEditForm(w, r, token)
	case "patch":
		patchEvent(w, r, token)
	case "delete":
		deleteEvent(w, r, token)
	case "view":
		viewEvent(w, r, token)
	case "add-participant":
		addEventParticipant(w, r, token)
	case "remove-participant":
		removeEventParticipant(w, r, token)
	case "add-packinglist":

	case "view-packinglist":
	default:
		showEventlist(w, token)
	}

}
