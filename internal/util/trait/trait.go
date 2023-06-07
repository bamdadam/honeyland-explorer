package trait

import (
	"context"
)

type NftTraitsRarityMap interface {
	Init(context.Context) (map[string]map[string]uint32, error)
	SetWriteMap(context.Context, map[string]map[string]uint32) error
	UpdateWriteMap(context.Context, interface{}) error
	Convert(context.Context, map[string]map[string]uint32) (map[string]map[string]string, error)
	SetReadMap(context.Context, map[string]map[string]string) error
	FindRarityMaxValue(map[string]uint32) uint32
}

// type TraitsRarity struct {
// 	BeeTRMap        *BeeTraitsRarityMap
// 	QueenReadRarity map[string]map[string]string
// 	QReadRarityMu   sync.RWMutex
// 	QueenRarity     map[string]map[string]uint32
// 	QRarityMu       sync.RWMutex

// 	NFTNumFieldName string
// }

// func (t *TraitsRarity) FindRarityMaxValue(m map[string]uint32) uint32 {
// 	var maxValue uint32
// 	// Iterate over the map and update the maximum value if a greater value is found
// 	for _, value := range m {
// 		if value > maxValue {
// 			maxValue = value
// 		}
// 	}
// 	return maxValue
// }
