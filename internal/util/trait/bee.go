package trait

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
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
	NFTNumFieldName string
}

func (m *BeeTraitsRarityMap) InitGenesis(c context.Context) error {
	brm := make(map[string]map[string]uint32)
	brm[m.NFTNumFieldName] = make(map[string]uint32)
	bq := queries.NewBeeQuery(m.DB)
	stmt, err := bq.GetTraits(true)
	if err != nil {
		logrus.Error("Can't prepare traits query: ", err.Error())
		return err
	}
	defer stmt.Close()
	bees := []model.BeeTrait{}
	err = stmt.SelectContext(c, &bees)
	if err != nil {
		logrus.Error("can't get traits to preprocess: ", err.Error())
		return err
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
	for traitName, value := range brm {
		for traitValue, num := range value {
			_, err := m.RDB.HSet(c, fmt.Sprint("Genesis-", traitName), traitValue, num)
			if err != nil {
				logrus.Error("can't set trait: ", err.Error())
				return err
			}
		}
	}
	return nil
}

func (m *BeeTraitsRarityMap) InitGeneration(c context.Context) (map[string]map[string]uint32, error) {
	brm := make(map[string]map[string]uint32)
	brm[m.NFTNumFieldName] = make(map[string]uint32)
	bq := queries.NewBeeQuery(m.DB)
	stmt, err := bq.GetTraits(false)
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
	for traitName, value := range brm {
		for traitValue, num := range value {
			_, err := m.RDB.HSet(c, fmt.Sprint("Generation-", traitName), traitValue, num)
			if err != nil {
				logrus.Error("can't set trait: ", err.Error())
				return brm, err
			}
		}
	}
	return brm, nil
}

// func (m *BeeTraitsRarityMap) Convert(ctx context.Context) (map[string]map[string]string, error) {
// 	brm := make(map[string]map[string]string)
// 	mn, err := m.RDB.HGet(ctx, m.NFTNumFieldName, "0")
// 	if err != nil {
// 		logrus.Error("error while converting: ", err)
// 	}
// 	nft_num, err := strconv.ParseUint(mn, 10, 64)
// 	if err != nil {
// 		logrus.Error("error while parsing nft number: ", err)
// 	}
// 	rv := reflect.TypeOf(model.BeeTrait{})
// 	for i := 0; i < rv.NumField(); i++ {
// 		fn := rv.Field(i).Name
// 		tt, err := m.RDB.HGetAll(ctx, fn)
// 		if err != nil {
// 			logrus.Error("error while getting trait from redis: ", err)
// 		}
// 		brm[fn] = tt

// 	}
// 	// for trait, traitVals := range m.BeeRarity {
// 	// 	brm[trait] = make(map[string]string)
// 	// 	for traitVal, value := range traitVals {
// 	// 		n := value
// 	// 		rp := (float64(n) / float64(mn)) * 100
// 	// 		brm[trait][traitVal] = strconv.FormatFloat(rp, 'f', -1, 64)
// 	// 	}
// 	// }
// 	return brm, nil
// }

func (m *BeeTraitsRarityMap) UpdateWriteMap(ctx context.Context, t model.BeeTrait) error {
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
			_, err := m.RDB.HIncryBy(ctx, fn, strconv.FormatUint(uint64(castedFv), 10), 1)
			if err != nil {
				logrus.Error("error while incr trait: ", err)
				return err
			}
		}
	}
	_, err := m.RDB.HIncryBy(ctx, m.NFTNumFieldName, "0", 1)
	if err != nil {
		logrus.Error("error while incr trait: ", err)
		return err
	}
	return nil
}

// func (m *BeeTraitsRarityMap) SetReadMap(ctx context.Context, rm map[string]map[string]string) error {
// 	m.BReadRarityMu.Lock()
// 	defer m.BReadRarityMu.Unlock()
// 	m.BeeReadRarity = rm
// 	return nil
// }

// func (m *BeeTraitsRarityMap) SetWriteMap(ctx context.Context, rm map[string]map[string]uint32) error {
// 	m.BRarityMu.Lock()
// 	defer m.BRarityMu.Unlock()
// 	m.BeeRarity = rm
// 	return nil
// }

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
	// m.BRarityMu.RLock()
	// defer m.BRarityMu.RUnlock()
	// br := m.BeeRarity
	maxVals := make(map[string]uint32)
	eyes, err := m.RDB.HGetAll(ctx, "Genesis-Eyes")
	if err != nil {
		logrus.Error("error while getting redis hash: ", err)
	}
	mouths, err := m.RDB.HGetAll(ctx, "Genesis-Mouth")
	if err != nil {
		logrus.Error("error while getting redis hash: ", err)
	}
	clothes, err := m.RDB.HGetAll(ctx, "Genesis-Clothes")
	if err != nil {
		logrus.Error("error while getting redis hash: ", err)
	}
	hats, err := m.RDB.HGetAll(ctx, "Genesis-Hat")
	if err != nil {
		logrus.Error("error while getting redis hash: ", err)
	}
	hands, err := m.RDB.HGetAll(ctx, "Genesis-BackHandAccessory")
	if err != nil {
		logrus.Error("error while getting redis hash: ", err)
	}
	backgrounds, err := m.RDB.HGetAll(ctx, "Genesis-Background")
	if err != nil {
		logrus.Error("error while getting redis hash: ", err)
	}
	maxVals["Eyes"] = m.FindRarityMaxValue(eyes)
	maxVals["Mouth"] = m.FindRarityMaxValue(mouths)
	maxVals["Clothes"] = m.FindRarityMaxValue(clothes)
	maxVals["Hat"] = m.FindRarityMaxValue(hats)
	maxVals["BackHandAccessory"] = m.FindRarityMaxValue(hands)
	maxVals["Background"] = m.FindRarityMaxValue(backgrounds)

	for _, trait := range genesisBeeTraits {
		eye, err := m.RDB.HGet(ctx, "Genesis-Eyes", strconv.FormatUint(uint64(trait.Eyes), 10))
		if err != nil {
			logrus.Error("error while getting redis hash: ", err)
		}
		pEye, err := strconv.ParseUint(eye, 10, 64)
		if err != nil {
			logrus.Error("error while parsing redis hash value: ", err)
		}
		mouth, err := m.RDB.HGet(ctx, "Genesis-Mouth", strconv.FormatUint(uint64(trait.Mouth), 10))
		if err != nil {
			logrus.Error("error while getting redis hash: ", err)
		}
		pMouth, err := strconv.ParseUint(mouth, 10, 64)
		if err != nil {
			logrus.Error("error while parsing redis hash value: ", err)
		}
		clothe, err := m.RDB.HGet(ctx, "Genesis-Clothes", strconv.FormatUint(uint64(trait.Clothes), 10))
		if err != nil {
			logrus.Error("error while getting redis hash: ", err)
		}
		pClothe, err := strconv.ParseUint(clothe, 10, 64)
		if err != nil {
			logrus.Error("error while parsing redis hash value: ", err)
		}
		hat, err := m.RDB.HGet(ctx, "Genesis-Hat", strconv.FormatUint(uint64(trait.Hat), 10))
		if err != nil {
			logrus.Error("error while getting redis hash: ", err)
		}
		pHat, err := strconv.ParseUint(hat, 10, 64)
		if err != nil {
			logrus.Error("error while parsing redis hash value: ", err)
		}
		hand, err := m.RDB.HGet(ctx, "Genesis-BackHandAccessory", strconv.FormatUint(uint64(trait.BackHandAccessory), 10))
		if err != nil {
			logrus.Error("error while getting redis hash: ", err)
		}
		pHand, err := strconv.ParseUint(hand, 10, 64)
		if err != nil {
			logrus.Error("error while parsing redis hash value: ", err)
		}
		background, err := m.RDB.HGet(ctx, "Genesis-Background", strconv.FormatUint(uint64(trait.Background), 10))
		if err != nil {
			logrus.Error("error while getting redis hash: ", err)
		}
		pBackground, err := strconv.ParseUint(background, 10, 64)
		if err != nil {
			logrus.Error("error while parsing redis hash value: ", err)
		}
		id := strconv.FormatUint(uint64(trait.NftNumber), 10)
		cosmeticScores[id] += float64(maxVals["Eyes"]) / float64(pEye)
		cosmeticScores[id] += float64(maxVals["Mouth"]) / float64(pMouth)
		cosmeticScores[id] += float64(maxVals["Clothes"]) / float64(pClothe)
		cosmeticScores[id] += float64(maxVals["Hat"]) / float64(pHat)
		cosmeticScores[id] += float64(maxVals["BackHandAccessory"]) / float64(pHand)
		cosmeticScores[id] += float64(maxVals["Background"]) / float64(pBackground)
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

func (*BeeTraitsRarityMap) FindRarityMaxValue(m map[string]string) uint32 {
	var maxValue uint32
	// Iterate over the map and update the maximum value if a greater value is found
	for _, value := range m {
		pValue, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			logrus.Error("error while parsing map max val: ", value, " err: ", err)
			return 0
		}
		if uint32(pValue) > maxValue {
			maxValue = uint32(pValue)
		}
	}
	return maxValue
}
