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

type Land struct {
	DB *sqlx.DB
}

func (l *Land) getLandById(c *fiber.Ctx) error {
	id := c.Params("id")
	uid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		logrus.Error("can't parse id to uint64: ", err.Error())
		return fiber.ErrBadRequest
	}
	if uid > 2000 {
		logrus.Error("lands higher than 2000 are not supported yet")
		return c.Status(http.StatusOK).JSON(make([]string, 0))
	}
	lq := queries.NewLandQuery(l.DB)
	stmt, err := lq.LoadById()
	if err != nil {
		logrus.Error("error while preparing query: ", err.Error())
		return fiber.ErrInternalServerError
	}
	defer stmt.Close()
	land := new(model.Land)
	err = stmt.GetContext(c.Context(), land, sql.Named("NftId", uid))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(http.StatusOK).JSON(make([]string, 0))
		}
		logrus.Error("error while getting land: ", err.Error())
		return fiber.ErrInternalServerError
	}
	return c.Status(http.StatusOK).JSON(land)
}

func (l *Land) getLandByPage(c *fiber.Ctx) error {
	id := c.Params("id")
	uid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		logrus.Error("can't parse id to uint64: ", err.Error())
		return fiber.ErrBadRequest
	}
	if uid > 1971 {
		logrus.Error("lands higher than 2000 are not supported yet")
		return c.Status(http.StatusOK).JSON(make([]string, 0))
	}
	body := new(request.Lands)
	err = c.BodyParser(body)
	if err != nil {
		logrus.Error("can't parse body: ", err.Error())
		return fiber.ErrBadRequest
	}
	lq := queries.NewLandQuery(l.DB)
	stmt, err := lq.LoadByPage(body.OrderAsc)
	if err != nil {
		logrus.Error("error while preparing query: ", err.Error())
		return fiber.ErrInternalServerError
	}
	defer stmt.Close()
	lands := new([]model.Land)
	err = stmt.SelectContext(c.Context(), lands, sql.Named("NftId", uid))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(http.StatusOK).JSON(make([]string, 0))
		}
		logrus.Error("error while getting land: ", err.Error())
		return fiber.ErrInternalServerError
	}
	return c.Status(http.StatusOK).JSON(lands)
}

func (l *Land) RegisterHandlers(g fiber.Router) {
	g.Get("land/:id", l.getLandById)
	g.Get("lands/:id", l.getLandByPage)
}
