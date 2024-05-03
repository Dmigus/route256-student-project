package multipostgres

import (
	"context"
	"errors"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/rs/xid"
	lop "github.com/samber/lo/parallel"
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
	// distributedTransactionCoordinator координирует работу нескольких транзакций, используя двухфахный коммит
	distributedTransactionCoordinator struct {
		mu                 sync.Mutex
		openedTransactions map[TxBeginner]pgx.Tx
		preparedForCommit  map[pgx.Tx]string
	}
)

// newDistributedTransactionCoordinator возвращает новый distributedTransactionCoordinator
func newDistributedTransactionCoordinator() *distributedTransactionCoordinator {
	return &distributedTransactionCoordinator{
		openedTransactions: make(map[TxBeginner]pgx.Tx),
		preparedForCommit:  make(map[pgx.Tx]string),
	}
}

// WithinTransaction исполняет команду command в рамках распределённой транзакции
func (tx *distributedTransactionCoordinator) WithinTransaction(ctx context.Context, command command) (err error) {
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

// GetTransaction возвращает имеющуюся транзакцию, полученную от beginner или создаёт новую и возвращает её, если её не было
func (tx *distributedTransactionCoordinator) GetTransaction(ctx context.Context, beginner TxBeginner) (pgx.Tx, error) {
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

// rollback производит откат всех открытых транзакций
func (tx *distributedTransactionCoordinator) rollback(ctx context.Context) error {
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

// prepareForCommit переводит все транзакции в prepared. Не завершает открытые коннекты
func (tx *distributedTransactionCoordinator) prepareForCommit(ctx context.Context) error {
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

func (tx *distributedTransactionCoordinator) createTransactionID() string {
	guid := xid.New()
	return guid.String()
}

// commitPrepared выполняет коммит для подготовленных транзакций
func (tx *distributedTransactionCoordinator) commitPrepared(ctx context.Context) error {
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

// excludeTransaction удаляет информацию о транзакции из координатора
func (tx *distributedTransactionCoordinator) excludeTransaction(conn TxBeginner) {
	tx.mu.Lock()
	defer tx.mu.Unlock()
	tr, ok := tx.openedTransactions[conn]
	if ok {
		delete(tx.openedTransactions, conn)
		delete(tx.preparedForCommit, tr)
	}
}

func (tx *distributedTransactionCoordinator) forEachTransaction(ctx context.Context, f func(context.Context, TxBeginner, pgx.Tx) error) error {
	transactions := tx.transactionsAsSlice()
	errs := lop.Map(transactions, func(item connWithTr, _ int) error {
		return f(ctx, item.conn, item.tr)
	})
	return errors.Join(errs...)
}

func (tx *distributedTransactionCoordinator) transactionsAsSlice() []connWithTr {
	tx.mu.Lock()
	defer tx.mu.Unlock()
	trs := make([]connWithTr, 0, len(tx.openedTransactions))
	for conn, tr := range tx.openedTransactions {
		trs = append(trs, connWithTr{conn: conn, tr: tr})
	}
	return trs
}
