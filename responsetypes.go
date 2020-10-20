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
	EventID      int           `gorm:"primary_key;AUTO_INCREMENT;not null"`
	Name         string        `gorm:"not null"`
	Start        time.Time     `gorm:"not null"`
	End          time.Time     `gorm:"not null"`
	Adress       string        `gorm:"not null"`
	Participants []Participant `gorm:"foreignkey:EventID;association_foreignkey:EventID"`
}

type Store struct {
	StoreID   int    `gorm:"primary_key;AUTO_INCREMENT;not null"`
	Name      string `gorm:"not null"`
	Adress    string `gorm:"not null"`
	Manager   User   `gorm:"not null"`
	ManagerID int    `gorm:"foreignkey:ManagerID;not null"`
	Boxes     []Box  `gorm:"foreignkey:StoreID;association_foreignkey:StoreID"`
}

type UserRight int

const (
	USERRIGHT_MEMBER UserRight = 1 + iota
	USERRIGHT_ADMIN
)

type User struct {
	UserID   int       `json:"id" gorm:"primary_key;AUTO_INCREMENT;not null"`
	Username string    `json:"username" gorm:"not null"`
	Password string    `json:"password" gorm:"not null"`
	Salt     string    `json:"-" gorm:"not null"`
	Email    string    `json:"email" gorm:"not null"`
	Right    UserRight `json:"userright" gorm:"not null"`
}

type eventParticipiantsResponse struct {
	User      User
	Arrival   time.Time
	Departure time.Time
}

type Box struct {
	BoxID       int    `gorm:"primary_key;AUTO_INCREMENT;not null"`
	StoreID     int    `gorm:"not null"`
	Items       []Item `gorm:"foreignkey:BoxID;association_foreignkey:BoxID"`
	Code        int    `gorm:"type:integer(13)"`
	Description string `gorm:"not null"`
	Weight      int    `gorm:"not null;default:0"`
}

type Item struct {
	ItemID      int `gorm:"primary_key;AUTO_INCREMENT;not null"`
	BoxID       int
	EquipmentID int       `gorm:"not null"`
	Equipment   Equipment `gorm:"not null"`
	Code        int       `gorm:"type:integer(13)"`
	Faults      []Fault   `gorm:"foreignkey:ItemID;association_foreignkey:ItemID"`
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

func getAllFaultstates() []FaultStatus {
	var res []FaultStatus
	res = append(res, FaultStatusNew)
	res = append(res, FaultStatusInRepair)
	res = append(res, FaultStatusFixed)
	res = append(res, FaultStatusUnfixable)
	return res
}

type Fault struct {
	FaultID int         `gorm:"primary_key;AUTO_INCREMENT;not null"`
	ItemID  int         `gorm:"not null"`
	Status  FaultStatus `gorm:"not null"`
	Comment string      `gorm:"not null"`
}

type Equipment struct {
	EquipmentID int    `gorm:"primary_key;AUTO_INCREMENT;not null"`
	Name        string `gorm:"not null"`
}

type Packinglist struct {
	PackinglistID int    `gorm:"primary_key;AUTO_INCREMENT;not null"`
	Name          string `gorm:"not null"`
	EventID       int    `gorm:"foreignkey:EventID;not null"`
	Event         Event  `gorm:"not null"`
	Boxes         []Box  `gorm:"many2many:packinglist_boxes;"`
	Weight        int    `gorm:"not null;default:0"`
}

type Participant struct {
	UserID    int       `gorm:"type:integer;primary_key;not null"`
	User      User      `gorm:"not null;foreignkey:UserID;association_foreignkey:UserID"`
	EventID   int       `gorm:"type:integer;primary_key;not null"`
	Arrival   time.Time `gorm:"not null"`
	Departure time.Time `gorm:"not null"`
}

type itemResponse struct {
	Item      Item
	Store     Store
	Box       Box
	Equipment Equipment
}

type boxResponse struct {
	Box   Box
	Store Store
	User  User
}

type faultResponse struct {
	Fault Fault
	Code  int
	Name  string
}
