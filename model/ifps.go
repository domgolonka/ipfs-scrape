package model

type Ipfs struct {
	Image       string `json:"image"`
	Description string `json:"description"`
	Name        string `json:"name"`
	Cid         string `json:"cid"`
}
