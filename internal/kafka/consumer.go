package kafka

import (
	"context"
	"errors"
	"fmt"
	"game-player-data/internal/config"
	"game-player-data/internal/repository"
	pbmsg "github.com/emortalmc/proto-specs/gen/go/message/gameplayerdata"
	pbmodel "github.com/emortalmc/proto-specs/gen/go/model/gameplayerdata"
	"github.com/emortalmc/proto-specs/gen/go/nongenerated/kafkautils"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"sync"
)

const GamePlayerDataTopic = "game-player-data"

type consumer struct {
	logger *zap.SugaredLogger
	repo   repository.Repository

	reader *kafka.Reader
}

func NewConsumer(ctx context.Context, wg *sync.WaitGroup, cfg *config.KafkaConfig, logger *zap.SugaredLogger,
	repo repository.Repository) {

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)},
		GroupID:     "game-player-data",
		GroupTopics: []string{GamePlayerDataTopic},

		Logger: kafka.LoggerFunc(func(format string, args ...interface{}) {
			logger.Infow(fmt.Sprintf(format, args...))
		}),
		ErrorLogger: kafka.LoggerFunc(func(format string, args ...interface{}) {
			logger.Errorw(fmt.Sprintf(format, args...))
		}),
	})

	c := &consumer{
		logger: logger,
		repo:   repo,

		reader: reader,
	}

	handler := kafkautils.NewConsumerHandler(logger, reader)
	handler.RegisterHandler(&pbmsg.UpdateGamePlayerDataMessage{}, c.handleUpdateGamePlayerDataMessage)

	logger.Infow("starting listening for kafka messages", "topics", reader.Config().GroupTopics)

	wg.Add(1)
	go func() {
		defer wg.Done()
		handler.Run(ctx) // Run is blocking until the context is cancelled
		if err := reader.Close(); err != nil {
			logger.Errorw("error closing kafka reader", "error", err)
		}
	}()
}

func (c *consumer) handleUpdateGamePlayerDataMessage(ctx context.Context, _ *kafka.Message, uncast proto.Message) {
	msg := uncast.(*pbmsg.UpdateGamePlayerDataMessage)

	pId, err := uuid.Parse(msg.PlayerId)
	if err != nil {
		c.logger.Errorw("failed to parse player id", "error", err)
		return
	}

	switch msg.GameMode {
	case pbmodel.GameDataGameMode_BLOCK_SUMO:
		err = c.handleBlockSumoUpdate(ctx, pId, msg)
	}

	if err != nil {
		c.logger.Errorw("failed to handle update", "error", err, "playerId", pId, "gameMode", msg.GameMode)
		return
	}
}

func (c *consumer) handleBlockSumoUpdate(ctx context.Context, pId uuid.UUID, msg *pbmsg.UpdateGamePlayerDataMessage) error {
	player, err := c.repo.GetBlockSumoData(ctx, pId)

	if err != nil {
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("failed to get block sumo data: %w", err)
		}

		player = config.CreateDefaultBlockSumoData(pId)
	}

	msgData := &pbmodel.V1BlockSumoPlayerData{}

	if err := anypb.UnmarshalTo(msg.Data, msgData, proto.UnmarshalOptions{}); err != nil {
		return fmt.Errorf("failed to unmarshal data: %w", err)
	}

	for _, path := range msg.DataMask.Paths {
		switch path {
		case "block_slot":
			player.BlockSlot = msgData.BlockSlot
		case "shears_slot":
			player.ShearsSlot = msgData.ShearsSlot
		}
	}

	return c.repo.SaveBlockSumoPlayer(ctx, player)
}
