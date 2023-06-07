package response

type Bee struct {
	Id                 uint64   `json:"id"`
	Type               uint16   `json:"type"`
	Generation         uint16   `json:"generation"`
	Universe           uint16   `json:"universe"`
	LandformSpecialty  uint16   `json:"landformSpecialty"`
	Like               uint16   `json:"like"`
	Dislike            uint16   `json:"dislike"`
	Mood               uint16   `json:"mood"`
	HxdPerMinute       float64  `json:"hxdPerMinute"`
	HxdPerTwoHour      float64  `json:"hxdPerTwoHour"`
	HxdCapacity        float64  `json:"hxdCapacity"`
	RecoveryTime       int16    `json:"recoveryTime"`
	Health             uint16   `json:"health"`
	Attack             uint16   `json:"attack"`
	Defense            uint16   `json:"defense"`
	Agility            uint16   `json:"agility"`
	Luck               uint16   `json:"luck"`
	Capacity           uint16   `json:"capacity"`
	Recovery           uint16   `json:"recovery"`
	Endurance          uint16   `json:"endurance"`
	Level              uint16   `json:"level"`
	LevelCap           uint16   `json:"levelCap"`
	MateCap            uint16   `json:"mateCap"`
	NormalAttack1      uint16   `json:"normalAttack1"`
	NormalAttack2      uint16   `json:"normalAttack2"`
	SpecialAttack      uint16   `json:"specialAttack"`
	DateOfBirth        string   `json:"dateOfBirth"`
	Mother             string   `json:"mother"`
	Father             string   `json:"father"`
	Head               uint16   `json:"head"`
	Eyes               uint16   `json:"eyes"`
	Mouth              uint16   `json:"mouth"`
	Clothes            uint16   `json:"clothes"`
	BackFootAccessory  uint16   `json:"backFootAccessory"`
	BackHandAccessory  uint16   `json:"backHandAccessory"`
	FrontFootAccessory uint16   `json:"frontFootAccessory"`
	FrontHandAccessory uint16   `json:"frontHandAccessory"`
	BodyVisualTrait    uint16   `json:"bodyVisualTrait"`
	Background         uint16   `json:"background"`
	Utility            int16    `json:"utilityRank"`
	Cosmetic           int16    `json:"cosmeticRank"`
	Traits             BeeTrait `json:"traits"`
}

type BeeTrait struct {
	Generation         string `json:"beeGeneration"`
	Universe           string `json:"universeId"`
	LandformSpecialty  string `json:"specialty1"`
	Like               string `json:"like"`
	Dislike            string `json:"dislike"`
	Mood               string `json:"mood"`
	Level              string `json:"level"`
	LevelCap           string `json:"levelCap"`
	MateCap            string `json:"mateCap"`
	MateCount          string `json:"mateCount"`
	NormalAttack1      string `json:"attackProfile1"`
	NormalAttack2      string `json:"attackProfile2"`
	SpecialAttack      string `json:"attackProfile3"`
	Head               string `json:"head"`
	Eyes               string `json:"eyes"`
	Mouth              string `json:"mouth"`
	Feet               string `json:"feet"`
	Clothes            string `json:"clothes"`
	Hand               string `json:"hand"`
	Hat                string `json:"hat"`
	BackFootAccessory  string `json:"backFootAccessory"`
	BackHandAccessory  string `json:"backHandAccessory"`
	FrontFootAccessory string `json:"frontFootAccessory"`
	FrontHandAccessory string `json:"frontHandAccessory"`
	BodyVisualTrait    string `json:"bodyVisualTrait"`
	Background         string `json:"background"`
}
