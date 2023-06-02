package service

import (
	"context"
	"game-player-data/internal/repository"
	pb "github.com/emortalmc/proto-specs/gen/go/grpc/gameplayerdata"
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

func (s *gamePlayerDataService) GetBlockSumoPlayerData(ctx context.Context, req *pb.PlayerIdRequest) (*pb.GetBlockSumoPlayerDataResponse, error) {
	pId, err := uuid.Parse(req.PlayerId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid player id")
	}

	data, err := s.repo.GetBlockSumoPlayer(ctx, pId)
	if err != nil {
		return nil, createDbErr(err)
	}

	return &pb.GetBlockSumoPlayerDataResponse{PlayerData: data.ToProto()}, nil
}

func (s *gamePlayerDataService) GetTowerDefencePlayerData(ctx context.Context, req *pb.PlayerIdRequest) (*pb.GetTowerDefencePlayerDataResponse, error) {
	pId, err := uuid.Parse(req.PlayerId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid player id")
	}

	data, err := s.repo.GetTowerDefencePlayer(ctx, pId)
	if err != nil {
		return nil, createDbErr(err)
	}

	return &pb.GetTowerDefencePlayerDataResponse{PlayerData: data.ToProto()}, nil
}

func (s *gamePlayerDataService) GetMinesweeperPlayerData(ctx context.Context, req *pb.PlayerIdRequest) (*pb.GetMinesweeperPlayerDataResponse, error) {
	pId, err := uuid.Parse(req.PlayerId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid player id")
	}

	data, err := s.repo.GetMinesweeperPlayer(ctx, pId)
	if err != nil {
		return nil, createDbErr(err)
	}

	return &pb.GetMinesweeperPlayerDataResponse{PlayerData: data.ToProto()}, nil
}

func createDbErr(err error) error {
	if err == mongo.ErrNoDocuments {
		return status.Error(codes.NotFound, "player not found")
	} else {
		return status.Error(codes.Internal, "failed to get player data")
	}
}
