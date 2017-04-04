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

type DefaultPage struct {
	Message  template.HTML
	Sidebar  template.HTML
	Username string
}

type Mainpage struct {
	Default DefaultPage
	Event   FrontEvent
	Item    FrontItem
	Member  FrontMember
	Fault   FrontFault
}

type EquipmentListPage struct {
	Default    DefaultPage
	Equipments []Equipment
}

type EquipmentAddPage struct {
	Default DefaultPage
}

type EquipmentEditPage struct {
	Default DefaultPage
	Equip   Equipment
}

type ListEvent struct {
	Eve           Event
	Participiants int
}

type EventListPage struct {
	Default DefaultPage
	Events  []ListEvent
}

type EventViewPage struct {
	Default       DefaultPage
	Eve           Event
	Us            User
	Participiants []eventParticipiantsResponse
	IsPart        bool
	Lists         []Packinglist
}

type EventEditPage struct {
	Default DefaultPage
	Eve     Event
}
