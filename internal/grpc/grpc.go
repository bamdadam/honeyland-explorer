package grpc

import (
	"context"
	"net"

	"github.com/bamdadam/honeyland-explorer/internal/db/model"
	"github.com/bamdadam/honeyland-explorer/internal/notifier"
	"github.com/bamdadam/honeyland-explorer/internal/util/trait"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Grpc struct {
	Server *HoneylandNotifyServiceServer
}

type HoneylandNotifyServiceServer struct {
	notifier.UnimplementedNotifyServiceServer
	beeTChannel      chan model.BeeTrait
	beeTRMap         *trait.BeeTraitsRarityMap
	updateBeeChannel chan model.UpdateBeeGrpc
}

func NewGrpcServer(btm *trait.BeeTraitsRarityMap) *Grpc {
	btc := make(chan model.BeeTrait, 100)
	ubc := make(chan model.UpdateBeeGrpc, 100)
	return &Grpc{
		Server: &HoneylandNotifyServiceServer{
			beeTChannel:      btc,
			beeTRMap:         btm,
			updateBeeChannel: ubc,
		},
	}
}

func (g *Grpc) Init(ctx context.Context) {
	logrus.Info("Starting grpc server on port 8443")
	listener, err := net.Listen("tcp", ":8443")
	if err != nil {
		logrus.Fatal("Can't create grpc listener: ", err.Error())
	}
	serverRegistrar := grpc.NewServer()
	go g.UpdateBeeWriteMap(ctx)
	notifier.RegisterNotifyServiceServer(serverRegistrar, g.Server)
	if err := serverRegistrar.Serve(listener); err != nil {
		logrus.Fatal("Can't server grpc server: ", err.Error())
	}
}
