package response

type LandsResponse struct {
	Type  string `json:"type"`
	Lands []Land `json:"lands"`
}

type Land struct {
	Name                 string  `json:"name"`
	Level                uint16  `json:"level"`
	Universe             string  `json:"universe"`
	LandForm             string  `json:"landform"`
	Climate              string  `json:"climate"`
	Feature              string  `json:"feature"`
	HXDPerDay            float64 `json:"hxd_per_day"`
	HoneypotPerDay       float64 `json:"honeypot_per_day"`
	EnduranceToReach     uint32  `json:"endurance_to_reach"`
	HoneyProductionSpeed uint64  `json:"honey_production_speed"`
	HoneyProductionGrade string  `json:"honey_production_grade"`
	HoneypotDropRate     uint16  `json:"honeypot_drop_rate"`
	HoneypotDropGrade    string  `json:"honeypot_drop_grade"`
	Zone                 uint8   `json:"zone"`
	MaxCommission        float32 `json:"max_commission"`
	Utility              uint64  `json:"utility_rank"`
	Cosmetic             uint64  `json:"cosmetic_rank"`
	NumberOfLands        uint64  `json:"number_of_lands"`
}
