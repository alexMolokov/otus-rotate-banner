package rotatorstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/alexMolokov/otus-rotate-banner/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" //nolint
)

var (
	ErrConnectDB       = errors.New("can't connect to database")
	ErrBannerNotExists = errors.New("banner not exists")
	ErrSlotNotExists   = errors.New("slot not exists")
	ErrBannerNotInSlot = errors.New("banner not in slot")
)

type Banner struct {
	ID          int64  `db:"banner_id"`
	Description string `db:"description"`
}

type Slot struct {
	ID          int64  `db:"slot_id"`
	Description string `db:"description"`
}

type Storage struct {
	db  *sqlx.DB
	cfg config.DBConf
}

func NewRotatorStorage(cfg config.DBConf) *Storage {
	return &Storage{
		cfg: cfg,
	}
}

func (s *Storage) Connect() error {
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=%s",
		s.cfg.User, s.cfg.Password, s.cfg.Name, s.cfg.Host, s.cfg.Port, s.cfg.SslMode)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return fmt.Errorf("%s %w", ErrConnectDB, err)
	}

	s.db = db
	if s.cfg.MaxConnectionPool > 0 {
		s.db.SetMaxOpenConns(s.cfg.MaxConnectionPool)
	}

	return nil
}

func (s *Storage) Close() error {
	if s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *Storage) GetBannerByID(ctx context.Context, bannerID int64) (*Banner, error) {
	row := s.db.QueryRowxContext(ctx, "SELECT banner_id, description FROM banner WHERE banner_id = $1", bannerID)
	if err := row.Err(); err != nil {
		return nil, fmt.Errorf("can't get banner %d: %w", bannerID, err)
	}

	var b Banner
	if err := row.StructScan(&b); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrBannerNotExists
		}
		return nil, fmt.Errorf("can't get banner id %d row scan : %w", bannerID, err)
	}
	return &b, nil
}

func (s *Storage) GetSlotByID(ctx context.Context, slotID int64) (*Slot, error) {
	row := s.db.QueryRowxContext(ctx, "SELECT slot_id, description FROM slot WHERE slot_id = $1", slotID)
	if err := row.Err(); err != nil {
		return nil, fmt.Errorf("can't get slot %d: %w", slotID, err)
	}

	var sl Slot
	if err := row.StructScan(&sl); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSlotNotExists
		}
		return nil, fmt.Errorf("can't get slot id %d row scan : %w", slotID, err)
	}
	return &sl, nil
}

func (s *Storage) AddBannerToSlot(ctx context.Context, bannerID, slotID int64) error {
	_, err := s.GetSlotByID(ctx, slotID)
	if err != nil {
		return err
	}

	_, err = s.GetBannerByID(ctx, bannerID)
	if err != nil {
		return err
	}

	query := `INSERT INTO banner_to_slot
	(banner_id, slot_id)
	VALUES ($1, $2)
	ON CONFLICT DO NOTHING`

	_, err = s.db.ExecContext(ctx, query, bannerID, slotID)
	if err != nil {
		return fmt.Errorf("can't bind  banner %d, slot : %d, %w", bannerID, slotID, err)
	}

	return nil
}

func (s *Storage) RemoveBannerFromSlot(ctx context.Context, bannerID, slotID int64) error {
	result, err := s.db.ExecContext(
		ctx,
		"DELETE FROM banner_to_slot WHERE banner_id = $1 AND slot_id = $2",
		bannerID, slotID,
	)
	if err != nil {
		return fmt.Errorf("can't delete banner %d from slot = %d %w", bannerID, slotID, err)
	}
	count, _ := result.RowsAffected()
	if count == 0 {
		return ErrBannerNotInSlot
	}
	return nil
}
