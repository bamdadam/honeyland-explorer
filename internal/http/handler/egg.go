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

type Egg struct {
	DB *sqlx.DB
}

func (e *Egg) GetEggById(c *fiber.Ctx) error {
	id := c.Params("id")
	body := new(request.Egg)
	err := c.BodyParser(body)
	if err != nil {
		logrus.Error("can't parse request body: ", err.Error())
		return fiber.ErrBadRequest
	}
	//conver uid to u64
	uid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		logrus.Error("can't parse id to uint64: ", err.Error())
		return fiber.ErrBadRequest
	}
	eq := queries.NewEggQuery(e.DB)
	stmt, err := eq.LoadById(body.IsGenesis)
	if err != nil {
		logrus.Error("Error while preparing query: ", err.Error())
		return fiber.ErrInternalServerError
	}
	defer stmt.Close()
	egg := new(model.Egg)
	err = stmt.GetContext(c.Context(), egg, sql.Named("NftId", uid))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(http.StatusOK).JSON(make([]string, 0))
		}
		logrus.Error("Error while getting egg: ", err.Error())
		return fiber.ErrInternalServerError
	}
	return c.Status(http.StatusOK).JSON(egg)
}

func (e *Egg) GetEggsById(c *fiber.Ctx) error {
	id := c.Params("id")
	body := new(request.Eggs)
	err := c.BodyParser(body)
	if err != nil {
		logrus.Error("can't parse request body: ", err.Error())
		return fiber.ErrBadRequest
	}
	//conver uid to u64
	uid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		logrus.Error("can't parse id to uint64: ", err.Error())
		return fiber.ErrBadRequest
	}
	eq := queries.NewEggQuery(e.DB)
	stmt, err := eq.LoadByPage(body.IsGenesis, body.OrderAsc)
	if err != nil {
		logrus.Error("Error while preparing query: ", err.Error())
		return fiber.ErrInternalServerError
	}
	defer stmt.Close()
	eggs := new([]model.Egg)
	err = stmt.SelectContext(c.Context(), eggs, sql.Named("NftId", uid))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(http.StatusOK).JSON(make([]string, 0))
		}
		logrus.Error("Error while getting egg: ", err.Error())
		return fiber.ErrInternalServerError
	}
	return c.Status(http.StatusOK).JSON(eggs)
}

func (e *Egg) RegisterHandlers(g fiber.Router) {
	g.Get("egg/:id", e.GetEggById)
	g.Get("eggs/:id", e.GetEggsById)
}
