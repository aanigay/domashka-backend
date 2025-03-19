package geo

type Address struct {
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	AddressType string  `json:"address_type"`
	Name        string  `json:"name"`
	Address     string  `json:"geo"`
	Comment     string  `json:"comment"`
}

type Chef struct {
	ChefID  int     `json:"chef_id"`
	Address Address `json:"geo"`
}

type Client struct {
	ClientID  int       `json:"client_id"`
	Addresses []Address `json:"addresses"`
}
