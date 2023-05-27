package model

type Egg struct {
	Id         uint64 `db:"NftId"`
	Type       string `db:"Type"`
	Generation string `db:"Generation"`
	Universe   string `db:"UniverseId"`
	Like       string `db:"Like"`
	Dislike    string `db:"Dislike"`
}

type EggTrait struct {
	Generation uint16 `db:"Generation"`
	Like       uint16 `db:"Like"`
	Dislike    uint16 `db:"Dislike"`
}
