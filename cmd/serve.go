package cmd

import (
	"context"
	"time"

	"github.com/bamdadam/honeyland-explorer/internal/db/rdb"
	"github.com/bamdadam/honeyland-explorer/internal/grpc"
	"github.com/bamdadam/honeyland-explorer/internal/http/handler"
	"github.com/bamdadam/honeyland-explorer/internal/util/trait"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	_ "github.com/microsoft/go-mssqldb"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func Register(root *cobra.Command) {
	root.AddCommand(&cobra.Command{
		Use:   "serve",
		Short: "run server",
		Run: func(cmd *cobra.Command, args []string) {
			main()
		},
	})
}

func main() {
	log.Info("starting server on port 1373")

	app := fiber.New(
		fiber.Config{
			AppName: "Honeyland explorer",
		},
	)
	connString := ("Data Source=Honeyreplica.database.windows.net;Initial Catalog=Honeyverse;user id=onlyRead; password=readOn1y;Connection Timeout=30")

	db, err := sqlx.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
	}
	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal("Error pinging database: ", err.Error())
	}
	rdb, err := rdb.New(ctx)
	if err != nil {
		log.Fatal("Error while making rdb client: ", err.Error())
	}
	btm := trait.BeeTraitsRarityMap{
		DB:              db,
		RDB:             rdb,
		NFTNumFieldName: "NFTNum",
	}
	err = btm.InitGenesis(ctx)
	if err != nil {
		log.Fatal("Error while getting genesis bee rarity map: ", err.Error())
	}
	err = btm.InitGeneration(ctx)
	if err != nil {
		log.Fatal("Error while getting generation bee rarity map: ", err.Error())
	}
	go btm.SetBeeSetsScheduler(ctx, 5*time.Minute)
	grpcService := grpc.NewGrpcServer(&btm)
	go grpcService.Init(ctx)
	bh := handler.Bee{
		DB:  db,
		Btm: &btm,
		RDB: rdb,
	}
	eh := handler.Egg{
		DB: db,
	}
	qh := handler.Queen{
		DB: db,
	}
	lh := handler.Land{
		DB: db,
	}
	g := app.Group("/")
	bh.RegisterHandlers(g)
	eh.RegisterHandlers(g)
	qh.RegisterHandlers(g)
	lh.RegisterHandlers(g)

	if err := app.Listen(":1373"); err != nil {
		log.Fatal("Can't start server: ", err.Error())
	}
}
