package models

type Vote struct {
	ID     int
	UserID int
	Action string
	ItemID int
	Item   string
}
