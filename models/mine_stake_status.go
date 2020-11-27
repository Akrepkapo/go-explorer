package models

import (
	"github.com/IBAX-io/go-explorer/conf"
	"github.com/shopspring/decimal"
	//"lib.venas.io/IBAX/go-ibax/packages/converter"
)

//MinePledgeStatus example
type MinePledgeStatus struct {
	Id           int64
	Number       int64           `gorm:"null" example:"1"`                                        //
	Devid        int64           `gorm:"primary_key;not null" example:"1823-6253-5248-2211-6348"` //ID
	Keyid        int64           `gorm:"not null" example:"7994306939897545753"`                  //ID
	Poolid       int64           `gorm:"not null" example:"7994306939897545753"`
	MineType     int64           `gorm:"not null" example:"1"`          //
	MineNumber   string          `gorm:"not null" example:"P9Mv0FeQ73"` //
	MineCapacity int64           `gorm:"not null" example:"1"`
	Cycle        int64           `gorm:"not null" example:"30"`            //
	Amount       decimal.Decimal `gorm:"not null default 0" example:"100"` //
	Expired      int64           `gorm:"null" `
	Status       int64           `gorm:"null" example:"1"`                       //
	Online       int64           `gorm:"null" example:"1"`                       //
	Review       int64           `gorm:"null default 0" example:"1"`             //
	Count        int64           `gorm:"null default 0" example:"1"`             //
	Pledges      int64           `gorm:"null"  example:"0"`                      //
	Transfers    int64           `gorm:"null"  example:"0"`                      //
	Stime        int64           `gorm:"not null" example:"2019-07-19 17:45:31"` //
	Etime        int64           `gorm:"not null" example:"2019-07-19 17:45:31"` //
	DateUpdated  int64           `gorm:"not null" example:"2019-07-19 17:45:31"` //
	DateCreated  int64           `gorm:"not null default 0"`                     //
}

// TableName returns name of table
func (MinePledgeStatus) TableName() string {
	return `1_v_mine_pledge_status_info`
		err := conf.GetDbConn().Conn().Table("1_v_mine_stake_status_info").Select("count(*)").Where("(mine_type = ? or mine_type = ?) and online = ? ", 2, 1, 1).Row().Scan(&in)
		if err != nil {
			return gcount, in, err
		}

		err = conf.GetDbConn().Conn().Table("1_v_mine_stake_status_info").Select("count(*)").Where("(mine_type = ?) and online = ? ", 2, 1).Row().Scan(&gcount)
		if err != nil {
			return gcount, in, err
		}
	}
	return gcount, in, nil

}
