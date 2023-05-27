package request

type Queen struct {
	IsGenesis bool `json:"is_genesis"`
}

type Queens struct {
	IsGenesis bool `json:"is_genesis"`
	OrderAsc  bool `json:"order_asc"`
}
