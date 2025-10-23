package model

type Error struct {
	Error       string `json:"error"`
	Code        int    `json:"code"`
	Description string `json:"description"`
}
