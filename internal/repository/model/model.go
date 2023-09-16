package model

import (
	"github.com/emortalmc/proto-specs/gen/go/model/gameplayerdata"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/anypb"
)

type GameData interface {
	ToAnyProto() (*anypb.Any, error)
}

type BlockSumoData struct {
	PlayerId uuid.UUID `bson:"_id"`

	BlockSlot  uint32 `bson:"blockSlot"`
	ShearsSlot uint32 `bson:"shearsSlot"`
}

func (d *BlockSumoData) ToAnyProto() (*anypb.Any, error) {
	return anypb.New(&gameplayerdata.BlockSumoPlayerData{
		BlockSlot:  d.BlockSlot,
		ShearsSlot: d.ShearsSlot,
	})
}

func (d *BlockSumoData) FromProto(pId uuid.UUID, data *gameplayerdata.BlockSumoPlayerData) {
	d.PlayerId = pId
	d.BlockSlot = data.BlockSlot
	d.ShearsSlot = data.ShearsSlot
}

// MinesweeperData TODO
type MinesweeperData struct {
	PlayerId uuid.UUID `bson:"_id"`
}

// TODO
func (d *MinesweeperData) ToProto() *gameplayerdata.MinesweeperPlayerData {
	return &gameplayerdata.MinesweeperPlayerData{}
}

// TowerDefenceData TODO
type TowerDefenceData struct {
	PlayerId uuid.UUID `bson:"_id"`
}

// TODO
func (d *TowerDefenceData) ToProto() *gameplayerdata.TowerDefencePlayerData {
	return &gameplayerdata.TowerDefencePlayerData{}
}
