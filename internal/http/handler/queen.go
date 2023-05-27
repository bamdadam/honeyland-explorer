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

type Queen struct {
	DB *sqlx.DB
}

func (q *Queen) GetQueenById(c *fiber.Ctx) error {
	// get id
	id := c.Params("id")
	body := new(request.Queen)
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
	qq := queries.NewQueenQuery(q.DB)
	stmt, err := qq.LoadById(body.IsGenesis)
	if err != nil {
		logrus.Error("error while preparing query: ", err.Error())
		return fiber.ErrInternalServerError
	}
	defer stmt.Close()
	queen := model.Queen{}
	err = stmt.GetContext(c.Context(), &queen, sql.Named("NftId", uid))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(http.StatusOK).JSON(make([]string, 0))
		}
		logrus.Error("error while getting queen: ", err.Error())
		return fiber.ErrInternalServerError
	}
	return c.Status(http.StatusOK).JSON(queen)
}

func (q *Queen) GetQueensById(c *fiber.Ctx) error {
	// get id
	id := c.Params("id")
	body := new(request.Queens)
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
	qq := queries.NewQueenQuery(q.DB)
	stmt, err := qq.LoadByPage(body.IsGenesis, body.OrderAsc)
	if err != nil {
		logrus.Error("error while preparing query: ", err.Error())
		return fiber.ErrInternalServerError
	}
	defer stmt.Close()
	queens := []model.Queen{}
	err = stmt.SelectContext(c.Context(), &queens, sql.Named("NftId", uid))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.Status(http.StatusOK).JSON(make([]string, 0))
		}
		logrus.Error("error while getting queens: ", err.Error())
		return fiber.ErrInternalServerError
	}
	return c.Status(http.StatusOK).JSON(queens)
}

/// register handlers
func (q Queen) RegisterHandlers(g fiber.Router) {
	// register handlers
	g.Get("queen/:id", q.GetQueenById)
	g.Get("queens/:id", q.GetQueensById)
}
