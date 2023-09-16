package service

import (
	"context"
	"errors"
	"game-player-data/internal/repository"
	"game-player-data/internal/repository/model"
	pb "github.com/emortalmc/proto-specs/gen/go/grpc/gameplayerdata"
	"github.com/emortalmc/proto-specs/gen/go/model/gameplayerdata"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type gamePlayerDataService struct {
	pb.UnimplementedGamePlayerDataServiceServer

	repo repository.Repository
}

func newGamePlayerDataService(repo repository.Repository) pb.GamePlayerDataServiceServer {
	return &gamePlayerDataService{
		repo: repo,
	}
}

func (s *gamePlayerDataService) GetGamePlayerData(ctx context.Context, req *pb.GetGamePlayerDataRequest) (*pb.GetGamePlayerDataResponse, error) {
	pId, err := uuid.Parse(req.PlayerId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid player id")
	}

	var data model.GameData

	switch req.GameMode {
	case gameplayerdata.GameDataGameMode_BLOCK_SUMO:
		data, err = s.repo.GetBlockSumoData(ctx, pId)
	}

	if err != nil {
		return nil, createDbErr(err)
	}

	anyData, err := data.ToAnyProto()
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to convert data to proto")
	}

	return &pb.GetGamePlayerDataResponse{
		Data: anyData,
	}, nil
}

func createDbErr(err error) error {
	if errors.Is(err, mongo.ErrNoDocuments) {
		return status.Error(codes.NotFound, "player not found")
	} else {
		return status.Error(codes.Internal, "failed to get player data")
	}
}
