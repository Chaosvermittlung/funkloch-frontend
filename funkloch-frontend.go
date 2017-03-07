package main

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Loginpage struct {
	Message template.HTML
}

func loginprocess(w http.ResponseWriter, r *http.Request) error {
	u := User{}
	u.Username = r.FormValue("username")
	err := u.Load()
	if err != nil {
		fmt.Println(err)
		return errors.New("Username or Password wrong")
	}
	h := hashPassword(r.FormValue("password"))
	if bytes.Compare(u.Password, h) != 0 {
		return errors.New("Username or Password wrong")
	}
	err = SetCookie(w, u.Username)
	if err != nil {
		return errors.New("WAWAWAWAWA Cookies not allowed")
	}
	http.Redirect(w, r, "/", http.StatusFound)
	return nil
}

func loginhandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var loginerror string
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

	lp := Loginpage{template.HTML(loginerror)}

	err = t.Execute(w, &lp)
	if err != nil {
		log.Println(err)
	}
}

func mainhandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/main.html")
	t.Execute(w, nil)

}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", mainhandler)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	http.ListenAndServe(":1234", r)

}
