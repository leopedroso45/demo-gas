package model

type GasStation struct {
	Name           string
	Address        string
	Phone          string
	OwnerOrManager Owner
	Services       []Service
	Prices         []Price
}

type Owner struct {
	Name  string
	Email string
	Phone string
}

type Service struct {
	Id          string
	Description string
}

type Price struct {
	ProductName string
	Value       float64
}
