package repository

import (
	"context"
	"game-player-data/internal/config"
	"game-player-data/internal/repository/model"
	"game-player-data/internal/repository/registrytypes"
	"game-player-data/internal/utils"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"sync"
	"time"
)

const (
	databaseName = "game-player-data"

	blockSumoCollection    = "blockSumo"
	towerDefenceCollection = "towerDefence"
	minesweeperCollection  = "minesweeper"
)

type mongoRepository struct {
	database *mongo.Database

	blockSumoCollection    *mongo.Collection
	towerDefenceCollection *mongo.Collection
	minesweeperCollection  *mongo.Collection
}

func NewMongoRepository(ctx context.Context, logger *zap.SugaredLogger, wg *sync.WaitGroup, cfg *config.MongoDBConfig) (Repository, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.URI).SetRegistry(createCodecRegistry()))
	if err != nil {
		return nil, err
	}

	database := client.Database(databaseName)
	repo := &mongoRepository{
		database:               database,
		blockSumoCollection:    database.Collection(blockSumoCollection),
		towerDefenceCollection: database.Collection(towerDefenceCollection),
		minesweeperCollection:  database.Collection(minesweeperCollection),
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		if err := client.Disconnect(ctx); err != nil {
			logger.Errorw("failed to disconnect from mongo", err)
		}
	}()

	return repo, nil
}

func (m *mongoRepository) GetBlockSumoPlayer(ctx context.Context, playerId uuid.UUID) (*model.BlockSumoData, error) {
	result := m.getData(ctx, playerId, m.blockSumoCollection)

	var data model.BlockSumoData
	if err := result.Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

func (m *mongoRepository) SaveBlockSumoPlayer(ctx context.Context, data *model.BlockSumoData) error {
	return m.saveData(ctx, data.PlayerId, data, m.blockSumoCollection)
}

func (m *mongoRepository) GetTowerDefencePlayer(ctx context.Context, playerId uuid.UUID) (*model.TowerDefenceData, error) {
	result := m.getData(ctx, playerId, m.blockSumoCollection)

	var data model.TowerDefenceData
	if err := result.Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

func (m *mongoRepository) SaveTowerDefencePlayer(ctx context.Context, data *model.TowerDefenceData) error {
	return m.saveData(ctx, data.PlayerId, data, m.towerDefenceCollection)
}

func (m *mongoRepository) GetMinesweeperPlayer(ctx context.Context, playerId uuid.UUID) (*model.MinesweeperData, error) {
	result := m.getData(ctx, playerId, m.blockSumoCollection)

	var data model.MinesweeperData
	if err := result.Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

func (m *mongoRepository) SaveMinesweeperPlayer(ctx context.Context, data *model.MinesweeperData) error {
	return m.saveData(ctx, data.PlayerId, data, m.minesweeperCollection)
}

func (m *mongoRepository) getData(ctx context.Context, playerId uuid.UUID, collection *mongo.Collection) *mongo.SingleResult {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return collection.FindOne(ctx, bson.M{"_id": playerId})
}

func (m *mongoRepository) saveData(ctx context.Context, playerId uuid.UUID, data interface{}, collection *mongo.Collection) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": playerId}

	if _, err := collection.ReplaceOne(ctx, filter, data, &options.ReplaceOptions{Upsert: utils.PointerOf(true)}); err != nil {
		return err
	}

	return nil
}

func createCodecRegistry() *bsoncodec.Registry {
	return bson.NewRegistryBuilder().
		RegisterTypeEncoder(registrytypes.UUIDType, bsoncodec.ValueEncoderFunc(registrytypes.UuidEncodeValue)).
		RegisterTypeDecoder(registrytypes.UUIDType, bsoncodec.ValueDecoderFunc(registrytypes.UuidDecodeValue)).
		Build()
}
