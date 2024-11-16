package event

type AddressEvent struct {
	EventType   string `json:"event_type"`
	AddressId   int    `json:"addressId"`
	City        string `json:"city"`
	Country     string `json:"country"`
	FullAddress string `json:"fullAddress"`
	UserId      string `json:"userId"`
}
