package trait

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
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
	pre := "Genesis-Bee-"
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
			_, err := m.RDB.HSet(c, pre+traitName, traitValue, num)
			if err != nil {
				logrus.Error("can't set trait: ", err.Error())
				return err
			}
		}
	}
	return nil
}

func (m *BeeTraitsRarityMap) InitGeneration(c context.Context) error {
	brm := make(map[string]map[string]uint32)
	brm[m.NFTNumFieldName] = make(map[string]uint32)
	pre := "Generation-Bee-"
	bq := queries.NewBeeQuery(m.DB)
	stmt, err := bq.GetTraits(false)
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
			_, err := m.RDB.HSet(c, pre+traitName, traitValue, num)
			if err != nil {
				logrus.Error("can't set trait: ", err.Error())
				return err
			}
		}
	}
	return nil
}

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

func (m *BeeTraitsRarityMap) CalcGenesisTraitScore(ctx context.Context) (map[string]float64, map[string]float64, error) {
	cosmeticScores := make(map[string]float64)
	utilityScores := make(map[string]float64)
	pre := "Genesis-Bee-"
	bq := queries.NewBeeQuery(m.DB)
	numBees, err := m.RDB.HGetUint(ctx, pre+m.NFTNumFieldName, "0")
	if err != nil {
		logrus.Error("error while getting number of bees: ", err)
		return cosmeticScores, utilityScores, err
	}
	genesisBeeTraits := []model.BeeRankingTrait{}
	totalPages := (int(numBees) / 200000) + 1
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
	maxVals := make(map[string]uint32)
	eyes, err := m.RDB.HGetAll(ctx, pre+"Eyes")
	if err != nil {
		logrus.Errorf("error while getting %v redis hash: %v", pre+"Eyes", err)
	}
	mouths, err := m.RDB.HGetAll(ctx, pre+"Mouth")
	if err != nil {
		logrus.Errorf("error while getting %v redis hash: %v", pre+"Mouth", err)
	}
	clothes, err := m.RDB.HGetAll(ctx, pre+"Clothes")
	if err != nil {
		logrus.Errorf("error while getting %v redis hash: %v", pre+"Clothes", err)
	}
	hats, err := m.RDB.HGetAll(ctx, pre+"Hat")
	if err != nil {
		logrus.Errorf("error while getting %v redis hash: %v", pre+"Hat", err)
	}
	hands, err := m.RDB.HGetAll(ctx, pre+"BackHandAccessory")
	if err != nil {
		logrus.Errorf("error while getting %v redis hash: %v", pre+"BackHandAccessory", err)
	}
	backgrounds, err := m.RDB.HGetAll(ctx, pre+"Background")
	if err != nil {
		logrus.Errorf("error while getting %v redis hash: %v", pre+"Background", err)
	}
	maxVals["Eyes"] = m.FindRarityMaxValue(eyes)
	maxVals["Mouth"] = m.FindRarityMaxValue(mouths)
	maxVals["Clothes"] = m.FindRarityMaxValue(clothes)
	maxVals["Hat"] = m.FindRarityMaxValue(hats)
	maxVals["BackHandAccessory"] = m.FindRarityMaxValue(hands)
	maxVals["Background"] = m.FindRarityMaxValue(backgrounds)

	for _, trait := range genesisBeeTraits {
		eye, err := m.RDB.HGetUint(ctx, pre+"Eyes", strconv.FormatUint(uint64(trait.Eyes), 10))
		if err != nil {
			logrus.Errorf("error while getting %v redis hash: %v", pre+"Eyes", err)
		}
		mouth, err := m.RDB.HGetUint(ctx, pre+"Mouth", strconv.FormatUint(uint64(trait.Mouth), 10))
		if err != nil {
			logrus.Errorf("error while getting %v redis hash: %v", pre+"Mouth", err)
		}
		clothe, err := m.RDB.HGetUint(ctx, pre+"Clothes", strconv.FormatUint(uint64(trait.Clothes), 10))
		if err != nil {
			logrus.Errorf("error while getting %v redis hash: %v", pre+"Clothes", err)
		}
		hat, err := m.RDB.HGetUint(ctx, pre+"Hat", strconv.FormatUint(uint64(trait.Hat), 10))
		if err != nil {
			logrus.Errorf("error while getting %v redis hash: %v", pre+"Hat", err)
		}
		hand, err := m.RDB.HGetUint(ctx, pre+"BackHandAccessory", strconv.FormatUint(uint64(trait.BackHandAccessory), 10))
		if err != nil {
			logrus.Errorf("error while getting %v redis hash: %v", pre+"BackHandAccessory", err)
		}
		background, err := m.RDB.HGetUint(ctx, pre+"Background", strconv.FormatUint(uint64(trait.Background), 10))
		if err != nil {
			logrus.Errorf("error while getting %v redis hash: %v", pre+"Background", err)
		}
		id := strconv.FormatUint(uint64(trait.NftNumber), 10)
		cosmeticScores[id] += float64(maxVals["Eyes"]) / float64(eye)
		cosmeticScores[id] += float64(maxVals["Mouth"]) / float64(mouth)
		cosmeticScores[id] += float64(maxVals["Clothes"]) / float64(clothe)
		cosmeticScores[id] += float64(maxVals["Hat"]) / float64(hat)
		cosmeticScores[id] += float64(maxVals["BackHandAccessory"]) / float64(hand)
		cosmeticScores[id] += float64(maxVals["Background"]) / float64(background)
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

func (m *BeeTraitsRarityMap) CalcGenerationTraitScore(ctx context.Context) (map[string]float64, map[string]float64, error) {
	cosmeticScores := make(map[string]float64)
	utilityScores := make(map[string]float64)
	pre := "Generation-Bee-"
	bq := queries.NewBeeQuery(m.DB)
	numBees, err := m.RDB.HGetUint(ctx, pre+m.NFTNumFieldName, "0")
	if err != nil {
		logrus.Error("error while getting number of bees: ", err)
		return cosmeticScores, utilityScores, err
	}
	generationBeeTraits := []model.BeeRankingTrait{}
	totalPages := (int(numBees) / 2000) + 1
	fmt.Println(totalPages)
	fmt.Println(numBees)
	for i := 0; i < totalPages; i++ {
		offset := i * 2000
		stmt, err := bq.GetBeeRankingTraits(false)
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
		generationBeeTraits = append(generationBeeTraits, beeTraits...)
	}
	fmt.Println(len(generationBeeTraits))
	tt, err := json.Marshal(generationBeeTraits)
	if err != nil {
		logrus.Error("Can't marshal traits: ", err.Error())
	}
	err = os.WriteFile("/home/bamdad/bee-ranking-trait.json", tt, 0644)
	if err != nil {
		logrus.Error("can't write to file: ", err)
	}
	maxVals := make(map[string]uint32)
	eyes, err := m.RDB.HGetAll(ctx, pre+"Eyes")
	if err != nil {
		logrus.Errorf("error while getting %v redis hash: %v", pre+"Eyes", err)
	}
	mouths, err := m.RDB.HGetAll(ctx, pre+"Mouth")
	if err != nil {
		logrus.Errorf("error while getting %v redis hash: %v", pre+"Mouth", err)
	}
	clothes, err := m.RDB.HGetAll(ctx, pre+"Clothes")
	if err != nil {
		logrus.Errorf("error while getting %v redis hash: %v", pre+"Clothes", err)
	}
	hats, err := m.RDB.HGetAll(ctx, pre+"Hat")
	if err != nil {
		logrus.Errorf("error while getting %v redis hash: %v", pre+"Hat", err)
	}
	hands, err := m.RDB.HGetAll(ctx, pre+"BackHandAccessory")
	if err != nil {
		logrus.Errorf("error while getting %v redis hash: %v", pre+"BackHandAccessory", err)
	}
	backgrounds, err := m.RDB.HGetAll(ctx, pre+"Background")
	if err != nil {
		logrus.Errorf("error while getting %v redis hash: %v", pre+"Background", err)
	}
	maxVals["Eyes"] = m.FindRarityMaxValue(eyes)
	maxVals["Mouth"] = m.FindRarityMaxValue(mouths)
	maxVals["Clothes"] = m.FindRarityMaxValue(clothes)
	maxVals["Hat"] = m.FindRarityMaxValue(hats)
	maxVals["BackHandAccessory"] = m.FindRarityMaxValue(hands)
	maxVals["Background"] = m.FindRarityMaxValue(backgrounds)

	for _, trait := range generationBeeTraits {
		eye, err := m.RDB.HGetUint(ctx, pre+"Eyes", strconv.FormatUint(uint64(trait.Eyes), 10))
		if err != nil {
			logrus.Errorf("error while getting %v redis hash: %v", pre+"Eyes", err)
		}
		mouth, err := m.RDB.HGetUint(ctx, pre+"Mouth", strconv.FormatUint(uint64(trait.Mouth), 10))
		if err != nil {
			logrus.Errorf("error while getting %v redis hash: %v", pre+"Mouth", err)
		}
		clothe, err := m.RDB.HGetUint(ctx, pre+"Clothes", strconv.FormatUint(uint64(trait.Clothes), 10))
		if err != nil {
			logrus.Errorf("error while getting %v redis hash: %v", pre+"Clothes", err)
		}
		hat, err := m.RDB.HGetUint(ctx, pre+"Hat", strconv.FormatUint(uint64(trait.Hat), 10))
		if err != nil {
			logrus.Errorf("error while getting %v redis hash: %v", pre+"Hat", err)
		}
		hand, err := m.RDB.HGetUint(ctx, pre+"BackHandAccessory", strconv.FormatUint(uint64(trait.BackHandAccessory), 10))
		if err != nil {
			logrus.Errorf("error while getting %v redis hash: %v", pre+"BackHandAccessory", err)
		}
		background, err := m.RDB.HGetUint(ctx, pre+"Background", strconv.FormatUint(uint64(trait.Background), 10))
		if err != nil {
			logrus.Errorf("error while getting %v redis hash %v : %v", pre+"Background", trait.Background, err)
		}
		id := strconv.FormatUint(uint64(trait.NftNumber), 10)
		cosmeticScores[id] += float64(maxVals["Eyes"]) / float64(eye)
		cosmeticScores[id] += float64(maxVals["Mouth"]) / float64(mouth)
		cosmeticScores[id] += float64(maxVals["Clothes"]) / float64(clothe)
		cosmeticScores[id] += float64(maxVals["Hat"]) / float64(hat)
		cosmeticScores[id] += float64(maxVals["BackHandAccessory"]) / float64(hand)
		cosmeticScores[id] += float64(maxVals["Background"]) / float64(background)
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

func (m *BeeTraitsRarityMap) setGenesisBeeSets(ctx context.Context) {
	cs, us, err := m.CalcGenesisTraitScore(ctx)
	pre := "Genesis-Bee-"
	if err != nil {
		logrus.Error("Error while calculating genesis trait score: ", err.Error())
	} else {
		beeSetC := pre + "BeeCosmetic"
		beeSetU := pre + "BeeUtility"
		for i, v := range cs {
			m.RDB.ZAdd(ctx, beeSetC, i, v)
		}
		for i, v := range us {
			m.RDB.ZAdd(ctx, beeSetU, i, v)
		}
		// cScores, err := m.RDB.ZRange(ctx, beeSetC, 0, -1)
		// if err != nil {
		// 	logrus.Error("Error while getting cosmetic values from redis: ", err.Error())
		// }

		// uScores, err := m.RDB.ZRange(ctx, beeSetU, 0, -1)
		// if err != nil {
		// 	logrus.Error("Error while getting cosmetic values from redis: ", err.Error())
		// }

		// jsonscores, err := json.Marshal(cScores)
		// if err != nil {
		// 	logrus.Error("Can't conver scores map to json: ", err.Error())
		// }
		// fmt.Println(string(jsonscores))
		// jsonscores, err = json.Marshal(uScores)
		// if err != nil {
		// 	logrus.Error("Can't conver scores map to json: ", err.Error())
		// }
		// fmt.Println(string(jsonscores))
	}
}

func (m *BeeTraitsRarityMap) setGenerationBeeSets(ctx context.Context) {
	cs, us, err := m.CalcGenerationTraitScore(ctx)
	pre := "Generation-Bee-"
	if err != nil {
		logrus.Error("Error while calculating generation trait score: ", err.Error())
	} else {
		beeSetC := pre + "BeeCosmetic"
		beeSetU := pre + "BeeUtility"
		for i, v := range cs {
			m.RDB.ZAdd(ctx, beeSetC, i, v)
		}
		for i, v := range us {
			m.RDB.ZAdd(ctx, beeSetU, i, v)
		}
		// cScores, err := m.RDB.ZRange(ctx, beeSetC, 0, -1)
		// if err != nil {
		// 	logrus.Error("Error while getting cosmetic values from redis: ", err.Error())
		// }

		// uScores, err := m.RDB.ZRange(ctx, beeSetU, 0, -1)
		// if err != nil {
		// 	logrus.Error("Error while getting cosmetic values from redis: ", err.Error())
		// }

		// jsonscores, err := json.Marshal(cScores)
		// if err != nil {
		// 	logrus.Error("Can't conver scores map to json: ", err.Error())
		// }
		// fmt.Println(string(jsonscores))
		// jsonscores, err = json.Marshal(uScores)
		// if err != nil {
		// 	logrus.Error("Can't conver scores map to json: ", err.Error())
		// }
		// fmt.Println(string(jsonscores))
	}
}

func (m *BeeTraitsRarityMap) SetBeeSetsScheduler(ctx context.Context, t time.Duration) {

	m.setGenesisBeeSets(ctx)
	m.setGenerationBeeSets(ctx)

	ticker := time.NewTicker(t)
	for range ticker.C {
		m.setGenesisBeeSets(ctx)
		m.setGenerationBeeSets(ctx)
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
