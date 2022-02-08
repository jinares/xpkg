package xmysql

import (
	"context"
	"github.com/doug-martin/goqu/v9"
	"github.com/jinares/xpkg/xlog"
	"github.com/jinares/xpkg/xtools"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"time"
)

type (
	GoquTranscationSQLHandler func(tx *goqu.TxDatabase) error
	//TranscationSQLHandler     func(tx *goqu.TxDatabase) (commmit bool, err error)
	GoquSQLHandler func(tx *goqu.Database) error
)

func GoquSQL(ctx context.Context, fn DBHandler, handler GoquSQLHandler) (err error) {
	if fn == nil {
		return xtools.XErr(codes.Internal, "db is nil")
	}
	db, err := fn()
	if err != nil {
		return err
	}
	godb := goqu.New("mysql", db)

	ulog := xlog.CtxLog(ctx).WithFields(logrus.Fields{
		"type":    "goqu",
		"goqu_id": xtools.GUID(),
	})
	if IsEnable() {
		godb.Logger(ulog)
	}
	st := time.Now()
	defer func() {

		exec_time := time.Since(st).Milliseconds()
		if exec_time > 1000 {
			if err != nil {
				ulog = xlog.Err(err, ulog)
			}
			ulog.WithField("exec_time", exec_time).Warn("mysql-run-slow")
		}
	}()
	return handler(godb)
}

func GoquTransactionSQL(ctx context.Context, fn DBHandler, handler GoquTranscationSQLHandler) error {
	nlog := xlog.CtxLog(ctx).WithFields(logrus.Fields{
		"type":    "goqu-tx",
		"goqu_id": xtools.GUID(),
	})
	if fn == nil {
		nlog.Error("mysql-get-empty")
		return xtools.XErr(codes.Internal, "db is nil")
	}

	gdb, err := fn()
	if err != nil {
		xlog.Err(err, nlog).Error("mysql-get-conn")
		return xtools.InternalErr(err, "请稍后再试", true)
	}
	db := goqu.New("mysql", gdb)
	tx, err := db.Begin()
	if err != nil {
		xlog.Err(err, nlog).Error("mysql-start-transaction")
		return xtools.InternalErr(err, "请稍后再试", true)
	}
	if IsEnable() {
		tx.Logger(nlog)
	}
	st := time.Now()
	err = handler(tx)

	exec_time := time.Since(st).Milliseconds()
	nlog = nlog.WithField("exec_time", exec_time)
	if exec_time > 1000 {
		nlog.Warn("mysql-run-slow")
	}
	if err != nil {
		xlog.Err(err, nlog).Error("mysql-rollback")
		rerr := tx.Rollback()
		if rerr != nil {
			nlog.WithField("rollback_fail", xtools.ErrString(rerr)).Error("mysql-rollback")
		}
		return err
	}
	err = tx.Commit()
	if err != nil {
		xlog.Err(err, nlog).Error("mysql-commit")
		return xtools.InternalErr(err, "请稍后再试", true)
	}
	return nil
}
