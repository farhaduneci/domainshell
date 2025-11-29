package domain

type DomainData struct {
	Available bool   `json:"available"`
	Domain    string `json:"domain"`
	OnSale    bool   `json:"on_sale"`
	Premium   bool   `json:"premium"`
	Prices    struct {
		Register struct {
			OneYear int `json:"1y"`
		} `json:"register"`
	} `json:"prices"`
	Reason string `json:"reason"`
}

type Response struct {
	Data []DomainData `json:"data"`
}
