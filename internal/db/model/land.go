package model

type Land struct {
	Id                   string  `db:"Number"`
	Level                uint16  `db:"Level"`
	Universe             string  `db:"UniverseId"`
	LandForm             string  `db:"Landform"`
	Climate              string  `db:"IsExtreme"`
	Feature              string  `db:"Feature"`
	BackGround           uint16  `db:"BackgroundId"`
	HoneyProductionSpeed uint64  `db:"ProductionSpeed"`
	HoneypotDropRate     uint16  `db:"DropTimeInMinute"`
	Capacity             uint16  `db:"Capacity"`
	Zone                 uint8   `db:"Zone"`
	MaxCommissionFee     float32 `db:"MaxCommissionFee"`
	MaxEntryFee          float32 `db:"MaxEntryFee"`
}

type LandTrait struct {
	Level                uint16 `db:"Level"`
	Universe             uint16 `db:"UniverseId"`
	LandForm             uint16 `db:"Landform"`
	Climate              bool   `db:"IsExtreme"`
	Feature              string `db:"Feature"`
	BackGround           uint16 `db:"BackgroundId"`
	HoneyProductionSpeed uint16 `db:"ProductionSpeed"`
	HoneypotDropRate     uint16 `db:"DropTimeInMinute"`
	Capacity             uint16 `db:"Capacity"`
	Zone                 uint16 `db:"Zone"`
}
