package db

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/ibraheemacara/tezos-delegation-service/config"
	"github.com/jackc/pgx/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	log "github.com/sirupsen/logrus"
)

type DbStore struct {
	DB      *gorm.DB
	pgxConn *pgx.Conn
}

type DBInterface interface {
	GetDelegations() ([]Delegations, error)
	GetDelegationsByYear(year string) ([]Delegations, error)
	GetLastBlock() (int32, error)
	InsertDelegations(delegator string, timestamp time.Time, block int32, amount int64) error
	BulkInsertDelegations(delegations []Delegations) error
}

func InitDB(cfg config.Config) (DBInterface, error) {
	dbPwd := url.QueryEscape(cfg.Db.Password)
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.Db.User, dbPwd, cfg.Db.Host, cfg.Db.Port, cfg.Db.Database)
	gormDB, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	pgxConn, err := pgx.Connect(context.Background(), connectionString)
	if err != nil {
		return nil, err
	}

	dbStore := &DbStore{
		DB:      gormDB,
		pgxConn: pgxConn,
	}

	if err := dbStore.DB.AutoMigrate(&Delegations{}); err != nil {
		return nil, err
	}

	log.Info("Database initialized successfully")

	return dbStore, nil
}

func (db *DbStore) GetDelegations() ([]Delegations, error) {
	var delegations []Delegations
	if err := db.DB.Order("block DESC").Limit(50).Find(&delegations).Error; err != nil {
		return nil, err
	}
	return delegations, nil
}

func (db *DbStore) GetDelegationsByYear(year string) ([]Delegations, error) {
	var delegations []Delegations
	if err := db.DB.Where("EXTRACT(YEAR FROM timestamp) = ?", year).Order("block DESC").Limit(50).Find(&delegations).Error; err != nil {
		return nil, err
	}
	return delegations, nil
}

func (db *DbStore) InsertDelegations(delegator string, timestamp time.Time, block int32, amount int64) error {
	delegation := Delegations{
		Delegator: delegator,
		Timestamp: timestamp,
		Block:     block,
		Amount:    amount,
	}
	return db.DB.Create(&delegation).Error
}

func (db *DbStore) GetLastBlock() (int32, error) {
	var delegation Delegations
	err := db.DB.Order("block DESC").Limit(1).Find(&delegation).Error
	if err == gorm.ErrRecordNotFound {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return delegation.Block, nil
}

func (db *DbStore) BulkInsertDelegations(delegations []Delegations) error {

	ctx := context.Background()
	defer db.pgxConn.Close(ctx)

	copyCount, err := db.pgxConn.CopyFrom(
		ctx,
		pgx.Identifier{"delegations"},
		[]string{"delegator", "timestamp", "block", "amount"},
		pgx.CopyFromSlice(len(delegations), func(i int) ([]interface{}, error) {
			delegation := delegations[i]
			return []interface{}{
				delegation.Delegator,
				delegation.Timestamp,
				delegation.Block,
				delegation.Amount,
			}, nil
		}),
	)

	if err != nil {
		return err
	}
	log.Infof("Copied %d delegations to database", copyCount)
	return nil
}
