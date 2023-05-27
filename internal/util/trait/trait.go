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
}
