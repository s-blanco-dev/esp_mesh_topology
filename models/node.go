package models

type Node struct {
	ParentMAC string  `json:"parentMAC"`
	SelfMAC   string  `json:"selfMAC"`
	Temp      float32 `json:"temp"`
	Humidity  float32 `json:"humidity"`
	IsRoot    bool    `json:"isRoot"`
}
