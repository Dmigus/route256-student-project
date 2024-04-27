package multipostgres

import (
	"context"
	"errors"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/rs/xid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
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
	connWithTr struct {
		conn TxBeginner
		tr   pgx.Tx
	}
	distributedTransaction struct {
		mu                 sync.Mutex
		openedTransactions map[TxBeginner]pgx.Tx
		preparedForCommit  map[pgx.Tx]string
	}
)

func newDistributedTransaction() *distributedTransaction {
	return &distributedTransaction{
		openedTransactions: make(map[TxBeginner]pgx.Tx),
		preparedForCommit:  make(map[pgx.Tx]string),
	}
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
	err = tx.commitPrepared(ctx)
	if err == nil {
		span.AddEvent("distributed transaction committed")
		span.SetStatus(codes.Ok, "")
	}
	return err
}

func (tx *distributedTransaction) GetTransaction(ctx context.Context, beginner TxBeginner) (pgx.Tx, error) {
	tx.mu.Lock()
	defer tx.mu.Unlock()
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
	return tx.forEachTransaction(ctx, func(ctx context.Context, _ TxBeginner, transaction pgx.Tx) error {
		tx.mu.Lock()
		transID, prepared := tx.preparedForCommit[transaction]
		tx.mu.Unlock()
		var err error
		if prepared {
			alwaysActualCtx := context.WithoutCancel(ctx)
			_, err = transaction.Exec(alwaysActualCtx, "ROLLBACK PREPARED '"+transID+"'")
			// release connection to pool
			_ = transaction.Rollback(alwaysActualCtx)
		} else {
			err = transaction.Rollback(ctx)
		}
		return err
	})
}

// prepareForCommit переводит все транзакции в prepared. Оставляет коннекты
func (tx *distributedTransaction) prepareForCommit(ctx context.Context) error {
	return tx.forEachTransaction(ctx, func(ctx context.Context, _ TxBeginner, transaction pgx.Tx) error {
		transID := tx.createTransactionID()
		_, err := transaction.Exec(ctx, "PREPARE TRANSACTION '"+transID+"'")
		if err == nil {
			tx.mu.Lock()
			tx.preparedForCommit[transaction] = transID
			tx.mu.Unlock()
		}
		return err
	})
}

func (tx *distributedTransaction) createTransactionID() string {
	guid := xid.New()
	return guid.String()
}

func (tx *distributedTransaction) commitPrepared(ctx context.Context) error {
	return tx.forEachTransaction(ctx, func(ctx context.Context, conn TxBeginner, transaction pgx.Tx) error {
		transID := tx.preparedForCommit[transaction]
		_, err := transaction.Exec(ctx, "COMMIT PREPARED '"+transID+"'")
		if err == nil {
			tx.excludeTransaction(conn)
			// release connection to pool
			_ = transaction.Commit(ctx)
		}
		return err
	})
}

func (tx *distributedTransaction) excludeTransaction(conn TxBeginner) {
	tx.mu.Lock()
	defer tx.mu.Unlock()
	tr, ok := tx.openedTransactions[conn]
	if ok {
		delete(tx.openedTransactions, conn)
		delete(tx.preparedForCommit, tr)
	}
}

func (tx *distributedTransaction) forEachTransaction(ctx context.Context, f func(context.Context, TxBeginner, pgx.Tx) error) error {
	transactions := tx.transactionsAsSlice()
	errs := make([]error, len(transactions))
	wg := sync.WaitGroup{}
	wg.Add(len(transactions))
	for ind, connectionWithTransaction := range transactions {
		go func() {
			defer wg.Done()
			errs[ind] = f(ctx, connectionWithTransaction.conn, connectionWithTransaction.tr)
		}()
	}
	wg.Wait()
	return errors.Join(errs...)
}

func (tx *distributedTransaction) transactionsAsSlice() []connWithTr {
	tx.mu.Lock()
	defer tx.mu.Unlock()
	trs := make([]connWithTr, 0, len(tx.openedTransactions))
	for conn, tr := range tx.openedTransactions {
		trs = append(trs, connWithTr{conn: conn, tr: tr})
	}
	return trs
}
