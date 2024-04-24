package multipostgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/rs/xid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"sync"
)

var tracer = otel.Tracer("postgres distributed transaction")

type (
	// TxBeginner объект, умеющий открывать новую транзакцию
	TxBeginner interface {
		Begin(ctx context.Context) (pgx.Tx, error)
	}
	command interface {
		execute(ctx context.Context) (hasToBeCommitted bool)
	}
	distributedTransaction struct {
		mu                 sync.Mutex
		openedTransactions map[TxBeginner]pgx.Tx
		preparedForCommit  map[pgx.Tx]string
	}
)

func newDistributedTransaction() *distributedTransaction {
	return &distributedTransaction{openedTransactions: make(map[TxBeginner]pgx.Tx)}
}

func (tx *distributedTransaction) WithinTransaction(ctx context.Context, command command) (err error) {
	ctx, span := tracer.Start(ctx, "distributed transaction")
	defer span.End()
	span.SetStatus(codes.Error, "transaction was not committed")
	defer func() {
		errRollback := tx.rollback(ctx)
		if err == nil {
			err = errRollback
		}
	}()
	if !command.execute(ctx) {
		return nil
	}
	err = tx.prepareForCommit(ctx)
	if err != nil {
		return err
	}
	err = tx.commit(ctx)
	if err == nil {
		span.AddEvent("distributed transaction committed")
		span.SetStatus(codes.Ok, "")
	}
	return err
}

func (tx *distributedTransaction) GetTransaction(ctx context.Context, beginner TxBeginner) (pgx.Tx, error) {
	if transaction, ok := tx.openedTransactions[beginner]; ok {
		return transaction, nil
	}
	transaction, err := beginner.Begin(ctx)
	if err != nil {
		return nil, err
	}
	tx.openedTransactions[beginner] = transaction
	return transaction, nil
}

func (tx *distributedTransaction) rollback(ctx context.Context) error {
	return tx.forEachTransaction(ctx, func(ctx context.Context, transaction pgx.Tx) error {
		tx.mu.Lock()
		transId, prepared := tx.preparedForCommit[transaction]
		tx.mu.Unlock()
		var err error
		if prepared {
			_, err = transaction.Exec(ctx, "ROLLBACK PREPARED $1;", transId)
		} else {
			err = transaction.Rollback(ctx)
		}
		if errors.Is(err, pgx.ErrTxClosed) {
			return nil
		}
		return err
	})
}

// переводит все транзакции в prepared
func (tx *distributedTransaction) prepareForCommit(ctx context.Context) error {
	return tx.forEachTransaction(ctx, func(ctx context.Context, transaction pgx.Tx) error {
		transId := tx.createTransactionID()
		_, err := transaction.Exec(ctx, "PREPARE TRANSACTION $1;", transId)
		if err == nil {
			tx.mu.Lock()
			tx.preparedForCommit[transaction] = transId
			tx.mu.Unlock()
		}
		return err
	})
}

func (tx *distributedTransaction) createTransactionID() string {
	guid := xid.New()
	return guid.String()
}

func (tx *distributedTransaction) commit(ctx context.Context) error {
	return tx.forEachTransaction(ctx, func(ctx context.Context, transaction pgx.Tx) error {
		transId, _ := tx.preparedForCommit[transaction]
		_, err := transaction.Exec(ctx, "COMMIT PREPARED $1;", transId)
		return err
	})
}

func (tx *distributedTransaction) forEachTransaction(ctx context.Context, f func(context.Context, pgx.Tx) error) error {
	transactions := tx.transactionsAsSlice()
	errs := make([]error, len(transactions))
	wg := sync.WaitGroup{}
	wg.Add(len(transactions))
	for ind, transaction := range transactions {
		go func() {
			defer wg.Done()
			err := f(ctx, transaction)
			errs[ind] = err
		}()
	}
	wg.Wait()
	return errors.Join(errs...)
}

func (tx *distributedTransaction) transactionsAsSlice() []pgx.Tx {
	tx.mu.Lock()
	defer tx.mu.Unlock()
	trs := make([]pgx.Tx, 0, len(tx.openedTransactions))
	for _, tx := range tx.openedTransactions {
		trs = append(trs, tx)
	}
	return trs
}
