package app

import (
	"context"
	"game-player-data/internal/config"
	"game-player-data/internal/repository"
	"game-player-data/internal/repository/model"
	"game-player-data/internal/service"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"os/signal"
	"sync"
	"syscall"
)

func Run(cfg *config.Config, logger *zap.SugaredLogger) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	wg := &sync.WaitGroup{}

	repoCtx, repoCancel := context.WithCancel(ctx)
	repoWg := &sync.WaitGroup{}

	repo, err := repository.NewMongoRepository(repoCtx, logger, repoWg, cfg.MongoDB)
	if err != nil {
		logger.Fatalw("failed to create repository", err)
	}

	if err := repo.SaveBlockSumoPlayer(ctx, &model.BlockSumoData{
		PlayerId:   uuid.MustParse("8d36737e-1c0a-4a71-87de-9906f577845e"),
		BlockSlot:  1,
		ShearsSlot: 2,
	}); err != nil {
		panic(err)
	}

	// todo kafka consumer

	service.RunServices(ctx, logger, wg, cfg, repo)

	wg.Wait()
	logger.Info("shutting down")

	logger.Info("shutting down repository")
	repoCancel()
	repoWg.Wait()
}
