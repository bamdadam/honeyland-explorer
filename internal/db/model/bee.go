package model

type Bee struct {
	Id                uint64 `db:"NftId"`
	Type              string `db:"BeeType"`
	Generation        string `db:"BeeGeneration"`
	Universe          string `db:"UniverseId"`
	LandformSpecialty string `db:"Specialty1"`

	Like    string `db:"Like"`
	Dislike string `db:"Dislike"`
	Mood    string `db:"Mood"`

	Health        int16  `db:"Health"`
	Attack        int16  `db:"Attack"`
	Defense       int16  `db:"Defense"`
	Agility       int16  `db:"Agility"`
	Luck          int16  `db:"Luck"`
	Capacity      int16  `db:"Capacity"`
	Recovery      int16  `db:"Recovery"`
	Endurance     int16  `db:"Endurance"`
	Level         int8   `db:"Level"`
	LevelCap      int8   `db:"LevelCap"`
	MateCap       int8   `db:"MateCap"`
	MateCount     uint16 `db:"MateCount"`
	NormalAttack1 string `db:"AttackProfile1"`
	NormalAttack2 string `db:"AttackProfile2"`
	SpecialAttack string `db:"AttackProfile3"`
	DateOfBirth   string `db:"DateOfBirth"`

	Mother string `db:"MotherNftToken"`
	Father string `db:"FatherNftToken"`

	Head               string `db:"Head"`
	Eyes               string `db:"Eyes"`
	Mouth              string `db:"Mouth"`
	Feet               string `db:"Feet"`
	Clothes            string `db:"Clothes"`
	Hand               string `db:"Hand"`
	Hat                string `db:"Hat"`
	BackFootAccessory  string `db:"BackFootAccessory"`
	BackHandAccessory  string `db:"BackHandAccessory"`
	FrontFootAccessory string `db:"FrontFootAccessory"`
	FrontHandAccessory string `db:"FrontHandAccessory"`
	BodyVisualTrait    string `db:"BodyVisualTrait"`
	Background         string `db:"Background"`
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
