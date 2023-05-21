package repository

import (
	"context"
	"log"
	"time"

	"github.com/asdine/storm/v3"
	"github.com/sagar23sj/go-ecommerce/internal/repository"
)

type BaseRepository struct {
	DB *storm.DB
}

type BaseTransaction struct {
	tx storm.Node
}

func (repo *BaseRepository) TimeNow() time.Time {
	return time.Now()
}

func (repo *BaseRepository) BeginTx(ctx context.Context) (repository.Transaction, error) {

	txObj, err := repo.DB.Begin(true)
	if err != nil {
		log.Printf("error occured while initiating database transaction: %v", err.Error())
		return nil, err
	}

	return &BaseTransaction{
		tx: txObj,
	}, nil
}

func (repo *BaseRepository) HandleTransaction(ctx context.Context, tx repository.Transaction, incomingErr error) (err error) {
	if incomingErr != nil {
		err = tx.Rollback()
		if err != nil {
			log.Printf("error occured while rollback database transaction: %v", err.Error())
			return
		}
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("error occured while commit database transaction: %v", err.Error())
		return
	}
	return
}

func (repo *BaseTransaction) Commit() error {
	return repo.tx.Commit()
}

func (repo *BaseTransaction) Rollback() error {
	return repo.tx.Rollback()
}

func (repo *BaseRepository) initiateQueryExecutor(tx repository.Transaction) (executor storm.Node) {
	//Populate the query executor so we can use a transaction if one is present.
	//If we are not running inside a transaction then the plain *storm.DB object is used.
	executor = repo.DB
	if tx != nil {
		txObj := tx.(*BaseTransaction)
		executor = txObj.tx
	}

	return executor
}
