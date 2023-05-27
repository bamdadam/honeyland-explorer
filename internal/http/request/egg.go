package request

type Egg struct {
	IsGenesis bool `json:"is_genesis"`
}

type Eggs struct {
	IsGenesis bool `json:"is_genesis"`
	OrderAsc  bool `json:"order_asc"`
}
