package main

import "time"

type ErrorResponse struct {
	Httpstatus   string `json:"httpstatus"`
	Errorcode    string `json:"errorcode"`
	Errormessage string `json:"errormessage"`
}

type authResponse struct {
	Token string `json:"token"`
}

type Event struct {
	EventID int
	Name    string
	Start   time.Time
	End     time.Time
	Adress  string
}

type Store struct {
	StoreID int
	Name    string
	Adress  string
	Manager int
}

type UserRight int

const (
	USERRIGHT_MEMBER UserRight = 1 + iota
	USERRIGHT_ADMIN
)

type User struct {
	UserID   int       `json:"id"`
	Username string    `json:"username"`
	Password string    `json:"password"`
	Salt     string    `json:"-"`
	Email    string    `json:"email"`
	Right    UserRight `json:"userright"`
}

type eventParticipiantsResponse struct {
	User      User
	Arrival   time.Time
	Departure time.Time
}

type StoreItem struct {
	StoreItemID int
	StoreID     int
	EquipmentID int
	Code        int
}

type FaultStatus int

const (
	FaultStatusNew FaultStatus = 0 + iota
	FaultStatusInRepair
	FaultStatusFixed
	FaultStatusUnfixable
)

func (f FaultStatus) String() string {
	switch f {
	case FaultStatusNew:
		return "New"
	case FaultStatusInRepair:
		return "In Repair"
	case FaultStatusFixed:
		return "Fixed"
	case FaultStatusUnfixable:
		return "Unfixable"
	default:
		return "Unkown Faultstatus"
	}
}

type Fault struct {
	FaultID     int
	StoreItemID int
	Status      FaultStatus
	Comment     string
}

type Equipment struct {
	EquipmentID int
	Name        string
}

type Packinglist struct {
	PackinglistID int
	Name          string
	EventID       int
}

type Participant struct {
	UserID    int
	EventID   int
	Arrival   time.Time
	Departure time.Time
}

type storeItemResponse struct {
	StoreItem StoreItem
	Store     Store
	Equipment Equipment
}

type faultResponse struct {
	Fault Fault
	Code  int
	Name  string
}
