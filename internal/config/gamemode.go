package config

import (
	"game-player-data/internal/repository/model"
	"github.com/google/uuid"
)

func CreateDefaultBlockSumoData(playerId uuid.UUID) *model.BlockSumoData {
	return &model.BlockSumoData{
		PlayerId:   playerId,
		BlockSlot:  2,
		ShearsSlot: 1,
	}
}
