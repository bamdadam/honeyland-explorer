package request

type Bee struct {
	IsGenesis bool `json:"is_genesis"`
}

type Bees struct {
	IsGenesis bool `json:"is_genesis"`
	OrderAsc  bool `json:"order_asc"`
}
