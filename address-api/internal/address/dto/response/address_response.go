package response

type AddressResponse struct {
	Id          int    `json:"id"`
	City        string `json:"city"`
	Country     string `json:"country"`
	FullAddress string `json:"fullAddress"`
	UserId      string `json:"userId"`
}
