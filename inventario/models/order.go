package models

type Product struct {
	Title       string  `json:"title"`
	Author      string  `json:"author"`
	Genre       string  `json:"genre"`
	Pages       int     `json:"pages"`
	Publication string  `json:"publication"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
}

type Location struct {
	Address1   string `json:"address1"`
	Address2   string `json:"address2"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postalCode"`
	Country    string `json:"country"`
}

type Customer struct {
	Name     string   `json:"name"`
	Lastname string   `json:"lastname"`
	Email    string   `json:"email"`
	Location Location `json:"location"`
	Phone    string   `json:"phone"`
}

type OrderMessage struct {
	OrderID  string    `json:"orderID"`
	Products []Product `json:"products"`
	Customer Customer  `json:"customer"`
}
