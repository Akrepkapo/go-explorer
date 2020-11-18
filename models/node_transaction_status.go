/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import (
	"encoding/hex"
	"time"

	"github.com/IBAX-io/go-explorer/conf"
	"gorm.io/gorm"

	log "github.com/sirupsen/logrus"
)

var (
	GNodeStatusTranHash map[string]TransactionStatus
)

func (ts *TransactionStatus) GetNodecount(db *gorm.DB) (int64, error) {
	var (
		count int64
	)
	err := db.Table("transactions_status").Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, err
}

func (ts *TransactionStatus) DBconnGetTransactionlist(db *gorm.DB) (*[]TransactionStatus, error) {
	var (
		tss []TransactionStatus
	)

	err := db.Order("time desc").Find(&tss).Error
	//num :=int64(len(tss))
	return &tss, err
}

func (ts *TransactionStatus) DBconnGetTimelimit(db *gorm.DB, time time.Time) (*[]TransactionStatus, error) {
	var (
		tss []TransactionStatus
	)

	err := db.Where("time >= ?", time.Unix()).Order("time desc").Find(&tss).Error
	if err != nil {
		return nil, err
	}

	return &tss, err
}

func (ts *TransactionStatus) DbconngetSqlite(transactionHash []byte) (bool, error) {
	return isFound(conf.GetDbConn().Conn().Where("hash = ?", transactionHash).First(ts))
}

func (ts *TransactionStatus) DBconnGetcount_Sqlite() (int64, error) {
	var (
		count int64
	)
	err := conf.GetDbConn().Conn().Table("transactions_status").Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, err
}

func (ts *TransactionStatus) DBconnGetTimelimit_Sqlite(time time.Time) (*[]TransactionStatus, error) {
	var (
		tss []TransactionStatus
	)
	err := conf.GetDbConn().Conn().Where("time >= ?", time.Unix()).Find(&tss).Order("time").Error
	if err != nil {
		return nil, err
	}

	return &tss, err
}

func DbconnbatchupdateSqlite(objarr *[]TransactionStatus) error {
		for i := 0; i < count; {
			if i+100 < count {
				s := dat[i : i+100]
				err := DbconnbatchupdateSqlite(&s)
				if err != nil {
					log.Info("node TransactionStatus update count err: " + err.Error())
				}
				i += 100
			} else {
				s := dat[i:]
				err := DbconnbatchupdateSqlite(&s)
				if err != nil {
					log.Info("node TransactionStatus update count err: " + err.Error())
				}
				i = count
			}
		}
	}

	if len(*ret1) != 0 {
		DbconnbatchupdateSqlite(ret1)
	}
	return nil
}

func DbconndealReduplictionTransactionstatus(objArr *[]TransactionStatus) (*[]TransactionStatus, *[]TransactionStatus) {
	var (
		ret  []TransactionStatus
		ret1 []TransactionStatus
	)
	if GNodeStatusTranHash == nil {
		GNodeStatusTranHash = make(map[string]TransactionStatus)
	}
	for _, val := range *objArr {
		key := hex.EncodeToString(val.Hash)
		dat, ok := GNodeStatusTranHash[key]
		if ok {
			if val.Error != dat.Error || val.BlockID != dat.BlockID {
				ret1 = append(ret1, val)
				GNodeStatusTranHash[key] = val
			}
		} else {
			GNodeStatusTranHash[key] = val
			ret = append(ret, val)
		}
	}
	return &ret, &ret1
}
