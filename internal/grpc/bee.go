package grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/bamdadam/honeyland-explorer/internal/db/model"
	"github.com/bamdadam/honeyland-explorer/internal/notifier"
	"github.com/sirupsen/logrus"
)

func (g *Grpc) SetBeeReadMapScheduler(ctx context.Context, d time.Duration) {
	ticker := time.NewTicker(d)
	go func() {
		for range ticker.C {
			bpm, err := g.Server.beeTRMap.Convert(ctx)
			if err != nil {
				logrus.Error("Error while converting Bee write map to Bee percentage map: ", err.Error())
			} else {
				g.Server.beeTRMap.SetReadMap(ctx, bpm)
				fmt.Println("Task executed!")
				jsonrm, err := json.Marshal(g.Server.beeTRMap.TR.BeeReadRarity)
				if err != nil {
					logrus.Error("Can't conver rarity map to json: ", err.Error())
				} else {
					fmt.Println(string(jsonrm))
				}
			}
		}
	}()
}

func (s *Grpc) NotifyBee(ctx context.Context, req *notifier.BeeNotifierRequest) (*notifier.BeeNotifierResponse, error) {
	bt := model.BeeTrait{}
	if req != nil {
		bt = model.BeeTrait{
			Generation:         uint16(req.GetGeneration()),
			Universe:           uint16(req.GetUniverse()),
			LandformSpecialty:  uint16(req.GetLandformSpecialty()),
			Like:               uint16(req.GetLike()),
			Dislike:            uint16(req.GetDislike()),
			Mood:               uint16(req.GetMood()),
			Level:              uint16(req.GetLevel()),
			LevelCap:           uint16(req.GetLevelCap()),
			MateCap:            uint16(req.GetMateCap()),
			MateCount:          uint16(req.GetMateCount()),
			NormalAttack1:      uint16(req.GetNormalAttack1()),
			NormalAttack2:      uint16(req.GetNormalAttack2()),
			SpecialAttack:      uint16(req.GetSpecialAttack()),
			Head:               uint16(req.GetHead()),
			Eyes:               uint16(req.GetEyes()),
			Mouth:              uint16(req.GetMouth()),
			Feet:               uint16(req.GetFeet()),
			Clothes:            uint16(req.GetClothes()),
			Hand:               uint16(req.GetHand()),
			Hat:                uint16(req.GetHat()),
			BackFootAccessory:  uint16(req.GetBackFootAccessory()),
			BackHandAccessory:  uint16(req.GetBackHandAccessory()),
			FrontFootAccessory: uint16(req.GetFrontFootAccessory()),
			FrontHandAccessory: uint16(req.GetFrontHandAccessory()),
			BodyVisualTrait:    uint16(req.GetBodyVisualTrait()),
			Background:         uint16(req.GetBackground()),
		}
		s.Server.beeTChannel <- bt
		return &notifier.BeeNotifierResponse{
			Success: true,
			Message: "",
		}, nil
	}
	return &notifier.BeeNotifierResponse{
		Success: false,
		Message: "Request is null",
	}, nil
}

func (s *Grpc) UpdateBeeWriteMap(ctx context.Context) {
	for bt := range s.Server.beeTChannel {
		err := s.Server.beeTRMap.UpdateWriteMap(ctx, bt)
		if err != nil {
			logrus.Error("Error while updating bee write map: ", err.Error())
		}
	}
}
