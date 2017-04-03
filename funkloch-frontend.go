package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var conf config

func showtemplate(w http.ResponseWriter, path string, data interface{}) {

	t, err := template.ParseFiles(path)
	if err != nil {
		fmt.Fprintln(w, "Error parsing template:", err)
		return
	}
	err = t.Execute(w, data)
	if err != nil {
		fmt.Fprintln(w, "Error executing template:", err)
		return
	}
}

func loginprocess(w http.ResponseWriter, r *http.Request) error {
	un := r.FormValue("username")
	pw := r.FormValue("password")

	req, err := getNewHTTPRequest("GET", "auth", nil)
	if err != nil {
		return errors.New("Error creating request: " + err.Error())
	}

	req.SetBasicAuth(un, pw)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.New("Error executing request: " + err.Error())
	}

	if resp.StatusCode > 299 {
		fmt.Println(resp.StatusCode)
		decoder := json.NewDecoder(resp.Body)
		var er ErrorResponse
		err = decoder.Decode(&er)
		if err != nil {
			data, _ := ioutil.ReadAll(resp.Body)
			fmt.Println(string(data))
			return errors.New("Error while decoding Error response. Only God can help you now:" + err.Error())
		}
		return errors.New("Got negativ status code")

	}

	decoder := json.NewDecoder(resp.Body)
	var ar authResponse
	err = decoder.Decode(&ar)
	if err != nil {
		return errors.New("Error while decoding Auth response: " + err.Error())
	}
	err = SetCookie(w, "token", ar.Token)
	if err != nil {
		return errors.New("WAWAWAWAWA Cookies not allowed")
	}
	http.Redirect(w, r, "/", http.StatusFound)
	return nil
}

func loginhandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var loginerror template.HTML
	a := r.FormValue("action")
	if a == "do" {
		err = loginprocess(w, r)
		if err == nil {
			return
		}
		loginerror = BuildMessage(errormessage, err.Error())
	}
	t, err := template.ParseFiles("templates/login.html")

	if err != nil {
		log.Fatal(err)
	}

	lp := Loginpage{loginerror}

	err = t.Execute(w, &lp)
	if err != nil {
		log.Println(err)
	}
}

func mainhandler(w http.ResponseWriter, r *http.Request) {
	var mp Mainpage
	tp := "templates/main.html"

	token, err := GetCookie(r, "token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}
	var ee []Event

	err = sendauthorizedHTTPRequest("GET", "event/list", token, nil, &ee)
	if err != nil {
		mp.Message = BuildMessage(errormessage, "Error creating event/list request: "+err.Error())
		showtemplate(w, tp, mp)
		return
	}

	mp.Event.Count = len(ee)

	var nextEvent Event
	err = sendauthorizedHTTPRequest("GET", "event/next", token, nil, &nextEvent)
	if err != nil {
		mp.Message = BuildMessage(errormessage, "Error creating event/next request: "+err.Error())
		showtemplate(w, tp, mp)
		return
	}

	dur := nextEvent.Start.Sub(time.Now())
	mp.Event.Days = int(Round(dur.Hours()/24, 0.5, 0))

	var eprs []eventParticipiantsResponse
	err = sendauthorizedHTTPRequest("GET", "event/"+strconv.Itoa(nextEvent.EventID)+"/Participants", token, nil, &eprs)
	if err != nil {
		mp.Message = BuildMessage(errormessage, "Error creating event/Participants request: "+err.Error())
		showtemplate(w, tp, mp)
		return
	}
	mp.Event.Registred = len(eprs)

	var ss []StoreItem
	err = sendauthorizedHTTPRequest("GET", "storeitem/list", token, nil, &ss)
	if err != nil {
		mp.Message = BuildMessage(errormessage, "Error creating storeitem request: "+err.Error())
		showtemplate(w, tp, mp)
		return
	}
	mp.Item.Count = len(ss)

	var sts []Store
	err = sendauthorizedHTTPRequest("GET", "store/list", token, nil, &sts)
	if err != nil {
		mp.Message = BuildMessage(errormessage, "Error creating store request: "+err.Error())
		showtemplate(w, tp, mp)
		return
	}
	mp.Item.Stores = len(sts)

	var uu []User
	err = sendauthorizedHTTPRequest("GET", "user/list", token, nil, &uu)
	if err != nil {
		mp.Message = BuildMessage(errormessage, "Error creating store request: "+err.Error())
		showtemplate(w, tp, mp)
		return
	}
	mp.Member.Registred = len(uu)

	var ff []Fault
	err = sendauthorizedHTTPRequest("GET", "fault/list", token, nil, &ff)
	if err != nil {
		mp.Message = BuildMessage(errormessage, "Error creating store request: "+err.Error())
		showtemplate(w, tp, mp)
		return
	}
	count := 0
	for _, f := range ff {
		if (f.Status != FaultStatusUnfixable) && (f.Status != FaultStatusFixed) {
			count++
		}
	}
	mp.Fault.Count = count

	showtemplate(w, tp, mp)

}

func main() {
	err := conf.Load()
	if err != nil {
		log.Fatal(err)
	}
	r := mux.NewRouter()
	r.HandleFunc("/", mainhandler)
	r.HandleFunc("/equipment", equipmentHandler)
	r.HandleFunc("/event", eventHandler)
	r.HandleFunc("/login", loginhandler)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(conf.Port), r))

}
