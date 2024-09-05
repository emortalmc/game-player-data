package repository

import (
	"context"
	"game-player-data/internal/repository/model"
	"github.com/google/uuid"
)

type Repository interface {
	GetBlockSumoData(ctx context.Context, playerId uuid.UUID) (*model.BlockSumoData, error)
	GetBlockSumoDataForPlayers(ctx context.Context, playerIds []uuid.UUID) ([]*model.BlockSumoData, error)
	SaveBlockSumoPlayer(ctx context.Context, data *model.BlockSumoData) error

	GetTowerDefencePlayer(ctx context.Context, playerId uuid.UUID) (*model.TowerDefenceData, error)
	SaveTowerDefencePlayer(ctx context.Context, data *model.TowerDefenceData) error

	GetMinesweeperPlayer(ctx context.Context, playerId uuid.UUID) (*model.MinesweeperData, error)
	SaveMinesweeperPlayer(ctx context.Context, data *model.MinesweeperData) error
}
