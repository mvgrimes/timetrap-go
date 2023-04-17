package models

type Meta struct {
	CurrentSheet string `json:"current_sheet"`
	LastSheet    string `json:"last_sheet"`
	LastCheckout int    `json:"last_checkout"`
}
