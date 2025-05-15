package geo

type Address struct {
	ID             int64   `json:"id"`
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	AddressType    string  `json:"address_type"`
	Name           string  `json:"name"`
	Address        *string `json:"address"`
	Comment        *string `json:"comment"`
	FlatNumber     *string `json:"flat_number,omitempty"`
	FloorNumber    *string `json:"floor_number,omitempty"`
	EntranceNumber *string `json:"entrance_number,omitempty"`
	IntercomNumber *string `json:"intercom_number,omitempty"`
	CourierComment *string `json:"courier_comment,omitempty"`
}

type Chef struct {
	ChefID  int     `json:"chef_id"`
	Address Address `json:"geo"`
}

type Client struct {
	ClientID  int       `json:"client_id"`
	Addresses []Address `json:"addresses"`
}
