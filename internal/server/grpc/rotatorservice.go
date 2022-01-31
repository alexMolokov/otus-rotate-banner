package internalgrpc

import (
	"context"

	"github.com/alexMolokov/otus-rotate-banner/api/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type RotatorService struct {
	app    Application
	logger Logger

	pb.UnimplementedRotatorServer
}

func (rs *RotatorService) AddBannerToSlot(ctx context.Context, r *pb.AddBannerToSlotRequest) (*emptypb.Empty, error) {
	err := rs.app.AddBannerToSlot(ctx, r.BannerId, r.SlotId)
	if err != nil {
		rs.logger.Error("can't add banner to slot bannerId=%d slotId=%d %v", r.BannerId, r.SlotId, err)
		return nil, status.Errorf(codes.Internal, "can't add banner to slot")
	}
	return nil, nil
}

func (rs *RotatorService) RemoveBannerFromSlot(ctx context.Context, r *pb.RemoveBannerFromSlotRequest) (*emptypb.Empty, error) { //nolint
	err := rs.app.RemoveBannerFromSlot(ctx, r.BannerId, r.SlotId)
	if err != nil {
		rs.logger.Error("can't remove banner from slot bannerId=%d slotId=%d %v", r.BannerId, r.SlotId, err)
		return nil, status.Errorf(codes.Internal, "can't remove banner from slot")
	}

	return nil, nil
}

func (rs *RotatorService) CountTransition(ctx context.Context, r *pb.CountTransitionRequest) (*emptypb.Empty, error) {
	err := rs.app.CountTransition(ctx, r.BannerId, r.SlotId, r.SgId)
	if err != nil {
		rs.logger.Error("can't count the transition bannerId=%d slotId=%d sgId=%d %v", r.BannerId, r.SlotId, r.SgId, err)
		return nil, status.Errorf(codes.Internal, "can't count the transition")
	}

	return nil, nil
}

func (rs *RotatorService) ChooseBanner(ctx context.Context, r *pb.ChooseBannerRequest) (*pb.ChooseBannerResponse, error) { //nolint
	bannerID, err := rs.app.ChooseBanner(ctx, r.SlotId, r.SgId)
	if err != nil {
		rs.logger.Error("can't choose banner for slotId=%d sgId=%d %v", r.SlotId, r.SgId, err)
		return nil, status.Errorf(codes.Internal, "can't choose banner")
	}

	return &pb.ChooseBannerResponse{BannerId: bannerID}, nil
}
