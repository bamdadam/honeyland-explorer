package trait

import (
	"context"
	"errors"
	"reflect"
	"strconv"

	"github.com/bamdadam/honeyland-explorer/internal/db/model"
	"github.com/bamdadam/honeyland-explorer/internal/queries"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type BeeTraitsRarityMap struct {
	DB *sqlx.DB
	TR *model.TraitsRarity
}

func (m *BeeTraitsRarityMap) Init(c context.Context) (map[string]map[string]uint32, error) {
	brm := make(map[string]map[string]uint32)
	brm["num_bees"] = make(map[string]uint32)
	bq := queries.NewBeeQuery(m.DB)
	stmt, err := bq.GetTraits()
	if err != nil {
		logrus.Error("Can't prepare traits query: ", err.Error())
		return brm, err
	}
	defer stmt.Close()
	bees := []model.BeeTrait{}
	err = stmt.SelectContext(c, &bees)
	if err != nil {
		logrus.Error("can't get traits to preprocess: ", err.Error())
		return brm, err
	}
	for _, value := range bees {
		rv := reflect.TypeOf(value)
		for i := 0; i < rv.NumField(); i++ {
			field := rv.Field(i)
			fn := field.Name
			fv := reflect.ValueOf(value).FieldByName(fn)
			fvKind := fv.Interface()
			switch v := fvKind.(type) {
			case uint16:
				castedFv := uint(v)
				_, ok := brm[fn]
				if !ok {
					brm[fn] = make(map[string]uint32)
				}
				brm[fn][strconv.FormatUint(uint64(castedFv), 10)] += 1
			}
		}
		brm["num_bees"]["0"] += 1
	}
	return brm, nil
}

// func (m *BeeTraitsRarityMap) Convert(ctx context.Context, rm map[string]map[string]uint32) (map[string]map[string]string, error) {
// 	brm := make(map[string]map[string]string)

// 	mn := rm[m.TR.NFTNumFieldName]["0"]
// 	for trait, traitVals := range rm {
// 		brm[trait] = make(map[string]string)
// 		for traitVal, value := range traitVals {
// 			n := value
// 			rp := (float64(n) / float64(mn)) * 100
// 			brm[trait][traitVal] = strconv.FormatFloat(rp, 'f', -1, 64)
// 		}
// 	}
// 	return brm, nil
// }

func (m *BeeTraitsRarityMap) Convert(ctx context.Context) (map[string]map[string]string, error) {
	brm := make(map[string]map[string]string)

	mn := m.TR.BeeRarity[m.TR.NFTNumFieldName]["0"]
	for trait, traitVals := range m.TR.BeeRarity {
		brm[trait] = make(map[string]string)
		for traitVal, value := range traitVals {
			n := value
			rp := (float64(n) / float64(mn)) * 100
			brm[trait][traitVal] = strconv.FormatFloat(rp, 'f', -1, 64)
		}
	}
	return brm, nil
}

func (m *BeeTraitsRarityMap) UpdateWriteMap(ctx context.Context, t interface{}) error {
	switch t.(type) {
	case model.BeeTrait:
		bee := t
		rv := reflect.TypeOf(bee)
		for i := 0; i < rv.NumField(); i++ {
			field := rv.Field(i)
			fn := field.Name
			fv := reflect.ValueOf(bee).FieldByName(fn)
			fvKind := fv.Interface()
			switch v := fvKind.(type) {
			case uint16:
				castedFv := uint(v)
				_, ok := m.TR.BeeRarity[fn]
				if !ok {
					m.TR.BeeRarity[fn] = make(map[string]uint32)
				}
				m.TR.BeeRarity[fn][strconv.FormatUint(uint64(castedFv), 10)] += 1
			default:
				logrus.Error("Bee trait value is not of type uint16: ", v)
			}
		}
		m.TR.BeeRarity["num_bees"]["0"] += 1
	default:
		return errors.New("invalid trait struct")
	}
	return nil
}

func (m *BeeTraitsRarityMap) SetReadMap(ctx context.Context, rm map[string]map[string]string) error {
	m.TR.BeeReadRarity = rm
	return nil
}

func (m *BeeTraitsRarityMap) SetWriteMap(ctx context.Context, rm map[string]map[string]uint32) error {
	m.TR.BeeRarity = rm
	return nil
}
