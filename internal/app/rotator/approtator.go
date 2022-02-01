package approtator

import (
	"context"
	"time"

	"github.com/alexMolokov/otus-rotate-banner/internal/algorithm/bandit"
	rotatorstorage "github.com/alexMolokov/otus-rotate-banner/internal/storage/rotator"
)

type App struct {
	Logger  Logger
	Storage Storage
}

func (a *App) AddBannerToSlot(ctx context.Context, bannerID, slotID int64) error {
	opCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	return a.Storage.AddBannerToSlot(opCtx, bannerID, slotID)
}

func (a *App) RemoveBannerFromSlot(ctx context.Context, bannerID, slotID int64) error {
	opCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	return a.Storage.RemoveBannerFromSlot(opCtx, bannerID, slotID)
}

func (a *App) CountTransition(ctx context.Context, bannerID, slotID, sgID int64) error {
	opCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	return a.Storage.CountTransition(opCtx, bannerID, slotID, sgID)
}

func (a *App) ChooseBanner(ctx context.Context, slotID, sgID int64) (bannerID int64, err error) {
	opCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	bannerStat, totalDisplay, err := a.Storage.GetBannersStat(opCtx, slotID, sgID)
	if err != nil {
		return 0, err
	}

	stat := make([]bandit.Stat, len(bannerStat))
	for i, v := range bannerStat {
		stat[i] = bandit.Stat{
			ID:     int(v.ID),
			Trials: int(v.Display),
			Reward: int(v.Click),
		}
	}

	bannerID = int64(bandit.Choice(stat, totalDisplay))
	err = a.Storage.CountTransition(ctx, bannerID, slotID, sgID)
	if err != nil {
		return 0, err
	}

	return bannerID, nil
}

type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warning(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

type Storage interface {
	GetBannerByID(ctx context.Context, bannerID int64) (*rotatorstorage.Banner, error)
	GetSlotByID(ctx context.Context, slotID int64) (*rotatorstorage.Slot, error)
	GetSocialGroupByID(ctx context.Context, sgID int64) (*rotatorstorage.SocialGroup, error)
	AddBannerToSlot(ctx context.Context, bannerID, slotID int64) error
	RemoveBannerFromSlot(ctx context.Context, bannerID, slotID int64) error
	CountTransition(ctx context.Context, bannerID, slotID, sgID int64) error
	CountDisplay(ctx context.Context, bannerID, slotID, sgID int64) error
	GetBannersStat(ctx context.Context, slotID, sgID int64) ([]rotatorstorage.BannerStat, int, error)
}

func NewAppRotator(logger Logger, storage Storage) *App {
	return &App{
		Logger:  logger,
		Storage: storage,
	}
}
