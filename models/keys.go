/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/IBAX-io/go-explorer/consts"

	"github.com/IBAX-io/go-explorer/conf"

	"github.com/shopspring/decimal"

	//"github.com/IBAX-io/go-explorer/conf"
	"strconv"
	"strings"

	"github.com/IBAX-io/go-ibax/packages/converter"
)

// Key is model
type Key struct {
	Ecosystem   int64
	ID          int64           `gorm:"primary_key;not null"`
	PublicKey   []byte          `gorm:"column:pub;not null"`
	Amount      decimal.Decimal `gorm:"not null"`
	Maxpay      decimal.Decimal `gorm:"not null"`
	MineLock    decimal.Decimal `gorm: "column:mine_lock;`
	PoolLock    decimal.Decimal `gorm: "column:pool_lock;`
	MintSurplus decimal.Decimal `gorm: "column:mintsurplus;`
	Multi       int64           `gorm:"not null"`
	Deleted     int64           `gorm:"not null"`
	Blocked     int64           `gorm:"not null"`
}

type KeyHex struct {
	Ecosystem     int64           `json:"ecosystem"`
	ID            string          `json:"id"`
	PublicKey     string          `json:"publickey"`
	Amount        decimal.Decimal `json:"amount"`
	Maxpay        decimal.Decimal `json:"maxpay"`
	Multi         int64           `json:"multi"`
	Deleted       int64           `json:"deleted"`
	Blocked       int64           `json:"blocked"`
	Ecosystemname string          `json:"ecosystemname"`
	Token_title   string          `json:"token_title"`
}

type EcosyKeyHex struct {
	Ecosystem int64 `json:"ecosystem"`
	//Ecosyname string `json:"Ecosyname"`
	IsValued        int64           `json:"isvalued"`
	Ecosystemname   string          `json:"ecosystemname"`
	Token_title     string          `json:"token_title"`
	Amount          decimal.Decimal `json:"amount"`
	Info            string          `json:"info"`
	Emission_amount string          `json:"emission_amount"`
	Type_emission   int64
	Type_withdraw   int64
}

type EcosyKeyTotalHex struct {
	Ecosystem       int64           `json:"ecosystem"`
	IsValued        int64           `json:"isvalued"`
	Ecosystemname   string          `json:"ecosystemname"`
	Token_title     string          `json:"token_title"`
	Amount          decimal.Decimal `json:"amount"`
	Info            string          `json:"info"`
	Emission_amount string          `json:"emission_amount"`
	MemberName      string          `json:"member_name"`
	MerberUrl       string          `json:"member_url"`
	LogoUrl         string          `json:"logo_url"`
	Type_emission   int64
	Type_withdraw   int64
	Transaction     int64           `json:"transaction"`
	StakeAmount     decimal.Decimal `json:"stake_amount"`
	FreezeAmount    decimal.Decimal `json:"freeze_amount"`
	Inamount        decimal.Decimal `json:"inamount"`
	Outamount       decimal.Decimal `json:"outamount"`
	Totalamount     decimal.Decimal `json:"total_amount"`
}

type KeysResult struct {
	Total    int64              `json:"total" `
	Page     int                `json:"page" `
	Limit    int                `json:"limit" `
	SysEcosy EcosyKeyTotalHex   `json:"sysecosy,omitempty"`
	Rets     []EcosyKeyTotalHex `json:"rets" `
}

// SetTablePrefix is setting table prefix
func (m *Key) SetTablePrefix(prefix int64) *Key {
	m.Ecosystem = prefix
	return m
}

// TableName returns name of table
func (m Key) TableName() string {
	if m.Ecosystem == 0 {
		m.Ecosystem = 1
	}
	return `1_keys`
}

// Get is retrieving model from database
func (m *Key) Get(id int64, wallet string) (*EcosyKeyHex, error) {

	var (
		ecosystems []Ecosystem
	)
	da := EcosyKeyHex{}

	key := strconv.FormatInt(id, 10)
	wid, err := strconv.ParseInt(wallet, 10, 64)
	if err == nil {
		//
		err := conf.GetDbConn().Conn().Table("1_ecosystems").Where("id = ?", key).Find(&ecosystems).Error
		if err == nil {
			da.Ecosystem = id
			da.Ecosystemname = ecosystems[0].Name
			da.IsValued = ecosystems[0].IsValued
			if da.Ecosystem == 1 {
				da.Token_title = consts.SysEcosytemTitle
			} else {
				da.Token_title = ecosystems[0].TokenTitle
			}
			//da.Token_title = ecosystems[0].Token_title
		}
		err = conf.GetDbConn().Conn().Table("1_keys").Where("id = ?", wid).Find(m).Error
		if err == nil {
			da.Ecosystem = id
			da.Amount = m.Amount
		}
	}

	//err := conf.GetDbConn().Conn().Where("id = ? and ecosystem = ?", wallet, m.ecosystem).First(m).Error
	return &da, err
}

// Get is retrieving model from database
func (m *Key) GetEcosykey(wallet string) (*[]EcosyKeyHex, error) {

	var (
		keys []Key
		ret  []EcosyKeyHex
	)

	//
	wid, err := strconv.ParseInt(wallet, 10, 64)
	if err == nil {
		//
		err = conf.GetDbConn().Conn().Table("1_keys").Where("id = ?", wid).Find(&keys).Error
		count := len(keys)
		if err == nil && count > 0 {
			for i := 0; i < len(keys); i++ {
				da := EcosyKeyHex{}
				da.Ecosystem = keys[i].Ecosystem
				da.Amount = keys[i].Amount

				var ecosystems Ecosystem
				//
				key := strconv.FormatInt(da.Ecosystem, 10)
				err := conf.GetDbConn().Conn().Table("1_ecosystems").Where("id = ?", key).Find(&ecosystems).Error
				if err == nil {
					//da.Ecosystem = id
					da.Ecosystemname = ecosystems.Name
					if da.Ecosystem == 1 {
						da.Token_title = consts.SysEcosytemTitle
					} else {
						da.Token_title = ecosystems.TokenTitle
					}
					//da.Token_title = ecosystems.Token_title

				}

				ret = append(ret, da)

			}
		}
	}

	return &ret, err
}

//Get is retrieving model from database
func (ts *Key) GetKeys(id int64, page int, size int, order string) (*[]KeyHex, int64, error) {
	var (
		tss    []Key
		ret    []KeyHex
		num    int64
		ioffet int
		i      int
	)

	if order == "" {
		order = "id asc"
	}
	num = 0

	key := strconv.FormatInt(id, 10)
	//err := conf.GetDbConn().Conn().Table(key + "_keys").Order(order).Find(&tss).Error
	err := conf.GetDbConn().Conn().Table("1_keys").Where("ecosystem = ?", key).Order(order).Find(&tss).Error
	if err != nil {
		return &ret, num, err
	}
	if page < 1 || size < 1 {
		return &ret, num, err
	}
	ioffet = (page - 1) * size
	num = int64(len(tss))
	if num < int64(page*size) {
		size = int(num) % size
	}
	if num < int64(ioffet) || num < 1 {
		return &ret, num, err
	}

	es := Ecosystem{}
	f, err := es.Get(id)
	if err != nil {
		return &ret, num, err
	}
	if !f {
		return &ret, num, err
	}
	for i = 0; i < size; i++ {
		da := KeyHex{}
		da.Ecosystem = id
		if da.Ecosystem == 1 {
			da.Token_title = consts.SysEcosytemTitle
		} else {
			da.Token_title = es.TokenTitle
		}
		//da.Token_title = es.Token_title
		da.Ecosystemname = es.Name
		da.ID = strconv.FormatInt(tss[ioffet].ID, 10)
		da.PublicKey = hex.EncodeToString(tss[ioffet].PublicKey)
		da.Maxpay = tss[ioffet].Maxpay
		da.Amount = tss[ioffet].Amount
		da.Deleted = tss[ioffet].Deleted
		da.Multi = tss[ioffet].Multi
		da.Blocked = tss[ioffet].Blocked
		//fmt.Println("ecosystem %d", id)
		ret = append(ret, da)
		ioffet++
	}

	return &ret, num, err
}

func (m *Key) GetTotal(page, limit int, order, wallet string) (int64, int, *[]EcosyKeyTotalHex, error) {

	var (
		tss   []Key
		total int64
	)
	var da []EcosyKeyTotalHex

	wid := converter.StringToAddress(wallet)
	err := errors.New("wallet err ")
	//wid, err := strconv.ParseInt(wallet, 10, 64)
	if wid != 0 {
		err = conf.GetDbConn().Conn().Table("1_keys").
			Where("id = ?", wid).
			Count(&total).Error
		if err != nil {
			return 0, 0, &da, err
		}
		err = conf.GetDbConn().Conn().Table("1_keys").Where("id = ?", wid).Order(order).Offset((page - 1) * limit).Limit(limit).Find(&tss).Error
		if err == nil {
			dlen := len(tss)
			for i := 0; i < dlen; i++ {
				ds := tss[i]
				d := EcosyKeyTotalHex{}
				d.Ecosystem = ds.Ecosystem
				d.Amount = ds.Amount

				//
				ems := Ecosystem{}
				f, err := ems.Get(ds.Ecosystem)
				if err != nil {
					return 0, 0, &da, err
				}
				if f {
					d.Ecosystemname = ems.Name
					d.IsValued = ems.IsValued
					if d.Ecosystem == 1 {
						d.Token_title = consts.SysEcosytemTitle
					} else {
						d.Token_title = ems.TokenTitle
					}
				}
				//
				ts := &History{}
				dh, err := ts.GetWalletHistoryTotals(ds.Ecosystem, wallet)
				if err != nil {
					return 0, 0, &da, err
				}
				d.Transaction = dh.Transaction
				d.Inamount = dh.Inamount
				d.Outamount = dh.Outamount

				da = append(da, d)
			}

			return total, limit, &da, nil

		}
	}

	return 0, 0, &da, err
}

func (m *Key) GetEcosyKey(keyid int64, wallet, logourl string) (*EcosyKeyTotalHex, error) {

	d := EcosyKeyTotalHex{}
	d.Ecosystem = m.Ecosystem
	d.Amount = m.Amount

	mb := Member{}
	fm, _ := mb.GetAccount(m.Ecosystem, wallet)
	if fm {
		if mb.ImageID != nil {
			if *mb.ImageID != int64(0) {
				mrl, err := Loadlogo(*mb.ImageID)
				if err != nil {
					return &d, err
				}
				if mrl != "" {
					d.MerberUrl = logourl + mrl
				}
			}
		}
		d.MemberName = mb.MemberName
	}

	escape := func(value interface{}) string {
		return strings.Replace(fmt.Sprint(value), `'`, `''`, -1)
	}

	//
	ems := Ecosystem{}
	f, err := ems.Get(m.Ecosystem)
	if err != nil {
		return &d, err
	}
	if f {
		d.Ecosystemname = ems.Name
		d.IsValued = ems.IsValued

		if ems.Info != "" {
			minfo := make(map[string]interface{})
			err := json.Unmarshal([]byte(ems.Info), &minfo)
			if err != nil {
				return &d, err
			}
			usid, ok := minfo["logo"]
			if ok {
				urid := escape(usid)
				uid, err := strconv.ParseInt(urid, 10, 64)
				if err != nil {
					return &d, err
				}

				url, err := Loadlogo(uid)
				if err != nil {
					return &d, err
				}
				if url != "" {
					d.LogoUrl = logourl + url
				}

			}
		}
		if d.Ecosystem == 1 {
			d.Token_title = consts.SysEcosytemTitle
			if d.Ecosystemname == "" {
				d.Ecosystemname = "platform ecosystem"
			}
		} else {
			d.Token_title = ems.TokenTitle
		}
	}
	//
	ts := &History{}
	dh, err := ts.GetWalletHistoryTotalbykeyid(m.Ecosystem, keyid)
	if err != nil {
		return &d, err
	}

	ag := &AssignGetInfo{}
	ba, fa, _, err := ag.GetBalance(nil, keyid)
	if err != nil {
		return &d, err
	}
	d.Transaction = dh.Transaction
	d.Inamount = dh.Inamount
	d.Outamount = dh.Outamount
	d.StakeAmount = m.MineLock.Add(m.PoolLock)
	if ba {
		d.FreezeAmount = fa
	}

	d.Totalamount = d.Totalamount.Add(d.Amount)
	d.Totalamount = d.Totalamount.Add(d.StakeAmount)
	d.Totalamount = d.Totalamount.Add(d.FreezeAmount)
	return &d, err
}

// Get is retrieving model from database
func (m *Key) GetWalletTotal(page, limit int, order string, wallet, logourl string) (*KeysResult, error) {

	var (
		tss   []Key
		ft    Key
		total int64
		ret   KeysResult
	)
	ret.Limit = limit
	ret.Page = page
	da := []EcosyKeyTotalHex{}

	wid := converter.StringToAddress(wallet)
	err := errors.New("wallet err ")
	if wallet == "0000-0000-0000-0000-0000" {
		d := EcosyKeyTotalHex{}
		ret.SysEcosy = d
		ret.Total = 1
		ret.Rets = da
		return &ret, nil
	} else if wid != 0 {

		err = conf.GetDbConn().Conn().Table("1_keys").
			Where("id = ?", wid).
			Count(&total).Error
		if err != nil {
			return &ret, err
		}

		err := conf.GetDbConn().Conn().Table("1_keys").Where("id = ? and ecosystem = ?", wid, 1).First(&ft).Error
		if err != nil {
			return &ret, err
		}
		df, err := ft.GetEcosyKey(wid, wallet, logourl)
		if err != nil {
			return &ret, err
		}
		ret.SysEcosy = *df

		err = conf.GetDbConn().Conn().Table("1_keys").Where("id = ? and ecosystem != ?", wid, 1).Order(order).Offset((page - 1) * limit).Limit(limit).Find(&tss).Error
		if err == nil {
			dlen := len(tss)
			for i := 0; i < dlen; i++ {
				ds := tss[i]
				d, err := ds.GetEcosyKey(wid, wallet, logourl)
				if err != nil {
					return &ret, err
				}
				da = append(da, *d)
			}
			ret.Total = int64(total)
			ret.Rets = da
			return &ret, nil

		}
	}
	var res result
	err := conf.GetDbConn().Conn().Table("1_keys").Select("SUM(amount) as amount").Where("ecosystem = 1").Scan(&res).Error
	return res.Amount, err
}

func (m *Key) GetStakeAmount() (string, error) {
	type result struct {
		Amount decimal.Decimal
	}
	var agi AssignGetInfo
	agm, err := agi.GetAllBalance(nil)
	if err != nil {
		return "0", err
	}

	if HasTableOrView(nil, "1_mine_stake") {
		var mine, pool result
		err = conf.GetDbConn().Conn().Table("1_keys").Select("SUM(mine_lock) as amount").Where("ecosystem = 1").Scan(&mine).Error
		if err != nil {
			b := strings.ContainsAny(err.Error(), "column mine_lock does not exist")
			if b {
				return agm.String(), nil
			}
			return "0", err
		}
		err = conf.GetDbConn().Conn().Table("1_keys").Select("SUM(pool_lock) as amount").Where("ecosystem = 1").Scan(&pool).Error
		if err != nil {
			b := strings.ContainsAny(err.Error(), "column pool_lock does not exist")
			if b {
				return agm.String(), nil
			}
			return "0", err
		}

		rt := mine.Amount.Add(pool.Amount)
		rt = rt.Add(agm)
		return rt.String(), nil
	}
	return agm.String(), nil
}
