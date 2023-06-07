package trait

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/bamdadam/honeyland-explorer/internal/db/model"
	"github.com/bamdadam/honeyland-explorer/internal/db/rdb"
	"github.com/bamdadam/honeyland-explorer/internal/queries"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type BeeTraitsRarityMap struct {
	DB              *sqlx.DB
	RDB             *rdb.RedisDB
	BeeReadRarity   map[string]map[string]string
	BReadRarityMu   sync.RWMutex
	BeeRarity       map[string]map[string]uint32
	BRarityMu       sync.RWMutex
	NFTNumFieldName string
}

func (m *BeeTraitsRarityMap) Init(c context.Context) (map[string]map[string]uint32, error) {
	brm := make(map[string]map[string]uint32)
	brm[m.NFTNumFieldName] = make(map[string]uint32)
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
		brm[m.NFTNumFieldName]["0"] += 1
	}
	return brm, nil
}

func (m *BeeTraitsRarityMap) Convert(ctx context.Context) (map[string]map[string]string, error) {
	brm := make(map[string]map[string]string)
	m.BRarityMu.RLock()
	defer m.BRarityMu.RUnlock()
	mn := m.BeeRarity[m.NFTNumFieldName]["0"]
	for trait, traitVals := range m.BeeRarity {
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
		m.BRarityMu.Lock()
		defer m.BRarityMu.Unlock()
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
				_, ok := m.BeeRarity[fn]
				if !ok {
					m.BeeRarity[fn] = make(map[string]uint32)
				}
				m.BeeRarity[fn][strconv.FormatUint(uint64(castedFv), 10)] += 1
			default:
				logrus.Error("Bee trait value is not of type uint16: ", v)
			}
		}
		m.BeeRarity[m.NFTNumFieldName]["0"] += 1
	default:
		return errors.New("invalid trait struct")
	}
	return nil
}

func (m *BeeTraitsRarityMap) SetReadMap(ctx context.Context, rm map[string]map[string]string) error {
	m.BReadRarityMu.Lock()
	defer m.BReadRarityMu.Unlock()
	m.BeeReadRarity = rm
	return nil
}

func (m *BeeTraitsRarityMap) SetWriteMap(ctx context.Context, rm map[string]map[string]uint32) error {
	m.BRarityMu.Lock()
	defer m.BRarityMu.Unlock()
	m.BeeRarity = rm
	return nil
}

func (m *BeeTraitsRarityMap) CalcGenesisTraitScore(ctx context.Context) (map[string]float64, map[string]float64, error) {
	cosmeticScores := make(map[string]float64)
	utilityScores := make(map[string]float64)
	numBees := model.BeeNum{}
	bq := queries.NewBeeQuery(m.DB)
	stmt, err := bq.GetCountBee(true)
	if err != nil {
		logrus.Error("Can't prepare count bees query: ", err.Error())
		return cosmeticScores, utilityScores, err
	}
	defer stmt.Close()
	err = stmt.GetContext(ctx, &numBees)
	if err != nil {
		logrus.Error("Can't get count bees: ", err.Error())
		return cosmeticScores, utilityScores, err
	}
	genesisBeeTraits := []model.BeeRankingTrait{}
	totalPages := (int(numBees.NumBees) / 200000) + 1
	for i := 0; i < totalPages; i++ {
		offset := i * 200000
		stmt, err := bq.GetBeeRankingTraits(true)
		if err != nil {
			logrus.Error("Can't prepare cosmetic traits query: ", err.Error())
			return cosmeticScores, utilityScores, err
		}
		defer stmt.Close()
		beeTraits := []model.BeeRankingTrait{}
		err = stmt.SelectContext(ctx, &beeTraits, sql.Named("NftId", offset))
		if err != nil {
			logrus.Error("Can't get cosmetic traits: ", err.Error())
			return cosmeticScores, utilityScores, err
		}
		genesisBeeTraits = append(genesisBeeTraits, beeTraits...)
	}
	m.BRarityMu.RLock()
	defer m.BRarityMu.RUnlock()
	br := m.BeeRarity
	maxVals := make(map[string]uint32)
	maxVals["Eyes"] = m.FindRarityMaxValue(br["Eyes"])
	maxVals["Mouth"] = m.FindRarityMaxValue(br["Mouth"])
	maxVals["Clothes"] = m.FindRarityMaxValue(br["Clothes"])
	maxVals["Hat"] = m.FindRarityMaxValue(br["Hat"])
	maxVals["BackHandAccessory"] = m.FindRarityMaxValue(br["BackHandAccessory"])
	maxVals["Background"] = m.FindRarityMaxValue(br["Background"])

	for _, trait := range genesisBeeTraits {
		id := strconv.FormatUint(uint64(trait.NftNumber), 10)
		cosmeticScores[id] += float64(maxVals["Eyes"]) / float64(br["Eyes"][strconv.FormatUint(uint64(trait.Eyes), 10)])
		cosmeticScores[id] += float64(maxVals["Mouth"]) / float64(br["Mouth"][strconv.FormatUint(uint64(trait.Mouth), 10)])
		cosmeticScores[id] += float64(maxVals["Clothes"]) / float64(br["Clothes"][strconv.FormatUint(uint64(trait.Clothes), 10)])
		cosmeticScores[id] += float64(maxVals["Hat"]) / float64(br["Hat"][strconv.FormatUint(uint64(trait.Hat), 10)])
		cosmeticScores[id] += float64(maxVals["BackHandAccessory"]) / float64(br["BackHandAccessory"][strconv.FormatUint(uint64(trait.BackHandAccessory), 10)])
		cosmeticScores[id] += float64(maxVals["Background"]) / float64(br["Background"][strconv.FormatUint(uint64(trait.Background), 10)])
		utilityScores[id] += float64(trait.Health)
		utilityScores[id] += float64(trait.Attack)
		utilityScores[id] += float64(trait.Defense)
		utilityScores[id] += float64(trait.Agility)
		utilityScores[id] += float64(trait.Luck)
		utilityScores[id] += float64(trait.Capacity)
		utilityScores[id] += float64(trait.Recovery)
		utilityScores[id] += float64(trait.Endurance)
	}
	return cosmeticScores, utilityScores, nil
}

func (m *BeeTraitsRarityMap) SetGenesisBeeSets(ctx context.Context) {
	cs, us, err := m.CalcGenesisTraitScore(ctx)
	if err != nil {
		logrus.Error("Error while calculating genesis trait score: ", err.Error())
	} else {
		genesisBeeSetC := "genesis bee cosmic"
		genesisBeeSetU := "genesis bee utility"
		for i, v := range cs {
			m.RDB.ZAdd(ctx, genesisBeeSetC, i, v)
		}
		for i, v := range us {
			m.RDB.ZAdd(ctx, genesisBeeSetU, i, v)
		}
		cScores, err := m.RDB.ZRange(ctx, genesisBeeSetC, 0, -1)
		if err != nil {
			logrus.Error("Error while getting cosmic values from redis: ", err.Error())
		}

		uScores, err := m.RDB.ZRange(ctx, genesisBeeSetU, 0, -1)
		if err != nil {
			logrus.Error("Error while getting cosmic values from redis: ", err.Error())
		}

		jsonscores, err := json.Marshal(cScores)
		if err != nil {
			logrus.Error("Can't conver scores map to json: ", err.Error())
		}
		fmt.Println(string(jsonscores))
		jsonscores, err = json.Marshal(uScores)
		if err != nil {
			logrus.Error("Can't conver scores map to json: ", err.Error())
		}
		fmt.Println(string(jsonscores))
	}
}

func (m *BeeTraitsRarityMap) SetGenesisBeeSetsScheduler(ctx context.Context, t time.Duration) {
	ticker := time.NewTicker(t)
	for range ticker.C {
		m.SetGenesisBeeSets(ctx)
	}
}

func (_ *BeeTraitsRarityMap) FindRarityMaxValue(m map[string]uint32) uint32 {
	var maxValue uint32
	// Iterate over the map and update the maximum value if a greater value is found
	for _, value := range m {
		if value > maxValue {
			maxValue = value
		}
	}
	return maxValue
}
