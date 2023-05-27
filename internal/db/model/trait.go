package model

type TraitsRarity struct {
	BeeReadRarity map[string]map[string]string
	BeeRarity     map[string]map[string]uint32

	QueenReadRarity map[string]map[string]string
	QueenRarity     map[string]map[string]uint32

	EggReadRarity map[string]map[string]string
	EggRarity     map[string]map[string]uint32

	LandReadRarity  map[string]map[string]string
	LandRarity      map[string]map[string]uint32
	NFTNumFieldName string
}
