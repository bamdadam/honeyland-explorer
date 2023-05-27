package trait

import (
	"context"
	"reflect"
	"strconv"

	"github.com/bamdadam/honeyland-explorer/internal/db/model"
	"github.com/bamdadam/honeyland-explorer/internal/queries"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type LandTraitsRarityMap struct {
	DB *sqlx.DB
}

func (m *LandTraitsRarityMap) Init(c context.Context) (map[string]map[string]uint32, error) {
	lfm := make(map[string]map[string]uint32)
	lfm["num_lands"] = make(map[string]uint32)
	lq := queries.NewLandQuery(m.DB)
	stmt, err := lq.GetTraits()
	if err != nil {
		logrus.Error("Can't prepare traits query: ", err.Error())
		return lfm, err
	}
	defer stmt.Close()
	lands := []model.LandTrait{}
	err = stmt.SelectContext(c, &lands)
	if err != nil {
		logrus.Error("can't get traits to preprocess: ", err.Error())
		return lfm, err
	}
	for _, value := range lands {
		rv := reflect.TypeOf(value)
		for i := 0; i < rv.NumField(); i++ {
			field := rv.Field(i)
			fn := field.Name
			fv := reflect.ValueOf(value).FieldByName(fn)
			fvKind := fv.Interface()
			switch v := fvKind.(type) {
			case uint16:
				castedFv := uint(v)
				_, ok := lfm[fn]
				if !ok {
					lfm[fn] = make(map[string]uint32)
				}
				lfm[fn][strconv.FormatUint(uint64(castedFv), 10)] += 1
			case float32:
				castedFv := v
				_, ok := lfm[fn]
				if !ok {
					lfm[fn] = make(map[string]uint32)
				}
				lfm[fn][strconv.FormatFloat(float64(castedFv), 'f', -1, 32)] += 1
			}
		}
		lfm["num_lands"]["0"] += 1
	}
	return lfm, nil
}
