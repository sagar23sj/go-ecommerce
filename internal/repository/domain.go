package repository

import "context"

// RepositoryTransaction interface holds transaction specific methods
type RepositoryTransaction interface {
	// return a transaction from a gorm connection
	BeginTx(ctx context.Context) (Transaction, error)
	HandleTransaction(ctx context.Context, tx Transaction, incomingErr error) error
}

type Transaction interface {
	Commit() error
	Rollback() error
}
