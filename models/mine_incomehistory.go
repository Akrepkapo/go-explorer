/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/
package models

import (
	"github.com/IBAX-io/go-explorer/conf"
	"github.com/shopspring/decimal"
)

type MineIncomehistory struct {
	ID                      int64           `gorm:"primary_key;not null"`
	Devid                   int64           `gorm:"not null"`
	Keyid                   int64           `gorm:"not null"`
	Mineid                  int64           `gorm:"not null"`

func (ts *MineIncomehistory) Get(hash []byte) (bool, error) {
	f, err := ts.GetRedisByhash(hash)
	if f && err == nil {
		return f, err
	}

	f, err = isFound(conf.GetDbConn().Conn().Where("mine_incomehistory_hash = ?", hash).First(ts))
	if f && err == nil {
		ts.Insert_redis()
		return f, err
	}
	return f, err
}

func (m *MineIncomehistory) GetID(id int64) (bool, error) {
	if HasTableOrView(nil, "1_mine_incomehistory") {
		return isFound(conf.GetDbConn().Conn().Where("block_id = ?", id).First(m))
	} else {
		return false, nil
	}
}
