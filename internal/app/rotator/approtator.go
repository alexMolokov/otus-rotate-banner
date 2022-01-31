package approtator

import "context"

type App struct {
	Logger  Logger
	Storage Storage
}

func (a *App) AddBannerToSlot(ctx context.Context, bannerID, slotID int64) error {
	return nil
}

func (a *App) RemoveBannerFromSlot(ctx context.Context, bannerID, slotID int64) error {
	return nil
}

func (a *App) CountTransition(ctx context.Context, bannerID, slotID, sgID int64) error {
	return nil
}

func (a *App) ChooseBanner(ctx context.Context, slotID, sgID int64) (bannerID int64, err error) {
	return 1, nil
}

type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warning(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

type Storage interface{}

func NewAppRotator(logger Logger, storage Storage) *App {
	return &App{
		Logger:  logger,
		Storage: storage,
	}
}
