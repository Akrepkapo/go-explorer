package storage

import (
	"fmt"
	"time"

	"github.com/IBAX-io/go-ibax/packages/smart"

	MaxIdle int    `yaml:"max_idle"`
	MaxOpen int    `yaml:"max_open"`
}

func (d *DatabaseModel) Initer() (err error) {
	dsn := fmt.Sprintf("%s TimeZone=UTC", d.Connect)
	pgdb, err = gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
	}), &gorm.Config{
		//AllowGlobalUpdate: true,                                  //allow global update
		Logger: logger.Default.LogMode(logger.Silent), // start Logger，show detail log
	})
	if err != nil {
		return err
	}
	sqlDB, err := pgdb.DB()
	if err != nil {
		return err
	}
	sqlDB.SetConnMaxLifetime(time.Minute * 10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(20)
	model.DBConn = pgdb
	if err = syspar.SysUpdate(nil); err != nil {
		return err
	}
	smart.InitVM()
	if err := smart.LoadContracts(); err != nil {
		return err
	}
	// Stats returns database statistics.
	//sqlDB.Stats()
	return nil

}

func (d *DatabaseModel) Conn() *gorm.DB {
	return pgdb
}

func (d *DatabaseModel) Close() error {
	if pgdb != nil {
		sqlDB, err := pgdb.DB()
		if err != nil {
			return err
		}
		if err = sqlDB.Close(); err != nil {
			return err
		}
		pgdb = nil
	}
	return nil
}

func GormDBInit(engine, connect string) (*gorm.DB, error) {
	conn, err := gorm.Open(postgres.New(postgres.Config{
		DSN: connect,
	}), &gorm.Config{
		AllowGlobalUpdate: true,                                  //allow global update
		Logger:            logger.Default.LogMode(logger.Silent), // start Logger，show detail log
	})
	if err != nil {
		return nil, err
	}
	db, err := conn.DB()
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(time.Minute * 5)
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	return conn, nil
}
