package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/bamdadam/honeyland-explorer/internal/db/model"
	"github.com/bamdadam/honeyland-explorer/internal/http/request"
	"github.com/bamdadam/honeyland-explorer/internal/http/response"
	"github.com/bamdadam/honeyland-explorer/internal/queries"
	"github.com/bamdadam/honeyland-explorer/internal/util/trait"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type Bee struct {
	DB  *sqlx.DB
	Btm *trait.BeeTraitsRarityMap
}

func (b *Bee) GetBeeById(c *fiber.Ctx) error {
	// get id
	id := c.Params("id")
	body := new(request.Bee)
	if err := c.BodyParser(body); err != nil {
		logrus.Error("cant parse request body: ", err.Error())
		return fiber.ErrBadRequest
	}
	logrus.Debug("body is: ", body.IsGenesis)
	// cast id to uint64
	uid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		logrus.Error("can't parse id to uint64: ", err.Error())
		return fiber.ErrBadRequest
	}
	// load bee
	bq := queries.NewBeeQuery(b.DB)
	stmt, err := bq.LoadById(body.IsGenesis)
	if err != nil {
		logrus.Error("error while preparing query: ", err.Error())
		return fiber.ErrInternalServerError
	}
	defer stmt.Close()
	bee := model.Bee{}
	err = stmt.GetContext(c.Context(), &bee, sql.Named("NftId", uid))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(http.StatusOK).JSON(make([]string, 0))
		}
		logrus.Error("error while getting bee: ", err.Error())
		return fiber.ErrInternalServerError
	}
	fmt.Println("bee is: ", bee)
	// load response
	beeResponse, err := b.FillResponseWithTrait(bee)
	if err != nil {
		logrus.Error("error while filling response")
		return fiber.ErrInternalServerError
	}
	cs, us, err := b.CalculateRanks(c, id, body.IsGenesis)
	if err != nil {
		logrus.Error("error while calculating ranks")
		return fiber.ErrInternalServerError
	}
	beeResponse.Cosmetic = cs
	beeResponse.Utility = us
	beeResponse.HxdPerMinute = b.CalculateHxdPerTimeFrame(bee.Agility, time.Minute)
	beeResponse.HxdPerTwoHour = b.CalculateHxdPerTimeFrame(bee.Agility, 2*time.Hour)
	beeResponse.HxdCapacity = b.CalculateHxdCapacity(bee.Capacity)
	beeResponse.RecoveryTime = b.CalculateRecoveryTime(bee.Recovery)
	return c.Status(http.StatusOK).JSON(beeResponse)
}

func (b *Bee) GetBeesById(c *fiber.Ctx) error {
	// get id
	id := c.Params("id")
	body := new(request.Bees)
	if err := c.BodyParser(body); err != nil {
		logrus.Error("cant parse request body: ", err.Error())
		return fiber.ErrBadRequest
	}
	logrus.Debug("is genesis: ", body.IsGenesis)
	logrus.Debug("order is: ", body.OrderAsc)
	// cast id to uint64
	uid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		logrus.Error("can't parse id to uint64: ", err.Error())
		return fiber.ErrBadRequest
	}
	// load bee
	bq := queries.NewBeeQuery(b.DB)
	stmt, err := bq.LoadByPage(body.IsGenesis, body.OrderAsc)
	if err != nil {
		logrus.Error("error while preparing query: ", err.Error())
		return fiber.ErrInternalServerError
	}
	defer stmt.Close()
	bees := []model.Bee{}
	err = stmt.SelectContext(c.Context(), &bees, sql.Named("NftId", uid))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(http.StatusOK).JSON(make([]string, 0))
		}
		logrus.Error("error while getting bees: ", err.Error())
		return fiber.ErrInternalServerError
	}
	beeResponses := []response.Bee{}
	for _, bee := range bees {
		beeResponse, err := b.FillResponseWithTrait(bee)
		if err != nil {
			logrus.Error("error while filling response")
			return fiber.ErrInternalServerError
		}
		cs, us, err := b.CalculateRanks(c, strconv.FormatUint(bee.Id, 10), body.IsGenesis)
		if err != nil {
			logrus.Error("error while calculating ranks")
			return fiber.ErrInternalServerError
		}
		beeResponse.Cosmetic = cs
		beeResponse.Utility = us
		beeResponse.HxdPerMinute = b.CalculateHxdPerTimeFrame(bee.Agility, time.Minute)
		beeResponse.HxdPerTwoHour = b.CalculateHxdPerTimeFrame(bee.Agility, 2*time.Hour)
		beeResponse.HxdCapacity = b.CalculateHxdCapacity(bee.Capacity)
		beeResponse.RecoveryTime = b.CalculateRecoveryTime(bee.Recovery)
		beeResponses = append(beeResponses, beeResponse)
	}
	return c.Status(http.StatusOK).JSON(beeResponses)
}

func (b *Bee) FillResponseWithTrait(m model.Bee) (response.Bee, error) {
	// aquire lock
	//to do
	beeResponse := response.Bee{
		Id:                 m.Id,
		Type:               m.Type,
		Generation:         m.Generation,
		Universe:           m.Universe,
		LandformSpecialty:  m.LandformSpecialty,
		Like:               m.Like,
		Dislike:            m.Dislike,
		Mood:               m.Mood,
		Health:             m.Health,
		Attack:             m.Attack,
		Defense:            m.Defense,
		Agility:            m.Agility,
		Luck:               m.Luck,
		Capacity:           m.Capacity,
		Recovery:           m.Recovery,
		Endurance:          m.Endurance,
		Level:              m.Level,
		LevelCap:           m.LevelCap,
		MateCap:            m.MateCap,
		NormalAttack1:      m.NormalAttack1,
		NormalAttack2:      m.NormalAttack2,
		SpecialAttack:      m.SpecialAttack,
		DateOfBirth:        m.DateOfBirth,
		Mother:             m.Mother,
		Father:             m.Father,
		Head:               m.Head,
		Eyes:               m.Eyes,
		Mouth:              m.Mouth,
		Clothes:            m.Clothes,
		BackFootAccessory:  m.BackFootAccessory,
		BackHandAccessory:  m.BackHandAccessory,
		FrontFootAccessory: m.FrontFootAccessory,
		FrontHandAccessory: m.FrontHandAccessory,
		BodyVisualTrait:    m.BodyVisualTrait,
		Background:         m.Background,
	}
	b.Btm.TR.BReadRarityMu.RLock()
	defer b.Btm.TR.BReadRarityMu.RUnlock()
	br := b.Btm.TR.BeeReadRarity
	_, ok := br["Generation"]
	if ok {
		beeResponse.Traits.Generation = br["Generation"][strconv.FormatUint(uint64(m.Generation), 10)]
	}
	_, ok = br["Universe"]
	if ok {
		beeResponse.Traits.Universe = br["Universe"][strconv.FormatUint(uint64(m.Universe), 10)]
	}
	_, ok = br["LandformSpecialty"]
	if ok {
		beeResponse.Traits.LandformSpecialty = br["LandformSpecialty"][strconv.FormatUint(uint64(m.LandformSpecialty), 10)]
	}
	_, ok = br["Like"]
	if ok {
		beeResponse.Traits.Like = br["Like"][strconv.FormatUint(uint64(m.Like), 10)]
	}
	_, ok = br["Dislike"]
	if ok {
		beeResponse.Traits.Dislike = br["Dislike"][strconv.FormatUint(uint64(m.Dislike), 10)]
	}
	_, ok = br["Mood"]
	if ok {
		beeResponse.Traits.Mood = br["Mood"][strconv.FormatUint(uint64(m.Mood), 10)]
	}
	_, ok = br["Level"]
	if ok {
		beeResponse.Traits.Level = br["Level"][strconv.FormatUint(uint64(m.Level), 10)]
	}
	_, ok = br["LevelCap"]
	if ok {
		beeResponse.Traits.LevelCap = br["LevelCap"][strconv.FormatUint(uint64(m.LevelCap), 10)]
	}
	_, ok = br["MateCap"]
	if ok {
		beeResponse.Traits.MateCap = br["MateCap"][strconv.FormatUint(uint64(m.MateCap), 10)]
	}
	_, ok = br["MateCount"]
	if ok {
		beeResponse.Traits.MateCount = br["MateCount"][strconv.FormatUint(uint64(m.MateCount), 10)]
	}
	_, ok = br["NormalAttack1"]
	if ok {
		beeResponse.Traits.NormalAttack1 = br["NormalAttack1"][strconv.FormatUint(uint64(m.NormalAttack1), 10)]
	}
	_, ok = br["NormalAttack2"]
	if ok {
		beeResponse.Traits.NormalAttack2 = br["NormalAttack2"][strconv.FormatUint(uint64(m.NormalAttack2), 10)]
	}
	_, ok = br["SpecialAttack"]
	if ok {
		beeResponse.Traits.SpecialAttack = br["SpecialAttack"][strconv.FormatUint(uint64(m.SpecialAttack), 10)]
	}
	_, ok = br["Head"]
	if ok {
		beeResponse.Traits.Head = br["Head"][strconv.FormatUint(uint64(m.Head), 10)]
	}
	_, ok = br["Eyes"]
	if ok {
		beeResponse.Traits.Eyes = br["Eyes"][strconv.FormatUint(uint64(m.Eyes), 10)]
	}
	_, ok = br["Mouth"]
	if ok {
		beeResponse.Traits.Mouth = br["Mouth"][strconv.FormatUint(uint64(m.Mouth), 10)]
	}
	_, ok = br["Feet"]
	if ok {
		beeResponse.Traits.Feet = br["Feet"][strconv.FormatUint(uint64(m.Feet), 10)]
	}
	_, ok = br["Clothes"]
	if ok {
		beeResponse.Traits.Clothes = br["Clothes"][strconv.FormatUint(uint64(m.Clothes), 10)]
	}
	_, ok = br["Hand"]
	if ok {
		beeResponse.Traits.Hand = br["Hand"][strconv.FormatUint(uint64(m.Hand), 10)]
	}
	_, ok = br["Hat"]
	if ok {
		beeResponse.Traits.Hat = br["Hat"][strconv.FormatUint(uint64(m.Hat), 10)]
	}
	_, ok = br["BackFootAccessory"]
	if ok {
		beeResponse.Traits.BackFootAccessory = br["BackFootAccessory"][strconv.FormatUint(uint64(m.BackFootAccessory), 10)]
	}
	_, ok = br["BackHandAccessory"]
	if ok {
		beeResponse.Traits.BackHandAccessory = br["BackHandAccessory"][strconv.FormatUint(uint64(m.BackHandAccessory), 10)]
	}
	_, ok = br["FrontFootAccessory"]
	if ok {
		beeResponse.Traits.FrontFootAccessory = br["FrontFootAccessory"][strconv.FormatUint(uint64(m.FrontFootAccessory), 10)]
	}
	_, ok = br["FrontHandAccessory"]
	if ok {
		beeResponse.Traits.FrontHandAccessory = br["FrontHandAccessory"][strconv.FormatUint(uint64(m.FrontHandAccessory), 10)]
	}
	_, ok = br["BodyVisualTrait"]
	if ok {
		beeResponse.Traits.BodyVisualTrait = br["BodyVisualTrait"][strconv.FormatUint(uint64(m.BodyVisualTrait), 10)]
	}
	_, ok = br["Background"]
	if ok {
		beeResponse.Traits.Background = br["Background"][strconv.FormatUint(uint64(m.Background), 10)]
	}
	return beeResponse, nil
}

func (b *Bee) CalculateRanks(c *fiber.Ctx, id string, isGenesis bool) (int16, int16, error) {
	if isGenesis {
		genesisBeeSetC := "genesis bee cosmic"
		genesisBeeSetU := "genesis bee utility"
		cosmeticS, err := b.Btm.RDB.ZRevRank(c.Context(), genesisBeeSetC, id)
		if err != nil {
			if errors.Is(err, redis.Nil) {
				logrus.Info("Genesis Bee %v is not part of sorted set", id)
				cosmeticS = -1
			} else {
				logrus.Error("error while getting cosmic score from redis: ", err.Error())
				return -1, -1, err
			}
		}
		utilityS, err := b.Btm.RDB.ZRevRank(c.Context(), genesisBeeSetU, id)
		if err != nil {
			if errors.Is(err, redis.Nil) {
				logrus.Info("Genesis Bee %v is not part of sorted set", id)
				utilityS = -1
			} else {
				logrus.Error("error while getting utility score from redis: ", err.Error())
				return -1, -1, err
			}
		}
		return int16(cosmeticS), int16(utilityS), nil
	} else {
		return -1, -1, nil
	}
}

func (b *Bee) CalculateHxdPerTimeFrame(stat uint16, d time.Duration) float64 {
	return 0.000012 * math.Pow(float64(stat), 1.3) * d.Minutes()
}
func (b *Bee) CalculateHxdCapacity(stat uint16) float64 {
	return 0.0008 * math.Pow(float64(stat), 1.3509)
}
func (b *Bee) CalculateRecoveryTime(stat uint16) int16 {
	return int16(math.Floor(3488 * math.Pow(float64(stat), -0.64)))
}

/// register handlers
func (b Bee) RegisterHandlers(g fiber.Router) {
	// register handlers
	g.Get("bee/:id", b.GetBeeById)
	g.Get("bees/:id", b.GetBeesById)
}
