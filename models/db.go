package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/IBAX-io/go-ibax/packages/consts"
	log "github.com/sirupsen/logrus"

	"github.com/IBAX-io/go-explorer/conf"

	"github.com/IBAX-io/go-ibax/packages/model"

	"gorm.io/gorm"
)

type DbTransaction struct {
	conn *gorm.DB
}
type DayBlock struct {
	Id     int64 `gorm:"not null"`
	Tx     int32 `gorm:"not null"`
	Length int64 `gorm:"not null"`
}

func isFound(db *gorm.DB) (bool, error) {
	if errors.Is(db.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return true, db.Error
}
func InitDatabase() {
	DatabaseInfo := conf.GetEnvConf().DatabaseInfo
	if err := DatabaseInfo.Initer(); err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("postgres database connect failed:")
	}
}

// GormClose is closing Gorm connection
func GormClose() error {
	if err := conf.GetEnvConf().DatabaseInfo.Close(); err != nil {
		return err
	}
	return nil
}

// StartTransaction is beginning transaction
func StartTransaction() (*DbTransaction, error) {
	conn := conf.GetDbConn().Conn().Begin()
	if conn.Error != nil {
		log.WithFields(log.Fields{"type": consts.DBError, "error": conn.Error}).Error("cannot start transaction because of connection error")
		return nil, conn.Error
	}
	return &DbTransaction{
		conn: conn,
	}, nil
}

// Rollback is transaction rollback
func (tr *DbTransaction) Rollback() {
	tr.conn.Rollback()
}

// Commit is transaction commit
func (tr *DbTransaction) Commit() error {
	return tr.conn.Commit().Error
}

// Connection returns connection of database
func (tr *DbTransaction) Connection() *gorm.DB {
	return tr.conn
}

func GetALL(tableName string, order string, v interface{}) error {
	return conf.GetDbConn().Conn().Table(tableName).Order(order).Find(v).Error
}

func GetEcosytem(id int64) (int, error) {
	var (
		keys []Key
	)
	err := conf.GetDbConn().Conn().Table("1_keys").Where("ecosystem = ?", id).Find(&keys).Error
	count := len(keys)

	return count, err
}

func GetBlockid(hash []byte) (int64, error) {
	lt := model.LogTransaction{}
	fount, err := lt.GetByHash(hash)
	if err == nil && fount {
		return lt.Block, nil
	}
	return -1, err
}

var (
	Gret []DBTransactionsInfo
)

func GetDBDealTraninfo(limit int) error {
	var (
		err error
	)
	if err = GetBlockInfoToRedis(limit); err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("GetDayBlockInfoToRedis err")
	}
	return err
}
func GetTraninfoFromRedis(limit int) (*[]ScanOutBlockTransactionRet, error) {
	var ret []ScanOutBlockTransactionRet
	var err error
	var transBlock []DayBlock

	rd := RedisParams{
		Key:   "blockChain-transaction",
		Value: "",
	}
	if err = rd.Get(); err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("GetTraninfoFromRedis getdb err")
		return nil, err
	}
	if err = json.Unmarshal([]byte(rd.Value), &transBlock); err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("GetTraninfoFromRedis json err")
		return nil, err
	}

	for i := 0; i < len(transBlock); i++ {
		var info = ScanOutBlockTransactionRet{
			BlockId:           transBlock[i].Id,
			BlockSizes:        transBlock[i].Length,
			BlockTranscations: int64(transBlock[i].Tx),
		}
		ret = append(ret, info)
	}
	return &ret, err
}
func SendTraninfoToWebsocket(dayblock []DayBlock) error {
	var ret []ScanOutBlockTransactionRet
	var err error
	for i := 0; i < len(dayblock); i++ {
		var info = ScanOutBlockTransactionRet{
			BlockId:           dayblock[i].Id,
			BlockSizes:        dayblock[i].Length,
			BlockTranscations: int64(dayblock[i].Tx),
		}
		ret = append(ret, info)
	}
	err = SendTopTransactiontps(&ret)
	if err != nil {
		return err
	}
	return nil
}

func SendTopTransactiontps(topBlockTps *[]ScanOutBlockTransactionRet) error {
	ds, err := json.Marshal(topBlockTps)
	if err != nil {
		return err
	}
	value, err := json.Marshal(trans)
	if err != nil {
		return err
	}
	rd := RedisParams{
		Key:   "blockChain-transaction",
		Value: string(value),
	}
	if err := rd.Set(); err != nil {
		return err
	}
	if err := SendTraninfoToWebsocket(trans); err != nil {
		return fmt.Errorf("SendTraninfoToWebsocket err:%s", err.Error())
	}
	return nil
}
func GetDayblockinfoFromRedis(t1, t2 int64, transBlock []Block) (int32, error) {
	var (
		dat int32
		err error
	)

	dlen := len(transBlock)
	dat = 0
	for i := 0; i < dlen; i++ {
		if transBlock[i].Time > t1 && transBlock[i].Time < t2 {
			dat += transBlock[i].Tx
		}
	}
	return dat, err
}

func GetDBDayTraninfo(day int) (*[]DBTransactionsInfo, error) {
	return &Gret, nil
}

func GetDBDayblockinfo(t1, t2 int64) (int32, error) {
	var (
		dat int32
	)
	trans := make([]model.Block, 0)
	err := conf.GetDbConn().Conn().Table("block_chain").Where(`time > ? and time < ?`, t1, t2).Find(&trans).Error
	dlen := len(trans)
	dat = 0
	for i := 0; i < dlen; i++ {
		dat += trans[i].Tx
	}
	return dat, err
}
func HasTableOrView(tr *DbTransaction, names string) bool {
	var name string
	conf.GetDbConn().Conn().Table("information_schema.tables").
		Where("table_type IN ('BASE TABLE', 'VIEW') AND table_schema NOT IN ('pg_catalog', 'information_schema') AND table_name=?", names).
		Select("table_name").Row().Scan(&name)

	return name == names
}

// IsTable returns is table exists
func IsTable(tblname string) bool {
	var name string
	conf.GetDbConn().Conn().Table("information_schema.tables").
		Where("table_type = 'BASE TABLE' AND table_schema NOT IN ('pg_catalog', 'information_schema') AND table_name=?", tblname).
		Select("table_name").Row().Scan(&name)

	return name == tblname
}
