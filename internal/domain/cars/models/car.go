package models

type Car struct {
	CID       string `json:"cid"`
	Lable     string `json:"lable"`
	Model     string `json:"model"`
	Year      int    `json:"year"`
	Price     int    `json:"price"`
	Available bool   `json:"available"`
	Count     int    `json:"count"`
}
