/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import (
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/IBAX-io/go-explorer/consts"

	"github.com/IBAX-io/go-explorer/conf"

	"github.com/IBAX-io/go-ibax/packages/converter"

	"github.com/shopspring/decimal"
)

// History represent record of history table
type History struct {
	ID               int64           `gorm:"primary_key;not null"`
	Senderid         int64           `gorm:"column:sender_id;not null"`
	Recipientid      int64           `gorm:"column:recipient_id;not null"`
	SenderBalance    decimal.Decimal `gorm:"column:sender_balance;not null"`
	RecipientBalance decimal.Decimal `gorm:"column:recipient_balance;not null"`
	Amount           decimal.Decimal `gorm:"column:amount;not null"`
	Comment          string          `gorm:"column:comment;not null"`
	Blockid          int64           `gorm:"column:block_id;not null"`
	Txhash           []byte          `gorm:"column:txhash;not null"`
	Createdat        int64           `gorm:"column:created_at;not null"`
	Ecosystem        int64           `gorm:"not null"`
	Type             int64           `gorm:"not null"`
	//Createdat   time.Time       `gorm:"column:created_at;not null"`
}

type HistoryHex struct {
	ID               int64           `json:"id,omitempty"`
	Senderid         string          `json:"sender_id"`
	Recipientid      string          `json:"recipient_id"`
	SenderBalance    decimal.Decimal `json:"sender_balance"`
	RecipientBalance decimal.Decimal `json:"recipient_balance"`
	Amount           decimal.Decimal `json:"amount"`
	Comment          string          `json:"comment"`
	Blockid          int64           `json:"block_id"`
	Txhash           string          `json:"txhash"`
	Createdat        time.Time       `json:"created_at"`
	Ecosystem        int64           `json:"ecosystem"`
	Type             int64           `json:"type"`
	Ecosystemname    string          `json:"ecosystemname"`
	Token_title      string          `json:"token_title"`
	ContractName     string          `json:"contract_name"`
}

type HistoryMergeHex struct {
	Ecosystem     int64           `json:"ecosystem"`
	ID            int64           `json:"id"`
	Senderid      string          `json:"sender_id"`
	Recipientid1  string          `json:"recipientid1"`
	Recipientid2  string          `json:"recipientid2,omitempty"`
	Recipientid3  string          `json:"recipientid3,omitempty"`
	Recipientid4  string          `json:"recipientid4,omitempty"`
	Amount1       decimal.Decimal `json:"amount1,omitempty"`
	Amount2       decimal.Decimal `json:"amount2,omitempty"`
	Amount3       decimal.Decimal `json:"amount3,omitempty"`
	Amount4       decimal.Decimal `json:"amount4,omitempty"`
	Comment       string          `json:"comment"`
	Blockid       int64           `json:"blockid"`
	Txhash        string          `json:"txhash"`
	Createdat     time.Time       `json:"created_at"`
	Ecosystemname string          `json:"ecosystemname"`
	Token_title   string          `json:"token_title"`
	//Ecosystem     int64    `json:"ecosystem"`
}

type HistoryItem struct {
	//Ecosystem        int64           `json:"ecosystem"`
	//Ecosystemname    string          `json:"ecosystemname"`
	//Token_title      string          `json:"token_title"`
	Senderid    string          `json:"sender_id"`
	Recipientid string          `json:"recipient_id"`
	Amount      decimal.Decimal `json:"amount"`
}

type HistoryExplorer struct {
	Ecosystem     int64       `json:"ecosystem"`
	ID            int64       `json:"id"`
	Senderid      string      `json:"sender_id"`
	Fees          HistoryItem `json:"fees"`
	Commission    HistoryItem `json:"commission"`
	Transfer      HistoryItem `json:"transfer"`
	PackFees      HistoryItem `json:"pack_fees,omitempty"`
	Comment       string      `json:"comment"`
	Blockid       int64       `json:"blockid"`
	Txhash        string      `json:"txhash"`
	Createdat     int64       `json:"created_at"`
	CreateSetup   int64       `json:"created_setup"`
	LogoUrl       string      `json:"logo_url"`
	Ecosystemname string      `json:"ecosystemname"`
	Token_title   string      `json:"token_title"`
}

type WalletHistoryHex struct {
	Transaction int64           `json:"transaction"`
	Inamount    decimal.Decimal `json:"inamount"`
	Outamount   decimal.Decimal `json:"outamount"`
	Amount      decimal.Decimal `json:"amount,omitempty"`
}

type HistorysResult struct {
	Total int64           `json:"total"1`
	Page  int             `json:"page" example:0"`
	Limit int             `json:"limit"`
	Sum   decimal.Decimal `json:"sum,omitempty"`
	Rets  []HistoryHex    `json:"rets"`
}

type HistoryTransaction struct {
	ID            int64  `json:"id"`
	Keyid         string `json:"key_id"`
	Blockid       int64  `json:"block_id"`
	Txhash        string `json:"txhash"`
	Createdat     int64  `json:"created_at"`
	Ecosystem     int64  `json:"ecosystem"`
	Ecosystemname string `json:"ecosystemname"`
	ContractName  string `json:"contract_name"`
	//Ecosystem     int64    `json:"ecosystem"`
}

type Historys []History

// TableName returns name of table
func (th *History) TableName() string {
	return "1_history"
}

func (th *History) Get(txHash []byte) (*HistoryMergeHex, error) {
	var (
		ts  []History
		tss HistoryMergeHex
	)

	err := conf.GetDbConn().Conn().Where("txhash = ?", txHash).Order("id ASC").Find(&ts).Error
	count := len(ts)
	if err == nil && count > 0 {
		if ts[0].Blockid > 0 {
			sort.Sort(Historys(ts))

			//fmt.Println(ts)
			tss.Ecosystem = ts[0].Ecosystem
			es := Ecosystem{}
			f, err := es.Get(tss.Ecosystem)
			if f && err == nil {
				tss.Ecosystemname = es.Name
				if tss.Ecosystem == 1 {
					tss.Token_title = consts.SysEcosytemTitle
					if tss.Ecosystemname == "" {
						tss.Ecosystemname = "platform ecosystem"
					}
				} else {
					tss.Token_title = es.TokenTitle
				}
				//tss.Token_title = es.Token_title
			}

			tss.ID = ts[0].ID
			tss.Senderid = converter.AddressToString(ts[0].Senderid) //strconv.FormatInt(ts[0].Senderid, 10)
			tss.Comment = ts[0].Comment
			tss.Blockid = ts[0].Blockid
			tss.Txhash = hex.EncodeToString(ts[0].Txhash)
			fmt.Println(tss.Txhash)
			fmt.Println(string(ts[0].Txhash))
			tss.Createdat = time.Unix(ts[0].Createdat, 0)
			if count == 4 {
				tss.Recipientid1 = converter.AddressToString(ts[3].Recipientid) //strconv.FormatInt(ts[3].Recipientid, 10)
				tss.Recipientid2 = converter.AddressToString(ts[2].Recipientid) //strconv.FormatInt(ts[2].Recipientid, 10)
				tss.Recipientid3 = converter.AddressToString(ts[1].Recipientid) //strconv.FormatInt(ts[1].Recipientid, 10)
				tss.Recipientid4 = converter.AddressToString(ts[0].Recipientid) //strconv.FormatInt(ts[0].Recipientid, 10)
				tss.Amount1 = ts[3].Amount
				tss.Amount2 = ts[2].Amount
				tss.Amount3 = ts[1].Amount
				tss.Amount4 = ts[0].Amount
				//				fmt.Println(ts[2].Amount)
				//				fmt.Println(ts[1].Amount)
				//				fmt.Println(ts[0].Amount)
				//				fmt.Println(tss)
			} else if count == 3 {
				tss.Recipientid1 = converter.AddressToString(ts[2].Recipientid) //strconv.FormatInt(ts[2].Recipientid, 10)
				tss.Recipientid2 = converter.AddressToString(ts[1].Recipientid) //strconv.FormatInt(ts[1].Recipientid, 10)
				tss.Recipientid3 = converter.AddressToString(ts[0].Recipientid) //strconv.FormatInt(ts[0].Recipientid, 10)
				tss.Amount1 = ts[2].Amount
				tss.Amount2 = ts[1].Amount
				tss.Amount3 = ts[0].Amount
				//				fmt.Println(ts[2].Amount)
				//				fmt.Println(ts[1].Amount)
				//				fmt.Println(ts[0].Amount)
				//				fmt.Println(tss)
			} else if count == 2 {
				tss.Recipientid2 = converter.AddressToString(ts[1].Recipientid) //strconv.FormatInt(ts[1].Recipientid, 10)
				tss.Recipientid3 = converter.AddressToString(ts[0].Recipientid) //strconv.FormatInt(ts[0].Recipientid, 10)
				tss.Amount2 = ts[1].Amount
				tss.Amount3 = ts[0].Amount
				//tss.Amount1 =decimal.NewFromFloat(0)
			} else if count == 1 {
				tss.Recipientid1 = converter.AddressToString(ts[0].Recipientid) //strconv.FormatInt(ts[0].Recipientid, 10)
				tss.Amount1 = ts[0].Amount
			}
		}
	}
	return &tss, err
}

func (th *History) GetExplorer(logourl string, txHash []byte) (*HistoryExplorer, error) {
	var (
		ts  []History
		tss HistoryExplorer
	)
	err := conf.GetDbConn().Conn().Where("txhash = ?", txHash).Order("id ASC").Find(&ts).Error
	count := len(ts)
	if err == nil && count > 0 {
		tss.ID = ts[0].ID
		tss.Senderid = converter.AddressToString(ts[0].Senderid) //strconv.FormatInt(ts[0].Senderid, 10)
		tss.Comment = ts[0].Comment
		tss.Blockid = ts[0].Blockid
		tss.Txhash = hex.EncodeToString(ts[0].Txhash)
		tss.Createdat = ts[0].Createdat
		tss.CreateSetup = ts[0].Createdat
		tss.Ecosystem = ts[0].Ecosystem
		es := Ecosystem{}
		f, err := es.Get(tss.Ecosystem)
		if f && err == nil {
			tss.Ecosystemname = es.Name
			if tss.Ecosystem == 1 {
				tss.Token_title = consts.SysEcosytemTitle
				if tss.Ecosystemname == "" {
					tss.Ecosystemname = "platform ecosystem"
				}
			} else {
				tss.Token_title = es.TokenTitle
			}
		}
		var ts1 TransactionStatus
		found, err := ts1.DbconngetSqlite(txHash)
		if err == nil && found {
			tss.CreateSetup = ts1.Time
		}

		mb := Member{}
		fm, _ := mb.GetAccount(tss.Ecosystem, tss.Senderid)
		if fm {
			if mb.ImageID != nil {
				if *mb.ImageID != int64(0) {
					mrl, err := Loadlogo(*mb.ImageID)
					if err == nil {
						if mrl != "" {
							tss.LogoUrl = logourl + mrl
						}
					}
					if mrl != "" {
						tss.LogoUrl = logourl + mrl
					}
				}
			}
		}

		//for
		for _, ts := range ts {
			var det HistoryItem
			if ts.Type == 1 {
				//det.Ecosystem = ts.Ecosystem
				//det.Token_title = ts.t
				//det.Ecosystemname = ts.e
				det.Senderid = converter.AddressToString(ts.Senderid)
				det.Recipientid = converter.AddressToString(ts.Recipientid)
				det.Amount = ts.Amount
				tss.Fees = det
			} else if ts.Type == 2 {
				det.Senderid = converter.AddressToString(ts.Senderid)
				det.Recipientid = converter.AddressToString(ts.Recipientid)
				det.Amount = ts.Amount
				tss.Commission = det
			} else if ts.Type == 12 {
				det.Senderid = converter.AddressToString(ts.Senderid)
				det.Recipientid = converter.AddressToString(ts.Recipientid)
				det.Amount = ts.Amount
				tss.PackFees = det
			} else {
				det.Senderid = converter.AddressToString(ts.Senderid)
				det.Recipientid = converter.AddressToString(ts.Recipientid)
				det.Amount = ts.Amount
				if tss.Transfer.Senderid != "" {
					tss.Transfer.Amount.Add(det.Amount)
				} else {
					tss.Transfer = det
				}

			}
		}

	}

	return &tss, err
}

func (th *History) GetHistoryTimeList(time time.Time) (*[]History, error) {
	var (
		tss []History
	)

	err := conf.GetDbConn().Conn().Model(&History{}).Where("created_at >?", time.Unix()).Order("created_at desc").Find(&tss).Error
	return &tss, err
}

func (th *History) GetHistoryIdList(id int64) (*[]History, error) {
	var (
		tss []History
	)

	err := conf.GetDbConn().Conn().Model(&History{}).Where("id >?", id).Order("id desc").Find(&tss).Error
	return &tss, err
}

//Get is retrieving model from database
func (ts *History) GetHistoryList() (*[]History, error) {
	var (
		ret []History
	)

	err := conf.GetDbConn().Conn().Table("1_history").Order("id desc").Find(&ret).Error
	return &ret, err
}

//Get is retrieving model from database
func (th *History) GetHistorys(page int, size int, order string) (*[]HistoryHex, int64, error) {
	var (
		tss []History
		ret []HistoryHex
		num int64
	)

	err := conf.GetDbConn().Conn().Limit(size).Offset((page - 1) * size).Order(order).Find(&tss).Error
	if err != nil {
		return &ret, num, err
	}
				da.Token_title = consts.SysEcosytemTitle
				if da.Ecosystemname == "" {
					da.Ecosystemname = "platform ecosystem"
				}
			} else {
				da.Token_title = es.TokenTitle
			}
			///da.Token_title = es.Token_title
		}
		da.ID = tss[i].ID
		da.Senderid = strconv.FormatInt(tss[i].Senderid, 10)
		da.Recipientid = strconv.FormatInt(tss[i].Recipientid, 10)
		da.Amount = tss[i].Amount
		da.Comment = tss[i].Comment
		da.Blockid = tss[i].Blockid
		da.Txhash = hex.EncodeToString(tss[i].Txhash)
		da.Createdat = time.Unix(tss[i].Createdat, 0)
		ret = append(ret, da)
	}

	return &ret, num, err
}

//Get is retrieving model from database
func (th *History) GetWallets(page int, size int, wallet string, searchType string) (*[]HistoryHex, int64, decimal.Decimal, error) {
	var (
		tss []History
		ret []HistoryHex
		num int64
		//ioffet int64
		i     int64
		keyId int64
		err   error
		total decimal.Decimal
	)

	num = 0
	keyId, err = strconv.ParseInt(wallet, 10, 64)
	if err != nil {
		return &ret, num, total, err
	}
	if page < 1 || size < 1 {
		return &ret, num, total, err
	}
	if searchType == "income" {
		if err = conf.GetDbConn().Conn().Table("1_history").Where("recipient_id = ?", keyId).Count(&num).Error; err != nil {
			return &ret, num, total, err
		}
		err = conf.GetDbConn().Conn().Table("1_history").Select("sum(amount)").Where("recipient_id = ?", keyId).Row().Scan(&total)
		if err != nil {
			return &ret, num, total, err
		}
		err = conf.GetDbConn().Conn().Table("1_history").
			Where("recipient_id = ?", keyId).
			Order("id desc").Offset((page - 1) * size).Limit(size).Find(&tss).Error
	} else if searchType == "outcome" {
		if err = conf.GetDbConn().Conn().Table("1_history").Where("sender_id = ?", keyId).Count(&num).Error; err != nil {
			return &ret, num, total, err
		}
		err = conf.GetDbConn().Conn().Table("1_history").Select("sum(amount)").Where("sender_id = ?", keyId).Row().Scan(&total)
		if err != nil {
			return &ret, num, total, err
		}
		err = conf.GetDbConn().Conn().Table("1_history").
			Where("sender_id = ?", keyId).
			Order("id desc").Offset((page - 1) * size).Limit(size).Find(&tss).Error
	} else {
		if err = conf.GetDbConn().Conn().Table("1_history").Where("recipient_id = ? OR sender_id = ?", keyId, keyId).Count(&num).Error; err != nil {
			return &ret, num, total, err
		}
		err = conf.GetDbConn().Conn().Table("1_history").Select("sum(amount)").Where("recipient_id = ? OR sender_id = ?", keyId, keyId).Row().Scan(&total)
		if err != nil {
			return &ret, num, total, err
		}
		err = conf.GetDbConn().Conn().Table("1_history").
			Where("recipient_id = ? OR sender_id = ?", keyId, keyId).
			Order("id desc").Offset((page - 1) * size).Limit(size).Find(&tss).Error
	}

	//total = deal_history_total(&tss)

	count := int64(len(tss))
	//fmt.Println("tr_blocks Error: %d", num)
	//ioffet = (page - 1) * size
	//if num < page*size {
	//	size = num % size
	//}
	//if num < ioffet || num < 1 {
	//	return &ret, num, total, err
	//}
	for i = 0; i < count; i++ {
		//fmt.Println("offet Error:%d ", ioffet)
		da := HistoryHex{}
		da.Ecosystem = tss[i].Ecosystem
		es := Ecosystem{}
		f, err := es.Get(da.Ecosystem)
		if f && err == nil {
			da.Ecosystemname = es.Name
			if da.Ecosystem == 1 {
				da.Token_title = consts.SysEcosytemTitle
			} else {
				da.Token_title = es.TokenTitle
			}
			//da.Token_title = es.Token_title
		}
		da.ID = tss[i].ID
		da.Senderid = strconv.FormatInt(tss[i].Senderid, 10)
		da.Recipientid = strconv.FormatInt(tss[i].Recipientid, 10)
		da.Amount = tss[i].Amount
		da.Comment = tss[i].Comment
		da.Blockid = tss[i].Blockid
		da.Txhash = hex.EncodeToString(tss[i].Txhash)
		da.Createdat = time.Unix(tss[i].Createdat, 0)
		ret = append(ret, da)
		//ioffet++
		//fmt.Println("offet Error:%d ", ioffet)
	}

	return &ret, num, total, err
}

//Get is retrieving model from database
func (th *History) GetEcosytemWallets(id int64, page int, size int, wallet string, searchType string) (*[]HistoryHex, int64, decimal.Decimal, error) {
	var (
		tss []History
		ret []HistoryHex
		num int64
		//ioffet int64
		i     int64
		keyId int64
		err   error
		total decimal.Decimal
	)

	num = 0
	keyId, err = strconv.ParseInt(wallet, 10, 64)
	if err != nil {
		return &ret, num, total, err
	}
	if page < 1 || size < 1 {
		return &ret, num, total, err
	}
	if searchType == "income" {
		if err = conf.GetDbConn().Conn().Table("1_history").Where("recipient_id = ? and ecosystem = ?", keyId, id).Count(&num).Error; err != nil {
			return &ret, num, total, err
		}
		err = conf.GetDbConn().Conn().Table("1_history").Select("sum(amount)").Where("recipient_id = ? and ecosystem = ?", keyId, id).Row().Scan(&total)
		if err != nil {
			return &ret, num, total, err
		}
		err = conf.GetDbConn().Conn().Table("1_history").
			Where("recipient_id = ? and ecosystem = ?", keyId, id).
			Order("id desc").Offset((page - 1) * size).Limit(size).Find(&tss).Error
	} else if searchType == "outcome" {
		if err = conf.GetDbConn().Conn().Table("1_history").Where("sender_id = ? and ecosystem = ?", keyId, id).Count(&num).Error; err != nil {
			return &ret, num, total, err
		}
		err = conf.GetDbConn().Conn().Table("1_history").Select("sum(amount)").Where("sender_id = ? and ecosystem = ?", keyId, id).Row().Scan(&total)
		if err != nil {
			return &ret, num, total, err
		}
		err = conf.GetDbConn().Conn().Table("1_history").
			Where("sender_id = ? and ecosystem = ?", keyId, id).
			Order("id desc").Offset((page - 1) * size).Limit(size).Find(&tss).Error
	} else {
		if err = conf.GetDbConn().Conn().Table("1_history").Where("(recipient_id = ? OR sender_id = ?) and ecosystem = ?", keyId, keyId, id).Count(&num).Error; err != nil {
			return &ret, num, total, err
		}
		err = conf.GetDbConn().Conn().Table("1_history").Select("sum(amount)").Where("(recipient_id = ? OR sender_id = ?) and ecosystem = ?", keyId, keyId, id).Row().Scan(&total)
		if err != nil {
			return &ret, num, total, err
		}
		err = conf.GetDbConn().Conn().Table("1_history").
			Where("(recipient_id = ? OR sender_id = ?) and ecosystem = ?", keyId, keyId, id).
			Order("id desc").Offset((page - 1) * size).Limit(size).Find(&tss).Error
	}

	//total = deal_history_total(&tss)

	count := int64(len(tss))
	//fmt.Println("tr_blocks Error: %d", num)
	//ioffet = (page - 1) * size
	//if num < page*size {
	//	size = num % size
	//}
	//if num < ioffet || num < 1 {
	//	return &ret, num, total, err
	//}
	for i = 0; i < count; i++ {
		//fmt.Println("offet Error:%d ", ioffet)
		da := HistoryHex{}
		da.Ecosystem = tss[i].Ecosystem
		es := Ecosystem{}
		f, err := es.Get(da.Ecosystem)
		if f && err == nil {
			da.Ecosystemname = es.Name
			if da.Ecosystem == 1 {
				da.Token_title = consts.SysEcosytemTitle
			} else {
				da.Token_title = es.TokenTitle
			}
			//da.Token_title = es.Token_title
		}
		da.ID = tss[i].ID
		da.Senderid = strconv.FormatInt(tss[i].Senderid, 10)
		da.Recipientid = strconv.FormatInt(tss[i].Recipientid, 10)
		da.Amount = tss[i].Amount
		da.Comment = tss[i].Comment
		da.Blockid = tss[i].Blockid
		da.Txhash = hex.EncodeToString(tss[i].Txhash)
		da.Createdat = time.Unix(tss[i].Createdat, 0)
		ret = append(ret, da)
		//i++
		//fmt.Println("offet Error:%d ", ioffet)
	}

	return &ret, num, total, err
}

func (th *History) GetEcosytemTransactionWallets(id int64, page int, size int, wallet string, searchType string) (*HistorysResult, error) {
	var (
		tss []History
		ret []HistoryHex
		num int64
		//ioffet int64
		i     int64
		keyId int64
		err   error
		total decimal.Decimal
		rets  HistorysResult
	)
	rets.Limit = size
	rets.Page = page

	num = 0

	var so ScanOut
	lb, err := so.GetRedisdashboard()
	if err != nil {
		return &rets, err
	}
	bid := lb.Blockid

	keyId = converter.StringToAddress(wallet)
	if wallet == "0000-0000-0000-0000-0000" {
	} else if keyId == 0 {
		return &rets, errors.New("wallet does not meet specifications")
	}
	if page < 1 || size < 1 {
		return &rets, err
	}
	if searchType == "income" {
		if err = conf.GetDbConn().Conn().Table("1_history").Where("recipient_id = ? and ecosystem = ? and block_id <= ?", keyId, id, bid).Count(&num).Error; err != nil {
			return &rets, err
		}
		err = conf.GetDbConn().Conn().Table("1_history").Select("sum(amount)").Where("recipient_id = ? and ecosystem = ? and block_id <= ? ", keyId, id, bid).Row().Scan(&total)
		if err != nil {
			return &rets, err
		}
		err = conf.GetDbConn().Conn().Table("1_history").
			Where("recipient_id = ? and ecosystem = ? and block_id <= ?", keyId, id, bid).
			Order("id desc").Offset((page - 1) * size).Limit(size).Find(&tss).Error
	} else if searchType == "outcome" {
		if err = conf.GetDbConn().Conn().Table("1_history").Where("sender_id = ? and ecosystem = ? and block_id <= ? ", keyId, id, bid).Count(&num).Error; err != nil {
			return &rets, err
		}
		err = conf.GetDbConn().Conn().Table("1_history").Select("sum(amount)").Where("sender_id = ? and ecosystem = ? and block_id <= ? ", keyId, id, bid).Row().Scan(&total)
		if err != nil {
			return &rets, err
		}
		err = conf.GetDbConn().Conn().Table("1_history").
			Where("sender_id = ? and ecosystem = ? and block_id <= ?", keyId, id, bid).
			Order("id desc").Offset((page - 1) * size).Limit(size).Find(&tss).Error
	} else {
		if err = conf.GetDbConn().Conn().Table("1_history").Where("(recipient_id = ? OR sender_id = ?) and ecosystem = ? and block_id <= ?", keyId, keyId, id, bid).Count(&num).Error; err != nil {
			return &rets, err
		}
		err = conf.GetDbConn().Conn().Table("1_history").Select("sum(amount)").Where("(recipient_id = ? OR sender_id = ?) and ecosystem = ? and block_id <= ?", keyId, keyId, id, bid).Row().Scan(&total)
		if err != nil {
			return &rets, err
		}
		err = conf.GetDbConn().Conn().Table("1_history").
			Where("(recipient_id = ? OR sender_id = ?) and ecosystem = ? and block_id <= ?", keyId, keyId, id, bid).
			Order("id desc").Offset((page - 1) * size).Limit(size).Find(&tss).Error
	}

	count := int64(len(tss))
	for i = 0; i < count; i++ {
		//fmt.Println("offet Error:%d ", ioffet)
		da := HistoryHex{}
		da.Ecosystem = tss[i].Ecosystem
		es := Ecosystem{}
		f, err := es.Get(da.Ecosystem)
		if f && err == nil {
			da.Ecosystemname = es.Name
			if da.Ecosystem == 1 {
				da.Token_title = consts.SysEcosytemTitle
				if da.Ecosystemname == "" {
					da.Ecosystemname = "platform ecosystem"
				}
			} else {
				da.Token_title = es.TokenTitle
			}
			//da.Token_title = es.Token_title
		}

		da.ID = tss[i].ID
		da.Senderid = converter.AddressToString(tss[i].Senderid)       //strconv.FormatInt(tss[i].Senderid, 10)
		da.Recipientid = converter.AddressToString(tss[i].Recipientid) //strconv.FormatInt(tss[i].Recipientid, 10)
		da.Type = tss[i].Type
		da.Amount = tss[i].Amount
		da.Comment = tss[i].Comment
		da.Blockid = tss[i].Blockid
		da.Txhash = hex.EncodeToString(tss[i].Txhash)
		da.Createdat = time.Unix(tss[i].Createdat, 0)
		var bt BlockTxDetailedInfoHex
		//ft, errt := bt.GetByHash_Sqlite(da.Txhash)
		ft, errt := bt.GetDb_txdetailedHash(da.Txhash)
		if errt == nil && ft {
			da.ContractName = bt.ContractName
		} else {
			//fmt.Println(errt)
		}

		ret = append(ret, da)
		//i++
		//fmt.Println("offet Error:%d ", ioffet)
	}
	rets.Total = num
	rets.Sum = total
	rets.Rets = ret
	return &rets, err
}

//Get is retrieving model from database
func (th *History) GetWalletTotals(wallet string) (*WalletHistoryHex, error) {
	var (
		tss1  []History
		tss2  []History
		ret   WalletHistoryHex
		keyId int64
		err   error
	)

	keyId, err = strconv.ParseInt(wallet, 10, 64)
	if err != nil {
		return &ret, err
	}

	err = conf.GetDbConn().Conn().Table("1_history").
		Where("recipient_id = ?", keyId).
		Order("created_at desc").Find(&tss1).Error
	if err != nil {
		return &ret, err
	}
	err = conf.GetDbConn().Conn().Table("1_history").
		Where("sender_id = ?", keyId).
		Order("created_at desc").Find(&tss2).Error
	if err != nil {
		return &ret, err
	}
	ret.Transaction = int64(len(tss1)) + int64(len(tss2))
	ret.Inamount = deal_history_total(&tss1)
	ret.Outamount = deal_history_total(&tss2)

	return &ret, err
}

//Get is retrieving model from database
func (th *History) GetWalletHistoryTotals(id int64, wallet string) (*WalletHistoryHex, error) {
	var (
		ret    WalletHistoryHex
		keyId  int64
		scount int64
		rcount int64
		in     string
		out    string
		err    error
	)

	keyId, err = strconv.ParseInt(wallet, 10, 64)
	if err != nil {
		return &ret, err
	}

	err = conf.GetDbConn().Conn().Table("1_history").
		Where("recipient_id = ? and ecosystem = ?", keyId, id).
		Count(&rcount).Error
	if err != nil {
		return &ret, err
	}
	if rcount > 0 {
		err = conf.GetDbConn().Conn().Table("1_history").Select("sum(amount)").Where("recipient_id = ? and ecosystem = ?", keyId, id).Row().Scan(&in)
		if err != nil {
			return &ret, err
		}
	} else {
		in = "0"
	}

	err = conf.GetDbConn().Conn().Table("1_history").
		Where("sender_id = ? and ecosystem = ?", keyId, id).
		Count(&scount).Error
	if err != nil {
		return &ret, err
	}
	if scount > 0 {
		err = conf.GetDbConn().Conn().Table("1_history").Select("sum(amount)").Where("sender_id = ? and ecosystem = ?", keyId, id).Row().Scan(&out)
		if err != nil {
			return &ret, err
		}
	} else {
		out = "0"
	}

	din, err := decimal.NewFromString(in)
	if err != nil {
		return &ret, err
	}
	dout, err := decimal.NewFromString(out)
	if err != nil {
		return &ret, err
	}
	ret.Transaction = int64(scount + rcount)
	ret.Inamount = din
	ret.Outamount = dout

	return &ret, err
}

//Get is retrieving model from database
func (th *History) GetWalletHistoryTotalbykeyid(id, keyId int64) (*WalletHistoryHex, error) {
	var (
		ret WalletHistoryHex
		//keyId  int64
		scount int64
		rcount int64
		in     string
		out    string
		err    error
	)

	//keyId, err = strconv.ParseInt(wallet, 10, 64)
	//if err != nil {
	//	return &ret, err
	//}

	err = conf.GetDbConn().Conn().Table("1_history").
		Where("recipient_id = ? and ecosystem = ?", keyId, id).
		Count(&rcount).Error
	if err != nil {
		return &ret, err
	}
	if rcount > 0 {
		err = conf.GetDbConn().Conn().Table("1_history").Select("sum(amount)").Where("recipient_id = ? and ecosystem = ?", keyId, id).Row().Scan(&in)
		if err != nil {
			return &ret, err
		}
	} else {
		in = "0"
	}

	err = conf.GetDbConn().Conn().Table("1_history").
		Where("sender_id = ? and ecosystem = ?", keyId, id).
		Count(&scount).Error
	if err != nil {
		return &ret, err
	}
	if scount > 0 {
		err = conf.GetDbConn().Conn().Table("1_history").Select("sum(amount)").Where("sender_id = ? and ecosystem = ?", keyId, id).Row().Scan(&out)
		if err != nil {
			return &ret, err
		}
	} else {
		out = "0"
	}

	din, err := decimal.NewFromString(in)
	if err != nil {
		return &ret, err
	}
	dout, err := decimal.NewFromString(out)
	if err != nil {
		return &ret, err
	}
	ret.Transaction = int64(scount + rcount)
	ret.Inamount = din
	ret.Outamount = dout

	return &ret, err
}

func (u Historys) Len() int {
	return len(u)
}

func (u Historys) Less(i, j int) bool {
	dat := u[i].Amount.Cmp(u[j].Amount)
	return dat < 0 //sort by id if id is the same sort by name...
}

func (u Historys) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}

func (th *History) Get_Sqlite(txHash []byte) (*HistoryMergeHex, error) {
	var (
		ts  []History
		tss HistoryMergeHex
		//i   int
	)
	err := conf.GetDbConn().Conn().Where("txhash = ?", txHash).Find(&ts).Error
	count := len(ts)
	if err == nil && count > 0 {
		if ts[0].Blockid > 0 {
			//fmt.Println(ts)
			sort.Sort(Historys(ts))
			//fmt.Println(ts)
			tss.Ecosystem = ts[0].Ecosystem
			es := Ecosystem{}
			f, err := es.Get(tss.Ecosystem)
			if f && err == nil {
				tss.Ecosystemname = es.Name
				if tss.Ecosystem == 1 {
					tss.Token_title = consts.SysEcosytemTitle
				} else {
					tss.Token_title = es.TokenTitle
				}
				//tss.Token_title = es.Token_title
			}
			tss.ID = ts[0].ID
			tss.Senderid = strconv.FormatInt(ts[0].Senderid, 10)
			tss.Comment = ts[0].Comment
			tss.Blockid = ts[0].Blockid
			tss.Txhash = hex.EncodeToString(ts[0].Txhash)
			tss.Createdat = time.Unix(ts[0].Createdat, 0)
			if count == 3 {
				tss.Recipientid1 = strconv.FormatInt(ts[2].Recipientid, 10)
				tss.Recipientid2 = strconv.FormatInt(ts[1].Recipientid, 10)
				tss.Recipientid3 = strconv.FormatInt(ts[0].Recipientid, 10)
				tss.Amount1 = ts[2].Amount
				tss.Amount2 = ts[1].Amount
				tss.Amount3 = ts[0].Amount
				//				fmt.Println(ts[2].Amount)
				//				fmt.Println(ts[1].Amount)
				//				fmt.Println(ts[0].Amount)
				//				fmt.Println(tss)
			} else if count == 2 {
				tss.Recipientid1 = strconv.FormatInt(ts[1].Recipientid, 10)
				tss.Recipientid2 = strconv.FormatInt(ts[0].Recipientid, 10)
				tss.Amount1 = ts[1].Amount
				tss.Amount2 = ts[0].Amount
			} else if count == 1 {
				tss.Recipientid1 = strconv.FormatInt(ts[0].Recipientid, 10)
				tss.Amount1 = ts[0].Amount
			}
		}
	}
	return &tss, err
}

//Get is retrieving model from database
func (th *History) GetHistorys_Sqlite(page int, size int, order string) (*[]HistoryHex, int64, error) {
	var (
		tss []History
		ret []HistoryHex
		num int64
	)

	err := conf.GetDbConn().Conn().Limit(size).Offset((page - 1) * size).Order(order).Find(&tss).Error
	if err != nil {
		return &ret, num, err
	}

	err = conf.GetDbConn().Conn().Table("1_history").Count(&num).Error
	if err != nil {
		return &ret, num, err
	}
	for i := 0; i < len(tss); i++ {
		//fmt.Println("offet Error:%d ", ioffet)
		da := HistoryHex{}
		da.Ecosystem = tss[i].Ecosystem

		es := Ecosystem{}
		f, err := es.Get(da.Ecosystem)
		if f && err == nil {
			da.Ecosystemname = es.Name
			if da.Ecosystem == 1 {
				da.Token_title = consts.SysEcosytemTitle
			} else {
				da.Token_title = es.TokenTitle
			}
			//da.Token_title = es.Token_title
		}
		da.ID = tss[i].ID
		da.Senderid = strconv.FormatInt(tss[i].Senderid, 10)
		da.Recipientid = strconv.FormatInt(tss[i].Recipientid, 10)
		da.Amount = tss[i].Amount
		da.Comment = tss[i].Comment
		da.Blockid = tss[i].Blockid
		da.Txhash = hex.EncodeToString(tss[i].Txhash)
		da.Createdat = time.Unix(tss[i].Createdat, 0)
		ret = append(ret, da)
	}

	return &ret, num, err
}

//Get is retrieving model from database
func (th *History) GetWallets_Sqlite(page int, size int, wallet string, searchType string) (*[]HistoryHex, int64, decimal.Decimal, error) {
	var (
		tss []History
		ret []HistoryHex
		num int64
		//ioffet int64
		i     int64
		keyId int64
		err   error
		total decimal.Decimal
	)
	num = 0
	keyId, err = strconv.ParseInt(wallet, 10, 64)
	if err != nil {
		return &ret, num, total, err
	}
	if page < 1 || size < 1 {
		return &ret, num, total, err
	}

	if searchType == "income" {
		if err = conf.GetDbConn().Conn().Table("1_history").Where("recipient_id = ?", keyId).Count(&num).Error; err != nil {
			return &ret, num, total, err
		}
		err = conf.GetDbConn().Conn().Table("1_history").Select("sum(amount)").Where("recipient_id = ?", keyId).Row().Scan(&total)
		if err != nil {
			return &ret, num, total, err
		}
		err = conf.GetDbConn().Conn().Table("1_history").
			Where("recipient_id = ?", keyId).
			Order("id desc").Offset((page - 1) * size).Limit(size).Find(&tss).Error
	} else if searchType == "outcome" {
		if err = conf.GetDbConn().Conn().Table("1_history").Where("sender_id = ?", keyId).Count(&num).Error; err != nil {
			return &ret, num, total, err
		}
		err = conf.GetDbConn().Conn().Table("1_history").Select("sum(amount)").Where("sender_id = ?", keyId).Row().Scan(&total)
		if err != nil {
			return &ret, num, total, err
		}
		err = conf.GetDbConn().Conn().Table("1_history").
			Where("sender_id = ?", keyId).
			Order("id desc").Offset((page - 1) * size).Limit(size).Find(&tss).Error
	} else {
		if err = conf.GetDbConn().Conn().Table("1_history").Where("recipient_id = ? OR sender_id = ?", keyId, keyId).Count(&num).Error; err != nil {
			return &ret, num, total, err
		}
		err = conf.GetDbConn().Conn().Table("1_history").Select("sum(amount)").Where("recipient_id = ? OR sender_id = ?", keyId, keyId).Row().Scan(&total)
		if err != nil {
			return &ret, num, total, err
		}
		err = conf.GetDbConn().Conn().Table("1_history").
			Where("recipient_id = ? OR sender_id = ?", keyId, keyId).
			Order("id desc").Offset((page - 1) * size).Limit(size).Find(&tss).Error
	}
	//total = deal_history_total(&tss)

	count := int64(len(tss))
	//num = int64(len(tss))
	//ioffet = (page - 1) * size
	//if num < page*size {
	//	size = num % size
	//}
	//if num < ioffet || num < 1 {
	//	return &ret, num, total, err
	//}

	for i = 0; i < count; i++ {
		//fmt.Println("offet Error:%d ", ioffet)
		da := HistoryHex{}
		da.Ecosystem = tss[i].Ecosystem
		es := Ecosystem{}
		f, err := es.Get(da.Ecosystem)
		if f && err == nil {
			da.Ecosystemname = es.Name
			if da.Ecosystem == 1 {
				da.Token_title = consts.SysEcosytemTitle
			} else {
				da.Token_title = es.TokenTitle
			}
			//da.Token_title = es.Token_title
		}
		da.ID = tss[i].ID
		da.Senderid = strconv.FormatInt(tss[i].Senderid, 10)
		da.Recipientid = strconv.FormatInt(tss[i].Recipientid, 10)
		da.Amount = tss[i].Amount
		da.Comment = tss[i].Comment
		da.Blockid = tss[i].Blockid
		da.Txhash = hex.EncodeToString(tss[i].Txhash)
		da.Createdat = time.Unix(tss[i].Createdat, 0)
		ret = append(ret, da)
		//ioffet++
		//fmt.Println("offet Error:%d ", ioffet)
	}

	return &ret, num, total, err
}

//Get is retrieving model from database
func (th *History) GetWallets_EcosytemSqlite(id int64, page int, size int, wallet string, searchType string) (*[]HistoryHex, int64, decimal.Decimal, error) {
	var (
		tss []History
		ret []HistoryHex
		num int64
		//ioffet int64
		i     int64
		keyId int64
		err   error
		total decimal.Decimal
	)
	num = 0
	keyId, err = strconv.ParseInt(wallet, 10, 64)
	if err != nil {
		return &ret, num, total, err
	}
	if page < 1 || size < 1 {
		return &ret, num, total, err
	}

	if searchType == "income" {
		if err = conf.GetDbConn().Conn().Table("1_history").Where("recipient_id = ? and ecosystem = ?", keyId, id).Count(&num).Error; err != nil {
			return &ret, num, total, err
		}
		err = conf.GetDbConn().Conn().Table("1_history").Select("sum(amount)").Where("recipient_id = ? and ecosystem = ?", keyId, id).Row().Scan(&total)
		if err != nil {
			return &ret, num, total, err
		}
		err = conf.GetDbConn().Conn().Table("1_history").
			Where("recipient_id = ? and ecosystem = ?", keyId, id).
			Order("id desc").Offset((page - 1) * size).Limit(size).Find(&tss).Error
	} else if searchType == "outcome" {
		if err = conf.GetDbConn().Conn().Table("1_history").Where("sender_id = ? and ecosystem = ?", keyId, id).Count(&num).Error; err != nil {
			return &ret, num, total, err
		}
		err = conf.GetDbConn().Conn().Table("1_history").Select("sum(amount)").Where("sender_id = ? and ecosystem = ?", keyId, id).Row().Scan(&total)
		if err != nil {
			return &ret, num, total, err
		}
		err = conf.GetDbConn().Conn().Table("1_history").
			Where("sender_id = ? and ecosystem = ?", keyId, id).
			Order("id desc").Offset((page - 1) * size).Limit(size).Find(&tss).Error
	} else {
		if err = conf.GetDbConn().Conn().Table("1_history").Where("(recipient_id = ? OR sender_id = ?) and ecosystem = ? ", keyId, keyId, id).Count(&num).Error; err != nil {
			return &ret, num, total, err
		}
		err = conf.GetDbConn().Conn().Table("1_history").Select("sum(amount)").Where("(recipient_id = ? OR sender_id = ?) and ecosystem = ? ", keyId, keyId, id).Row().Scan(&total)
		if err != nil {
			return &ret, num, total, err
		}
		err = conf.GetDbConn().Conn().Table("1_history").
			Where("(recipient_id = ? OR sender_id = ?) and ecosystem = ? ", keyId, keyId, id).
			Order("id desc").Offset((page - 1) * size).Limit(size).Find(&tss).Error
	}
	//total = deal_history_total(&tss)

	count := int64(len(tss))

	for i = 0; i < count; i++ {
		da := HistoryHex{}
		da.Ecosystem = tss[i].Ecosystem
		es := Ecosystem{}
		f, err := es.Get(da.Ecosystem)
		if f && err == nil {
			da.Ecosystemname = es.Name
			if da.Ecosystem == 1 {
				da.Token_title = consts.SysEcosytemTitle
			} else {
				da.Token_title = es.TokenTitle
			}
		}
		da.ID = tss[i].ID
		da.Senderid = strconv.FormatInt(tss[i].Senderid, 10)
		da.Recipientid = strconv.FormatInt(tss[i].Recipientid, 10)
		da.Amount = tss[i].Amount
		da.Comment = tss[i].Comment
		da.Blockid = tss[i].Blockid
		da.Txhash = hex.EncodeToString(tss[i].Txhash)
		da.Createdat = time.Unix(tss[i].Createdat, 0)
		ret = append(ret, da)
	}

	return &ret, num, total, err
}

func (th *History) GetWalletTotals_Sqlites(wallet string) (*WalletHistoryHex, error) {
	var (
		tss1  []History
		tss2  []History
		ret   WalletHistoryHex
		keyId int64
		err   error
	)

	keyId, err = strconv.ParseInt(wallet, 10, 64)
	if err != nil {
		return &ret, err
	}

	err = conf.GetDbConn().Conn().Table("1_history").
		Where("recipient_id = ?", keyId).
		Order("created_at desc").Find(&tss1).Error
	if err != nil {
		return &ret, err
	}
	err = conf.GetDbConn().Conn().Table("1_history").
		Where("sender_id = ?", keyId).
		Order("created_at desc").Find(&tss2).Error
	if err != nil {
		return &ret, err
	}
	ret.Transaction = int64(len(tss1)) + int64(len(tss2))
	ret.Inamount = deal_history_total(&tss1)
	ret.Outamount = deal_history_total(&tss2)

	return &ret, err
}
func deal_history_total(objArr *[]History) decimal.Decimal {
	var (
		total decimal.Decimal
	)
	for _, val := range *objArr {
		total = total.Add(val.Amount)
	}
	return total
}
