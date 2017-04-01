package main

import "html/template"

type Loginpage struct {
	Message template.HTML
}

type FrontEvent struct {
	Count     int
	Days      int
	Registred int
}

type FrontItem struct {
	Count  int
	Stores int
}

type FrontMember struct {
	Registred int
}

type FrontFault struct {
	Count int
}

type Mainpage struct {
	Message template.HTML
	Event   FrontEvent
	Item    FrontItem
	Member  FrontMember
	Fault   FrontFault
}
