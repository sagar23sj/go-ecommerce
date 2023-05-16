package repository

import (
	"context"
	"log"

	"gorm.io/gorm"
)

type BaseRepository struct {
	DB *gorm.DB
}

// RepositoryTransaction interface holds transaction specific methods
type RepositoryTransaction interface {
	// return a transaction from a gorm connection
	BeginTx(ctx context.Context) (tx *gorm.DB, err error)
	CommitTx(tx *gorm.DB) (err error)
	RollbackTx(tx *gorm.DB) (err error)
	HandleTransaction(ctx context.Context, tx *gorm.DB, incomingErr error) (err error)
	// get a db connection
	GetConn() (conn *gorm.DB)
}

func (repo *BaseRepository) BeginTx(ctx context.Context) (tx *gorm.DB, err error) {

	sqlDB := repo.DB.Begin()
	if sqlDB.Error != nil {
		log.Printf("error occured while initiating database transaction: %v", err.Error())
		return nil, err
	}

	return sqlDB, nil
}

func (repo *BaseRepository) GetConn() (conn *gorm.DB) {
	return repo.DB
}

func (repo *BaseRepository) CommitTx(tx *gorm.DB) (err error) {
	err = tx.Commit().Error
	return
}

func (repo *BaseRepository) RollbackTx(tx *gorm.DB) (err error) {
	err = tx.Rollback().Error
	return
}

func (repo *BaseRepository) HandleTransaction(ctx context.Context, tx *gorm.DB, incomingErr error) (err error) {
	if incomingErr != nil {
		err = tx.Rollback().Error
		if err != nil {
			return
		}
		return
	}

	err = tx.Commit().Error
	if err != nil {
		return
	}
	return
}

func (repo *BaseRepository) initiateQueryExecutor(tx *gorm.DB) *gorm.DB {
	//Populate the query executor so we can join/use a transaction if one is present.
	//If we are not running inside a transaction then the plain *gorm.DB object is used.
	executor := repo.DB
	if tx != nil {
		executor = tx
	}
	return executor
}
