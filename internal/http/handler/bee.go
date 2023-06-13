package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/bamdadam/honeyland-explorer/internal/db/model"
	"github.com/bamdadam/honeyland-explorer/internal/db/rdb"
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
	RDB *rdb.RedisDB
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
	beeResponse, err := b.FillResponseWithTrait(c.Context(), bee, body.IsGenesis)
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
		beeResponse, err := b.FillResponseWithTrait(c.Context(), bee, body.IsGenesis)
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

func (b *Bee) FillResponseWithTrait(ctx context.Context, m model.Bee, isGenesis bool) (response.Bee, error) {
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
	pre := "Generation-Bee-"
	if isGenesis {
		pre = "Genesis-Bee-"
	}
	nftNum, err := b.RDB.HGetUint(ctx, pre+b.Btm.NFTNumFieldName, "0")
	if err != nil {
		if errors.Is(err, redis.Nil) {
			logrus.Errorf("nftNum %v is not part of rarity map", b.Btm.NFTNumFieldName)
		} else {
			logrus.Error("error while getting nftNum rarity from redis: ", err.Error())
		}
		return beeResponse, nil
	}
	g, err := b.RDB.HGetUint(ctx, pre+"Generation", strconv.FormatUint(uint64(m.Generation), 10))
	if err != nil {
		if errors.Is(err, redis.Nil) {
			logrus.Errorf("Generation %v is not part of rarity map", m.Generation)
		} else {
			logrus.Error("error while getting generation rarity from redis: ", err.Error())
		}
	}
	beeResponse.Traits.Generation = strconv.FormatFloat((float64(g)/float64(nftNum))*100, 'f', 3, 64)
	u, err := b.RDB.HGetUint(ctx, pre+"Universe", strconv.FormatUint(uint64(m.Universe), 10))
	if err != nil {
		if errors.Is(err, redis.Nil) {
			logrus.Errorf("Universe %v is not part of rarity map", m.Universe)
		} else {
			logrus.Error("error while getting universe rarity from redis: ", err.Error())
		}
	}
	beeResponse.Traits.Universe = strconv.FormatFloat((float64(u)/float64(nftNum))*100, 'f', 3, 64)
	ls, err := b.RDB.HGetUint(ctx, pre+"LandformSpecialty", strconv.FormatUint(uint64(m.LandformSpecialty), 10))
	if err != nil {
		if errors.Is(err, redis.Nil) {
			logrus.Errorf("LandformSpecialty %v is not part of rarity map", m.LandformSpecialty)
		} else {
			logrus.Error("error while getting landformSpecialty rarity from redis: ", err.Error())
		}
	}
	beeResponse.Traits.LandformSpecialty = strconv.FormatFloat((float64(ls)/float64(nftNum))*100, 'f', 3, 64)
	l, err := b.RDB.HGetUint(ctx, pre+"Like", strconv.FormatUint(uint64(m.Like), 10))
	if err != nil {
		if errors.Is(err, redis.Nil) {
			logrus.Errorf("Like %v is not part of rarity map", m.Like)
		} else {
			logrus.Error("error while getting like rarity from redis: ", err.Error())
		}
	}
	beeResponse.Traits.Like = strconv.FormatFloat((float64(l)/float64(nftNum))*100, 'f', 3, 64)
	d, err := b.RDB.HGetUint(ctx, pre+"Dislike", strconv.FormatUint(uint64(m.Dislike), 10))
	if err != nil {
		if errors.Is(err, redis.Nil) {
			logrus.Errorf("Dislike %v is not part of rarity map", m.Dislike)
		} else {
			logrus.Error("error while getting dislike rarity from redis: ", err.Error())
		}
	}
	beeResponse.Traits.Dislike = strconv.FormatFloat((float64(d)/float64(nftNum))*100, 'f', 3, 64)
	mo, err := b.RDB.HGetUint(ctx, pre+"Mood", strconv.FormatUint(uint64(m.Mood), 10))
	if err != nil {
		if errors.Is(err, redis.Nil) {
			logrus.Errorf("Mood %v is not part of rarity map", m.Mood)
		} else {
			logrus.Error("error while getting mood rarity from redis: ", err.Error())
		}
	}
	beeResponse.Traits.Mood = strconv.FormatFloat((float64(mo)/float64(nftNum))*100, 'f', 3, 64)
	le, err := b.RDB.HGetUint(ctx, pre+"Level", strconv.FormatUint(uint64(m.Level), 10))
	if err != nil {
		if errors.Is(err, redis.Nil) {
			logrus.Errorf("Level %v is not part of rarity map", m.Level)
		} else {
			logrus.Error("error while getting Level rarity from redis: ", err.Error())
		}
	}
	beeResponse.Traits.Level = strconv.FormatFloat((float64(le)/float64(nftNum))*100, 'f', 3, 64)
	lc, err := b.RDB.HGetUint(ctx, pre+"LevelCap", strconv.FormatUint(uint64(m.LevelCap), 10))
	if err != nil {
		if errors.Is(err, redis.Nil) {
			logrus.Errorf("LevelCap %v is not part of rarity map", m.LevelCap)
		} else {
			logrus.Error("error while getting LevelCap rarity from redis: ", err.Error())
		}
	}
	beeResponse.Traits.LevelCap = strconv.FormatFloat((float64(lc)/float64(nftNum))*100, 'f', 3, 64)
	mc, err := b.RDB.HGetUint(ctx, pre+"MateCap", strconv.FormatUint(uint64(m.MateCap), 10))
	if err != nil {
		if errors.Is(err, redis.Nil) {
			logrus.Errorf("MateCap %v is not part of rarity map", m.MateCap)
		} else {
			logrus.Error("error while getting MateCap rarity from redis: ", err.Error())
		}
	}
	beeResponse.Traits.MateCap = strconv.FormatFloat((float64(mc)/float64(nftNum))*100, 'f', 3, 64)
	mco, err := b.RDB.HGetUint(ctx, pre+"MateCount", strconv.FormatUint(uint64(m.MateCount), 10))
	if err != nil {
		if errors.Is(err, redis.Nil) {
			logrus.Errorf("MateCount %v is not part of rarity map", m.MateCount)
		} else {
			logrus.Error("error while getting MateCount rarity from redis: ", err.Error())
		}
	}
	beeResponse.Traits.MateCount = strconv.FormatFloat((float64(mco)/float64(nftNum))*100, 'f', 3, 64)
	n1, err := b.RDB.HGetUint(ctx, pre+"NormalAttack1", strconv.FormatUint(uint64(m.NormalAttack1), 10))
	if err != nil {
		if errors.Is(err, redis.Nil) {
			logrus.Errorf("NormalAttack1 %v is not part of rarity map", m.NormalAttack1)
		} else {
			logrus.Error("error while getting NormalAttack1 rarity from redis: ", err.Error())
		}
	}
	beeResponse.Traits.NormalAttack1 = strconv.FormatFloat((float64(n1)/float64(nftNum))*100, 'f', 3, 64)
	n2, err := b.RDB.HGetUint(ctx, pre+"NormalAttack2", strconv.FormatUint(uint64(m.NormalAttack2), 10))
	if err != nil {
		if errors.Is(err, redis.Nil) {
			logrus.Errorf("NormalAttack2 %v is not part of rarity map", m.NormalAttack2)
		} else {
			logrus.Error("error while getting NormalAttack2 rarity from redis: ", err.Error())
		}
	}
	beeResponse.Traits.NormalAttack2 = strconv.FormatFloat((float64(n2)/float64(nftNum))*100, 'f', 3, 64)
	sa, err := b.RDB.HGetUint(ctx, pre+"SpecialAttack", strconv.FormatUint(uint64(m.SpecialAttack), 10))
	if err != nil {
		if errors.Is(err, redis.Nil) {
			logrus.Errorf("SpecialAttack %v is not part of rarity map", m.SpecialAttack)
		} else {
			logrus.Error("error while getting SpecialAttack rarity from redis: ", err.Error())
		}
	}
	beeResponse.Traits.SpecialAttack = strconv.FormatFloat((float64(sa)/float64(nftNum))*100, 'f', 3, 64)
	h, err := b.RDB.HGetUint(ctx, pre+"Head", strconv.FormatUint(uint64(m.Head), 10))
	if err != nil {
		if errors.Is(err, redis.Nil) {
			logrus.Errorf("Head %v is not part of rarity map", m.Head)
		} else {
			logrus.Error("error while getting Head rarity from redis: ", err.Error())
		}
	}
	beeResponse.Traits.Head = strconv.FormatFloat((float64(h)/float64(nftNum))*100, 'f', 3, 64)
	e, err := b.RDB.HGetUint(ctx, pre+"Eyes", strconv.FormatUint(uint64(m.Eyes), 10))
	if err != nil {
		if errors.Is(err, redis.Nil) {
			logrus.Errorf("Eyes %v is not part of rarity map", m.Eyes)
		} else {
			logrus.Error("error while getting Eyes rarity from redis: ", err.Error())
		}
	}
	beeResponse.Traits.Eyes = strconv.FormatFloat((float64(e)/float64(nftNum))*100, 'f', 3, 64)
	mou, err := b.RDB.HGetUint(ctx, pre+"Mouth", strconv.FormatUint(uint64(m.Mouth), 10))
	if err != nil {
		if errors.Is(err, redis.Nil) {
			logrus.Errorf("Mouth %v is not part of rarity map", m.Mouth)
		} else {
			logrus.Error("error while getting Mouth rarity from redis: ", err.Error())
		}
	}
	beeResponse.Traits.Mouth = strconv.FormatFloat((float64(mou)/float64(nftNum))*100, 'f', 3, 64)
	f, err := b.RDB.HGetUint(ctx, pre+"Feet", strconv.FormatUint(uint64(m.Feet), 10))
	if err != nil {
		if errors.Is(err, redis.Nil) {
			logrus.Errorf("Feet %v is not part of rarity map", m.Feet)
		} else {
			logrus.Error("error while getting Feet rarity from redis: ", err.Error())
		}
	}
	beeResponse.Traits.Feet = strconv.FormatFloat((float64(f)/float64(nftNum))*100, 'f', 3, 64)
	cl, err := b.RDB.HGetUint(ctx, pre+"Clothes", strconv.FormatUint(uint64(m.Clothes), 10))
	if err != nil {
		if errors.Is(err, redis.Nil) {
			logrus.Errorf("Clothes %v is not part of rarity map", m.Clothes)
		} else {
			logrus.Error("error while getting Clothes rarity from redis: ", err.Error())
		}
	}
	beeResponse.Traits.Clothes = strconv.FormatFloat((float64(cl)/float64(nftNum))*100, 'f', 3, 64)
	ha, err := b.RDB.HGetUint(ctx, pre+"Hand", strconv.FormatUint(uint64(m.Hand), 10))
	if err != nil {
		if errors.Is(err, redis.Nil) {
			logrus.Errorf("Hand %v is not part of rarity map", m.Hand)
		} else {
			logrus.Error("error while getting Hand rarity from redis: ", err.Error())
		}
	}
	beeResponse.Traits.Hand = strconv.FormatFloat((float64(ha)/float64(nftNum))*100, 'f', 3, 64)
	hat, err := b.RDB.HGetUint(ctx, pre+"Hat", strconv.FormatUint(uint64(m.Hat), 10))
	if err != nil {
		if errors.Is(err, redis.Nil) {
			logrus.Errorf("Hat %v is not part of rarity map", m.Hat)
		} else {
			logrus.Error("error while getting Hat rarity from redis: ", err.Error())
		}
	}
	beeResponse.Traits.Hat = strconv.FormatFloat((float64(hat)/float64(nftNum))*100, 'f', 3, 64)
	bf, err := b.RDB.HGetUint(ctx, pre+"BackFootAccessory", strconv.FormatUint(uint64(m.BackFootAccessory), 10))
	if err != nil {
		if errors.Is(err, redis.Nil) {
			logrus.Errorf("BackFootAccessory %v is not part of rarity map", m.BackFootAccessory)
		} else {
			logrus.Error("error while getting BackFootAccessory rarity from redis: ", err.Error())
		}
	}
	beeResponse.Traits.BackFootAccessory = strconv.FormatFloat((float64(bf)/float64(nftNum))*100, 'f', 3, 64)
	bh, err := b.RDB.HGetUint(ctx, pre+"BackHandAccessory", strconv.FormatUint(uint64(m.BackHandAccessory), 10))
	if err != nil {
		if errors.Is(err, redis.Nil) {
			logrus.Errorf("BackHandAccessory %v is not part of rarity map", m.BackHandAccessory)
		} else {
			logrus.Error("error while getting BackHandAccessory rarity from redis: ", err.Error())
		}
	}
	beeResponse.Traits.BackHandAccessory = strconv.FormatFloat((float64(bh)/float64(nftNum))*100, 'f', 3, 64)
	ff, err := b.RDB.HGetUint(ctx, pre+"FrontFootAccessory", strconv.FormatUint(uint64(m.FrontFootAccessory), 10))
	if err != nil {
		if errors.Is(err, redis.Nil) {
			logrus.Errorf("FrontFootAccessory %v is not part of rarity map", m.FrontFootAccessory)
		} else {
			logrus.Error("error while getting FrontFootAccessory rarity from redis: ", err.Error())
		}
	}
	beeResponse.Traits.FrontFootAccessory = strconv.FormatFloat((float64(ff)/float64(nftNum))*100, 'f', 3, 64)
	fh, err := b.RDB.HGetUint(ctx, pre+"FrontHandAccessory", strconv.FormatUint(uint64(m.FrontHandAccessory), 10))
	if err != nil {
		if errors.Is(err, redis.Nil) {
			logrus.Errorf("FrontHandAccessory %v is not part of rarity map", m.FrontHandAccessory)
		} else {
			logrus.Error("error while getting FrontHandAccessory rarity from redis: ", err.Error())
		}
	}
	beeResponse.Traits.FrontHandAccessory = strconv.FormatFloat((float64(fh)/float64(nftNum))*100, 'f', 3, 64)
	bv, err := b.RDB.HGetUint(ctx, pre+"BodyVisualTrait", strconv.FormatUint(uint64(m.BodyVisualTrait), 10))
	if err != nil {
		if errors.Is(err, redis.Nil) {
			logrus.Errorf("BodyVisualTrait %v is not part of rarity map", m.BodyVisualTrait)
		} else {
			logrus.Error("error while getting BodyVisualTrait rarity from redis: ", err.Error())
		}
	}
	beeResponse.Traits.BodyVisualTrait = strconv.FormatFloat((float64(bv)/float64(nftNum))*100, 'f', 3, 64)
	bc, err := b.RDB.HGetUint(ctx, pre+"Background", strconv.FormatUint(uint64(m.Background), 10))
	if err != nil {
		if errors.Is(err, redis.Nil) {
			logrus.Errorf("Background %v is not part of rarity map", m.Background)
		} else {
			logrus.Error("error while getting Background rarity from redis: ", err.Error())
		}
	}
	beeResponse.Traits.Background = strconv.FormatFloat((float64(bc)/float64(nftNum))*100, 'f', 3, 64)
	return beeResponse, nil
}

func (b *Bee) CalculateRanks(c *fiber.Ctx, id string, isGenesis bool) (int16, int16, error) {
	if isGenesis {
		genesisBeeSetC := "genesis bee cosmic"
		genesisBeeSetU := "genesis bee utility"
		cosmeticS, err := b.Btm.RDB.ZRevRank(c.Context(), genesisBeeSetC, id)
		if err != nil {
			if errors.Is(err, redis.Nil) {
				logrus.Infof("Genesis Bee %v is not part of sorted set", id)
				cosmeticS = -1
			} else {
				logrus.Error("error while getting cosmic score from redis: ", err.Error())
				return -1, -1, err
			}
		}
		utilityS, err := b.Btm.RDB.ZRevRank(c.Context(), genesisBeeSetU, id)
		if err != nil {
			if errors.Is(err, redis.Nil) {
				logrus.Infof("Genesis Bee %v is not part of sorted set", id)
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

// / register handlers
func (b Bee) RegisterHandlers(g fiber.Router) {
	// register handlers
	g.Get("bee/:id", b.GetBeeById)
	g.Get("bees/:id", b.GetBeesById)
}
