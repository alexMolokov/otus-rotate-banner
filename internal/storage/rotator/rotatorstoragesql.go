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
	ErrConnectDB            = errors.New("can't connect to database")
	ErrBannerNotExists      = errors.New("banner not exists")
	ErrSlotNotExists        = errors.New("slot not exists")
	ErrBannerNotInSlot      = errors.New("banner not in slot")
	ErrSocialGroupNotExists = errors.New("social group not exists")
	ErrCountTransition      = errors.New("can't register transition")
	ErrDisplayTransition    = errors.New("can't register transition")
)

type Banner struct {
	ID          int64  `db:"banner_id"`
	Description string `db:"description"`
}

type Slot struct {
	ID          int64  `db:"slot_id"`
	Description string `db:"description"`
}

type SocialGroup struct {
	ID          int64  `db:"social_group_id"`
	Description string `db:"description"`
}

type BannerStat struct {
	ID      int64 `db:"banner_id"`
	Display int64 `db:"display"`
	Click   int64 `db:"click"`
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

func (s *Storage) GetSocialGroupByID(ctx context.Context, sgID int64) (*SocialGroup, error) {
	row := s.db.QueryRowxContext(ctx,
		`SELECT social_group_id, description 
		FROM social_group WHERE social_group_id = $1`, sgID)
	if err := row.Err(); err != nil {
		return nil, fmt.Errorf("can't get social group %d: %w", sgID, err)
	}

	var sg SocialGroup
	if err := row.StructScan(&sg); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSocialGroupNotExists
		}
		return nil, fmt.Errorf("can't get social group id %d row scan : %w", sgID, err)
	}
	return &sg, nil
}

// AddBannerToSlot Добавляет связку баннер <-> слот.
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

// RemoveBannerFromSlot Удаляет связку баннер <-> слот.
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

// CountTransition Регистрирует переход (клик на баннере).
func (s *Storage) CountTransition(ctx context.Context, bannerID, slotID, sgID int64) error {
	query := `UPDATE stat SET click = click + 1
		WHERE slot_id = $1 AND banner_id = $2 AND sg_id = $3`

	result, err := s.db.ExecContext(ctx, query, bannerID, slotID, sgID)
	if err != nil {
		return fmt.Errorf("can't count transition slot %d banner = %d social group %d: %w", slotID, bannerID, sgID, err)
	}

	count, err := result.RowsAffected()
	if err != nil || count == 0 {
		return ErrCountTransition
	}

	return nil
}

// CountDisplay Регистрирует показ баннера.
func (s *Storage) CountDisplay(ctx context.Context, bannerID, slotID, sgID int64) error {
	query := `INSERT INTO stat (slot_id, banner_id, social_group_id)
		VALUES ($1, $2, $3)
		ON CONFLICT (slot_id, banner_id, social_group_id) 
    	DO UPDATE SET display = excluded.display + 1`

	result, err := s.db.ExecContext(ctx, query, bannerID, slotID, sgID)
	if err != nil {
		return fmt.Errorf("can't count display slot %d banner = %d social group %d: %w", slotID, bannerID, sgID, err)
	}

	count, err := result.RowsAffected()
	if err != nil || count == 0 {
		return ErrDisplayTransition
	}

	return nil
}

// GetBannersStat Выбирает баннеры с их статистиками
// которые могут быть показаны в указанном слоте и для указанной соц.группы.
func (s *Storage) GetBannersStat(ctx context.Context, slotID, sgID int64) ([]BannerStat, error) {
	result := make([]BannerStat, 0)
	query := `SELECT banner_id, display, click FROM stat
		WHERE slot_id = $1 AND social_group_id = $2`
	err := s.db.SelectContext(ctx, &result, query, slotID, sgID)
	return result, err
}
