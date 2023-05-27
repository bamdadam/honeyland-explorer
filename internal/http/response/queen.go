package response

type QueensResponse struct {
	Type   string  `json:"type"`
	Queens []Queen `json:"queens"`
}

type Queen struct {
	Id                uint64 `json:"id"`
	Type              string `json:"type"`
	Generation        string `json:"generation"`
	Universe          string `json:"universe"`
	LandformSpecialty string `json:"landform_specialty"`
	Like              string `json:"like"`
	Dislike           string `json:"dislike"`
	Mood              string `json:"mood"`
	Health            uint16 `json:"health"`
	Attack            uint16 `json:"attack"`
	Defense           uint16 `json:"defense"`
	Agility           uint16 `json:"agility"`
	Luck              uint16 `json:"luck"`
	Capacity          uint16 `json:"capacity"`
	Recovery          uint16 `json:"recovery"`
	Endurance         uint16 `json:"endurance"`
	Level             uint8  `json:"level"`
	LevelCap          uint8  `json:"levelcap"`
	MateCount         uint16 `json:"mate_count"`
	MateCap           int8   `json:"matecap"`
	NormalAttack1     string `json:"normal_attack1"`
	NormalAttack2     string `json:"normal_attack2"`
	SpecialAttack     string `json:"special_attack"`
	DateOfBirth       string `json:"date_of_birth"`
	Mother            string `json:"mother"`
	Father            string `json:"father"`
	QueenBody         string `json:"queen_body"`
	QueenType         string `json:"queen_type"`
	Background        string `json:"background"`
	Utility           uint64 `json:"utility_rank"`
	Cosmetic          uint64 `json:"cosmetic_rank"`
	NumberOfBees      uint64 `json:"number_of_bees"`
}
