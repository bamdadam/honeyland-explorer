package model

import (
	"database/sql"
)

type Queen struct {
	Id                uint64         `db:"NftId"`
	Type              string         `db:"BeeType"`
	Generation        string         `db:"BeeGeneration"`
	Universe          string         `db:"UniverseId"`
	LandformSpecialty string         `db:"Specialty1"`
	Like              string         `db:"Like"`
	Dislike           string         `db:"Dislike"`
	Mood              string         `db:"Mood"`
	Health            uint16         `db:"Health"`
	Attack            uint16         `db:"Attack"`
	Defense           uint16         `db:"Defense"`
	Agility           uint16         `db:"Agility"`
	Luck              uint16         `db:"Luck"`
	Capacity          uint16         `db:"Capacity"`
	Recovery          uint16         `db:"Recovery"`
	Endurance         uint16         `db:"Endurance"`
	Level             uint8          `db:"Level"`
	LevelCap          uint8          `db:"LevelCap"`
	MateCap           int8           `db:"MateCap"`
	MateCount         uint16         `db:"MateCount"`
	NormalAttack1     string         `db:"AttackProfile1"`
	NormalAttack2     string         `db:"AttackProfile2"`
	SpecialAttack     string         `db:"AttackProfile3"`
	DateOfBirth       string         `db:"DateOfBirth"`
	Mother            string         `db:"MotherNftToken"`
	Father            string         `db:"FatherNftToken"`
	QueenBody         string         `db:"QueenBody"`
	QueenType         string         `db:"QueenType"`
	Background        sql.NullString `db:"Background"`
}

type QueenTrait struct {
	Generation        uint16         `db:"BeeGeneration"`
	Universe          uint16         `db:"UniverseId"`
	LandformSpecialty uint16         `db:"Specialty1"`
	Like              uint16         `db:"Like"`
	Dislike           uint16         `db:"Dislike"`
	Mood              uint16         `db:"Mood"`
	Level             uint16         `db:"Level"`
	LevelCap          uint16         `db:"LevelCap"`
	MateCap           uint16         `db:"MateCap"`
	MateCount         uint16         `db:"MateCount"`
	NormalAttack1     uint16         `db:"AttackProfile1"`
	NormalAttack2     uint16         `db:"AttackProfile2"`
	SpecialAttack     uint16         `db:"AttackProfile3"`
	QueenBody         uint16         `db:"QueenBody"`
	QueenType         uint16         `db:"QueenType"`
	Background        sql.NullString `db:"Background"`
}
