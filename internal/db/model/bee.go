package model

type Bee struct {
	Id                 uint64 `db:"NftId"`
	Type               uint16 `db:"BeeType"`
	Generation         uint16 `db:"BeeGeneration"`
	Universe           uint16 `db:"UniverseId"`
	LandformSpecialty  uint16 `db:"Specialty1"`
	Like               uint16 `db:"Like"`
	Dislike            uint16 `db:"Dislike"`
	Mood               uint16 `db:"Mood"`
	Health             uint16 `db:"Health"`
	Attack             uint16 `db:"Attack"`
	Defense            uint16 `db:"Defense"`
	Agility            uint16 `db:"Agility"`
	Luck               uint16 `db:"Luck"`
	Capacity           uint16 `db:"Capacity"`
	Recovery           uint16 `db:"Recovery"`
	Endurance          uint16 `db:"Endurance"`
	Level              uint16 `db:"Level"`
	LevelCap           uint16 `db:"LevelCap"`
	MateCap            uint16 `db:"MateCap"`
	MateCount          uint16 `db:"MateCount"`
	NormalAttack1      uint16 `db:"AttackProfile1"`
	NormalAttack2      uint16 `db:"AttackProfile2"`
	SpecialAttack      uint16 `db:"AttackProfile3"`
	DateOfBirth        string `db:"DateOfBirth"`
	Mother             string `db:"MotherNftToken"`
	Father             string `db:"FatherNftToken"`
	Head               uint16 `db:"Head"`
	Eyes               uint16 `db:"Eyes"`
	Mouth              uint16 `db:"Mouth"`
	Feet               uint16 `db:"Feet"`
	Clothes            uint16 `db:"Clothes"`
	Hand               uint16 `db:"Hand"`
	Hat                uint16 `db:"Hat"`
	BackFootAccessory  uint16 `db:"BackFootAccessory"`
	BackHandAccessory  uint16 `db:"BackHandAccessory"`
	FrontFootAccessory uint16 `db:"FrontFootAccessory"`
	FrontHandAccessory uint16 `db:"FrontHandAccessory"`
	BodyVisualTrait    uint16 `db:"BodyVisualTrait"`
	Background         uint16 `db:"Background"`
}

type BeeTrait struct {
	Generation         uint16 `db:"beeGeneration"`
	Universe           uint16 `db:"UniverseId"`
	LandformSpecialty  uint16 `db:"Specialty1"`
	Like               uint16 `db:"Like"`
	Dislike            uint16 `db:"Dislike"`
	Mood               uint16 `db:"Mood"`
	Level              uint16 `db:"Level"`
	LevelCap           uint16 `db:"LevelCap"`
	MateCap            uint16 `db:"MateCap"`
	MateCount          uint16 `db:"MateCount"`
	NormalAttack1      uint16 `db:"AttackProfile1"`
	NormalAttack2      uint16 `db:"AttackProfile2"`
	SpecialAttack      uint16 `db:"AttackProfile3"`
	Head               uint16 `db:"Head"`
	Eyes               uint16 `db:"Eyes"`
	Mouth              uint16 `db:"Mouth"`
	Feet               uint16 `db:"Feet"`
	Clothes            uint16 `db:"Clothes"`
	Hand               uint16 `db:"Hand"`
	Hat                uint16 `db:"Hat"`
	BackFootAccessory  uint16 `db:"BackFootAccessory"`
	BackHandAccessory  uint16 `db:"BackHandAccessory"`
	FrontFootAccessory uint16 `db:"FrontFootAccessory"`
	FrontHandAccessory uint16 `db:"FrontHandAccessory"`
	BodyVisualTrait    uint16 `db:"BodyVisualTrait"`
	Background         uint16 `db:"Background"`
}

type BeeNum struct {
	NumBees uint32 `db:"numBees"`
}

type BeeRankingTrait struct {
	Eyes              uint16 `db:"Eyes"`
	Mouth             uint16 `db:"Mouth"`
	Clothes           uint16 `db:"Clothes"`
	Hat               uint16 `db:"Hat"`
	BackHandAccessory uint16 `db:"BackHandAccessory"`
	Background        uint16 `db:"Background"`
	NftNumber         uint32 `db:"NftNumber"`
	Health            uint16 `db:"Health"`
	Attack            uint16 `db:"Attack"`
	Defense           uint16 `db:"Defense"`
	Agility           uint16 `db:"Agility"`
	Luck              uint16 `db:"Luck"`
	Capacity          uint16 `db:"Capacity"`
	Recovery          uint16 `db:"Recovery"`
	Endurance         uint16 `db:"Endurance"`
}

type UpdateBeeGrpc struct {
	StatDifference int32
	Generation     uint16
	Id             uint32
}
