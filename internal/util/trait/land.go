package trait

import (
	"context"
	"errors"
	"reflect"
	"strconv"
	"sync"

	"github.com/bamdadam/honeyland-explorer/internal/db/model"
	"github.com/bamdadam/honeyland-explorer/internal/queries"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type LandTraitsRarityMap struct {
	DB              *sqlx.DB
	LandReadRarity  map[string]map[string]string
	LReadRarityMu   sync.RWMutex
	LandRarity      map[string]map[string]uint32
	LRarityMu       sync.RWMutex
	NFTNumFieldName string
}

func (l *LandTraitsRarityMap) Init(c context.Context) (map[string]map[string]uint32, error) {
	lfm := make(map[string]map[string]uint32)
	lfm[l.NFTNumFieldName] = make(map[string]uint32)
	lq := queries.NewLandQuery(l.DB)
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
		lfm[l.NFTNumFieldName]["0"] += 1
	}
	return lfm, nil
}

func (l *LandTraitsRarityMap) Convert(ctx context.Context) (map[string]map[string]string, error) {
	lrm := make(map[string]map[string]string)
	l.LRarityMu.RLock()
	defer l.LRarityMu.RUnlock()
	mn := l.LandRarity[l.NFTNumFieldName]["0"]
	for trait, traitVals := range l.LandRarity {
		lrm[trait] = make(map[string]string)
		for traitVal, value := range traitVals {
			n := value
			rp := (float64(n) / float64(mn)) * 100
			lrm[trait][traitVal] = strconv.FormatFloat(rp, 'f', -1, 64)
		}
	}
	return lrm, nil
}

func (l *LandTraitsRarityMap) UpdateWriteMap(ctx context.Context, t interface{}) error {
	switch t.(type) {
	case model.LandTrait:
		l.LRarityMu.Lock()
		defer l.LRarityMu.Unlock()
		land := t
		rv := reflect.TypeOf(land)
		for i := 0; i < rv.NumField(); i++ {
			field := rv.Field(i)
			fn := field.Name
			fv := reflect.ValueOf(land).FieldByName(fn)
			fvKind := fv.Interface()
			switch v := fvKind.(type) {
			case uint16:
				castedFv := uint(v)
				_, ok := l.LandRarity[fn]
				if !ok {
					l.LandRarity[fn] = make(map[string]uint32)
				}
				l.LandRarity[fn][strconv.FormatUint(uint64(castedFv), 10)] += 1
			case bool:
				castedFv := v
				_, ok := l.LandRarity[fn]
				if !ok {
					l.LandRarity[fn] = make(map[string]uint32)
				}
				l.LandRarity[fn][strconv.FormatBool(castedFv)] += 1
			case string:
				castedFv := v
				_, ok := l.LandRarity[fn]
				if !ok {
					l.LandRarity[fn] = make(map[string]uint32)
				}
				l.LandRarity[fn][castedFv] += 1
			default:
				logrus.Error("Bee trait value is not of accpeted types: ", v)
			}
		}
		l.LandRarity[l.NFTNumFieldName]["0"] += 1
	default:
		return errors.New("invalid trait struct")
	}
	return nil
}

func (l *LandTraitsRarityMap) SetReadMap(ctx context.Context, rm map[string]map[string]string) error {
	l.LReadRarityMu.Lock()
	defer l.LReadRarityMu.Unlock()
	l.LandReadRarity = rm
	return nil
}

func (l *LandTraitsRarityMap) SetWriteMap(ctx context.Context, rm map[string]map[string]uint32) error {
	l.LRarityMu.Lock()
	defer l.LRarityMu.Unlock()
	l.LandRarity = rm
	return nil
}

func (mp *LandTraitsRarityMap) FindRarityMaxValue(m map[string]uint32) uint32 {
	var maxValue uint32
	// Iterate over the map and update the maximum value if a greater value is found
	for _, value := range m {
		if value > maxValue {
			maxValue = value
		}
	}
	return maxValue
}
