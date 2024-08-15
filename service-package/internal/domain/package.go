package domain

type Package struct {
	Id          string  `json:"id"`
	PackageName string  `json:"package_name"`
	Price       float32 `json:"price"`
	ForLong     string  `json:"for_long"`
}
