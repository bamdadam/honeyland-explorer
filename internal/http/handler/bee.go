package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/bamdadam/honeyland-explorer/internal/db/model"
	"github.com/bamdadam/honeyland-explorer/internal/http/request"
	"github.com/bamdadam/honeyland-explorer/internal/queries"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type Bee struct {
	DB *sqlx.DB
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
	return c.Status(http.StatusOK).JSON(bee)
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
	return c.Status(http.StatusOK).JSON(bees)
}

/// register handlers
func (b Bee) RegisterHandlers(g fiber.Router) {
	// register handlers
	g.Get("bee/:id", b.GetBeeById)
	g.Get("bees/:id", b.GetBeesById)
}
