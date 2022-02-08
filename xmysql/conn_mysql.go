package xmysql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jinares/xpkg/xlog"
	"github.com/jinares/xpkg/xtools"
	"google.golang.org/grpc/codes"
	"sync"
	"time"
)

type (
	DBHandler func() (*sql.DB, error)
)

var (
	mysqlMutex sync.RWMutex
	mysqlPool  = map[string]*sql.DB{}
)

func ConnMYSQL(item MYSQLConfig) (*sql.DB, error) {
	if item.DNS == "" {
		return nil, xtools.XErr(codes.Unimplemented, "dns is empty", true)
	}
	mysqlMutex.RLock()
	mysqldb, isok := mysqlPool[item.DNS]
	mysqlMutex.RUnlock()
	if isok == true {
		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Second)
		defer cancel()
		err := mysqldb.PingContext(ctx)
		if err != nil {
			fmt.Println(fmt.Sprintf("mysql-ping-err:%s", err.Error()))
			xlog.Err(err).WithField("data", mysqldb.Stats()).Error("mysql-ping")
			return nil, xtools.MErr(err, codes.ResourceExhausted, "请稍后再试", true)
		}
		return mysqldb, nil

	}
	mysqlMutex.Lock()
	defer mysqlMutex.Unlock()
	//避免并发多次写入
	mysqldb, isok = mysqlPool[item.DNS]
	if isok {
		return mysqldb, nil
	}
	//delete(mysqlPool, item.DNS)
	mysqldb, err := sql.Open("mysql", item.DNS)
	if err != nil {
		fmt.Println(fmt.Sprintf("mysql-connection-err:%s", err.Error()))
		xlog.Err(err).Error("mysql-connection")
		return nil, xtools.XErr(codes.Internal, err.Error(), true)
	}
	mysqldb.SetMaxOpenConns(item.MaxPoolSize)
	mysqldb.SetMaxIdleConns(item.MaxIdleConns)
	if item.MaxIdleTime > 0 {
		mysqldb.SetConnMaxIdleTime(time.Duration(item.MaxIdleTime) * time.Second)
	}
	if item.ConnMaxLifetime > 0 {
		mysqldb.SetConnMaxLifetime(time.Duration(item.ConnMaxLifetime) * time.Second)
	} else if item.ConnMaxLifetime < 0 {
		mysqldb.SetConnMaxLifetime(30 * time.Second)
	} else {
		mysqldb.SetConnMaxLifetime(0)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err = mysqldb.PingContext(ctx)
	if err != nil {
		fmt.Println(fmt.Sprintf("mysql-ping2-err:%s", err.Error()))
		//mysqldb.Close()
		xlog.Err(err).Error("mysql-ping2")
		return nil, xtools.MErr(err, codes.ResourceExhausted, "请稍后再试", true)
	}

	mysqlPool[item.DNS] = mysqldb
	return mysqldb, nil
}
