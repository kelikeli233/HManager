package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"time"
)

type MySQLDB struct {
	db *sql.DB
}

func connectMySQL(args ConnectionArgs) (*MySQLDB, error) {
	if args.MasterDSN == "" {
		return nil, fmt.Errorf("数据库链接为空")
	}
	logrus.Debugf("连接MYSQL数据库信息：%s", args.MasterDSN)

	connection, err := sql.Open("mysql", args.MasterDSN)
	if err != nil {
		//logrus.Debugf("",err)
		return nil, fmt.Errorf("数据库连接失败：%s", err)
	}

	connection.SetMaxIdleConns(args.MaxIdleConns)
	connection.SetMaxOpenConns(args.MaxOpenConns)
	connection.SetConnMaxLifetime(time.Duration(args.MaxConnLifetimeSeconds) * time.Second)

	database := &MySQLDB{
		db: connection,
	}
	err = database.Ping()
	if err != nil {
		return nil, err
	}
	log.Println("数据库连接成功")
	logrus.Debugln(args.MasterDSN)
	return database, nil
}

func (d MySQLDB) GenericQuery(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := d.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil

}

func (d MySQLDB) GenericBegin(begin func(tx *sql.Tx) error) error {
	tx, err := d.db.Begin()
	if err != nil {
		logrus.Errorf(":%v", err)
		return err
	}

	if err := begin(tx); err != nil {
		logrus.Errorf(":%v\n", err)
		tx.Rollback()
		//发起ping 不通的话 尝试重连
		return err
	}

	if err := tx.Commit(); err != nil {
		logrus.Errorln(err)
		return err
	}

	return nil
}

func (d MySQLDB) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	fmt.Println(d.db.PingContext(ctx))
	if err := d.db.PingContext(ctx); err != nil {

		logrus.Debugf("Failed to ping database: %s", err)
		return err
	}
	logrus.Debugf("Pinged database...")
	return nil
}
