/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/IBAX-io/go-explorer/consts"

	"github.com/IBAX-io/go-explorer/conf"

	"github.com/IBAX-io/go-ibax/packages/converter"
	"github.com/sirupsen/logrus"
	"github.com/vmihailenco/msgpack/v5"
)

type BlockTxDetailedInfo struct {
	BlockID      int64  `gorm:"not null;index:btdblockid_idx" json:"block_id"`
	Hash         string `gorm:"primary_key;not null" json:"hash"`
	ContractName string `gorm:"not null" json:"contract_name"`
	Params       string `gorm:"not null" json:"params"`
	KeyID        string `gorm:"not null" json:"key_id"`
	Time         int64  `gorm:"not null" json:"time"`
	Type         int64  `gorm:"not null" json:"type"`
	Size         int64  `gorm:"not null" json:"size"`

	Ecosystemname string `gorm:"null" json:"ecosystemname"`
	Token_title   string `gorm:"null" json:"token_title"`
	Ecosystem     int64  `gorm:"null" json:"ecosystem"`
}

type BlockTxDetailedInfoHex struct {
	BlockID      int64  `gorm:"not null;index:blockid_idx" json:"block_id"`
	Hash         string `gorm:"primary_key;not null" json:"hash"`
	ContractName string `gorm:"not null" json:"contract_name"`
	//Params       map[string]interface{} `json:"params"`
	Params string `gorm:"not null" json:"params"`
	KeyID  string `gorm:"not null" json:"key_id"`
	Time   int64  `gorm:"not null" json:"time"`
	Type   int64  `gorm:"not null" json:"type"`
	Size   int64  `gorm:"not null" json:"size"`

	Ecosystemname string `gorm:"null" json:"ecosystemname"`
	Token_title   string `gorm:"null" json:"token_title"`
	Ecosystem     int64  `gorm:"null" json:"ecosystem"`
}

type BlockTxDetailedInfoDashboard struct {
	BlockID       int64  `gorm:"not null;index:blockid_idx" json:"block_id"`
	Hash          string `gorm:"primary_key;not null" json:"hash"`
	ContractName  string `gorm:"not null" json:"contract_name"`
	KeyID         string `gorm:"not null" json:"key_id"`
	Time          int64  `gorm:"not null" json:"time"`
	Ecosystemname string `gorm:"null" json:"ecosystemname"`
}

type HashTransactionResult struct {
	Total int64                    `json:"total" `
	Page  int                      `json:"page" `
	Limit int                      `json:"limit"`
	Rets  []BlockTxDetailedInfoHex `json:"rets"`
}

func (bt *BlockTxDetailedInfoHex) GetByHash_Db(hash string) (bool, error) {
	return bt.GetDb_txdetailedHash(hash)
}

func (bt *BlockTxDetailedInfoHex) GetByBlockid_BlockTransactionsLast_DB(id int64, page int, limit int, order string) (int64, *[]BlockTxDetailedInfoHex, error) {
	var (
		ret   []BlockTxDetailedInfoHex
		total int64
		err   error
	)
	ret, err = bt.GetDb_txdetailedId(id)
	if err != nil {
		return total, &ret, err
	}
	//if len(ret) > limit {
	//	ret = ret[:limit]
	//}
	total = int64(len(ret))
	return total, &ret, err
}

func (bt *BlockTxDetailedInfoHex) GetByKeyid_BlockTransactionsLast_Db(id string, page int, size int, order string) (int64, *[]BlockTxDetailedInfoHex, error) {
	var (
		ret   []BlockTxDetailedInfoHex
		total int64
		err   error
	)
	ret, total, err = bt.GetDb_txdetailedKey(id, order, size, page)
	if err != nil {
		return total, &ret, err
	}
	//total = int64(len(ret))
	return total, &ret, err

}

//Get is retrieving model from database
func (bt *BlockTxDetailedInfoHex) GetCommonTransactionSearch(page, limit int, search, order string) (*HashTransactionResult, error) {
	var (
		ret HashTransactionResult
		err error
	)
	ret.Page = page
	ret.Limit = limit

	bid, err := strconv.ParseInt(search, 10, 64)
	if err == nil && bid > 0 {
		//blockid
		total, rets, err := bt.GetByBlockid_BlockTransactionsLast_DB(bid, page, limit, order)
		if err != nil {
			return &ret, err
		}
		ret.Total = total
		ret.Rets = *rets
		return &ret, err
	} else {
		keyid := converter.StringToAddress(search)
		if keyid != 0 {
			//wallet

			total, rets, err := bt.GetByKeyid_BlockTransactionsLast_Db(search, page, limit, order)
			if err != nil {
				return &ret, err
			}
			ret.Total = total
			ret.Rets = *rets
			return &ret, err
		} else {
			//hash
			if search == "" {
				if page == 1 && limit == 10 {
					rets, total, err := GetTransactionBlockFromRedis()
					if err != nil {
						return &ret, err
					}
					ret.Total = int64(total)
					ret.Rets = *rets
					return &ret, err
				}
				rets, total, err := Get_Group_TransactionBlock(page, limit, order)
				//total, rets, err := bt.GetFind_BlockTransactionsLast_Sqlite(page, limit, order)
				if err != nil {
					return &ret, err
				}
				ret.Total = total
				ret.Rets = *rets
				return &ret, err
			} else {
				f, err := bt.GetByHash_Db(search)
				if err != nil {
					return &ret, err
				}
				if f {
					ret.Total = 1
					ret.Rets = append(ret.Rets, *bt)
					return &ret, err
				} else {
					return &ret, errors.New("not found hash")
				}
			}
		}
	}
}
func Get_Group_TransactionBlock(ids int, icount int, order string) (*[]BlockTxDetailedInfoHex, int64, error) {
	ts := &LogTransaction{}
	//bt := &models.BlockTxDetailedInfoHex{}

	//if models.GsqliteIsactive {
	//	ret, num, err := bt.Get_BlockTransactions_Sqlite(ids, icount, order)
	//	if err == nil && *ret != nil && num > 0 {
	//		//fmt.Println("Get_BlockTransactions_Sqlite  ok ids:%d icount:%d", ids, icount)
	//		return ret, num, err
	//	}
	//}
	ret, num, err := ts.Get_BlockTransactions(ids, icount, order)
	//fmt.Println("Get_BlockTransactions pg  ok ids:%d icount:%d", ids, icount)
	return ret, int64(num), err
	//return nil, 0, nil
}
func (bt *BlockTxDetailedInfoHex) Get_BlockTransactions_Sqlite(page int, size int, order string) (*[]BlockTxDetailedInfoHex, int, error) {
	var (
		ret []BlockTxDetailedInfoHex
		tss []BlockTxDetailedInfo
	)

	err := conf.GetDbConn().Conn().Limit(size).Offset((page - 1) * size).Order(order).Find(&tss).Error
	if err == nil {

		for i := 0; i < len(tss); i++ {

			bh := BlockTxDetailedInfoHex{}
			//params := map[string]interface{}
			bh.BlockID = tss[i].BlockID
			bh.ContractName = tss[i].ContractName
			bh.Hash = tss[i].Hash
			bh.KeyID = tss[i].KeyID
			//bh.Params = rt.Transactions[j].Params
			bh.Time = tss[i].Time
			bh.Type = tss[i].Type
			bh.Size = tss[i].Size

			bh.Ecosystem = tss[i].Ecosystem
			bh.Ecosystemname = tss[i].Ecosystemname
			if bh.Ecosystem == 1 {
				bh.Token_title = consts.SysEcosytemTitle
			} else {
				bh.Token_title = tss[i].Token_title
			}
			//bh.Token_title = tss[i].Token_title
			//es := Ecosystem{}
			//f, err := es.Get(tss[i]..EcosystemID)
			//if f && err == nil {
			//	bh.Ecosystem = tss[i].TxHeader.EcosystemID
			//	bh.Ecosystemname = es.Name
			//	bh.Token_title = es.Token_title
			//}

			if err := json.Unmarshal([]byte(tss[i].Params), &bh.Params); err == nil {
				//bh.Params
			}
			ret = append(ret, bh)
		}

	}

	return &ret, len(GLogTranHash), err

}

func Deal_LogTransactionBlockTxDetial(objArr *[]LogTransaction) (*[]BlockTxDetailedInfo, error) {
	var (
		ret    []BlockTxDetailedInfo
		Blocks []int64
	)

	ret1 := Deal_Redupliction_LogTransaction(objArr)
	count := len(*ret1)
	if len(*ret1) == 0 {
	}

	for i := int64(len(Blocks)); i > 0; i-- {
		bk := &Block{}
		found, err := bk.Get_Sqlite(Blocks[i-1])
		if err == nil && found {
			rt, err := GetBlocksDetailedInfoHex(bk)
			if err == nil {
				for j := 0; j < int(rt.Tx); j++ {
					bh := BlockTxDetailedInfo{}
					bh.BlockID = rt.Header.BlockID
					bh.ContractName = rt.Transactions[j].ContractName
					bh.Hash = rt.Transactions[j].Hash
					bh.KeyID = rt.Transactions[j].KeyID
					//bh.Params = rt.Transactions[j].Params
					bh.Time = rt.Transactions[j].Time
					bh.Type = rt.Transactions[j].Type

					bh.Ecosystemname = rt.Transactions[j].Ecosystemname
					bh.Ecosystem = rt.Transactions[j].Ecosystem
					if bh.Ecosystem == 1 {
						bh.Token_title = consts.SysEcosytemTitle
					} else {
						bh.Token_title = rt.Transactions[j].Token_title
					}
					bh.Size = rt.Transactions[j].Size

					lg1, err := json.Marshal(rt.Transactions[j].Params)
					if err == nil {
						bh.Params = string(lg1)
					}

					if Thash[rt.Transactions[j].Hash] {
						ret = append(ret, bh)
					}
				}
			} else {
				logrus.Info("logtran GetBlocksDetailedInfoHex: %s", err.Error())
			}

		} else if err != nil {
			logrus.Info("logtran GetBlocks  DetailedInfoHex: %s", err.Error())
		}
	}

	return &ret, nil
}

func Deal_Redupliction_LogTransaction(objArr *[]LogTransaction) *[]LogTransaction {
	var (
		ret []LogTransaction
	)
	if GLogTranHash == nil {
		GLogTranHash = make(map[string]int64)
	}
	for _, val := range *objArr {
		key := hex.EncodeToString(val.Hash)
		dat, ok := GLogTranHash[key]
		if ok {
			logrus.Info("GLogTranHash exist block:%d block1:%d key: "+key, dat, val.Block)
		} else {
			GLogTranHash[key] = val.Block
			ret = append(ret, val)
		}
	}
	return &ret
}

func (s *BlockTxDetailedInfoHex) Marshal() ([]byte, error) {
	if res, err := msgpack.Marshal(s); err != nil {
		return nil, err
	} else {
		return res, err
	}
}

func (s *BlockTxDetailedInfoHex) Unmarshal(bt []byte) error {
	if err := msgpack.Unmarshal(bt, &s); err != nil {
		return err
	}
	return nil
}

func (bt *BlockTxDetailedInfoHex) GetDb_txdetailedHash(hash string) (bool, error) {
	var bk Block
	hs, _ := hex.DecodeString(hash)

	dt := &LogTransaction{}
	fb, err := dt.getTransactionIdFromHash(hs)
	if !fb || err != nil {
		return false, err
	}
	fb, err = bk.Get(dt.Block)
	if err != nil {
		return false, err
	}
	if fb {
		_, _, bkts, err := Deal_TransactionBlockTxDetial(&bk)
		if err != nil {
			return false, err
		}
		for _, obj := range *bkts {
			if obj.Hash == hash {
				val, err := obj.Marshal()
				if err != nil {
					return false, err
				}
				err1 := bt.Unmarshal(val)
				if err1 != nil {
					return false, err1
				}
			}
		}
		if bt == nil {
			return false, nil
		}
	} else {
		return false, nil
	}
	return true, nil
}

func (bt *BlockTxDetailedInfoHex) GetDb_txdetailedId(id int64) ([]BlockTxDetailedInfoHex, error) {
	var bk Block
	var ret []BlockTxDetailedInfoHex
	fb, err := bk.Get(id)
	if err != nil {
		return nil, err
	}
	if fb {
		_, _, bkts, err := Deal_TransactionBlockTxDetial(&bk)
		if err != nil {
			return nil, err
		}
		ret = make([]BlockTxDetailedInfoHex, len(*bkts))
		for i, obj := range *bkts {
			ret[i] = obj
		}
	} else {
		return nil, nil
	}

	return ret, nil
}

func (bt *BlockTxDetailedInfoHex) GetDb_txdetailedKey(key string, order string, limit, page int) ([]BlockTxDetailedInfoHex, int64, error) {
	var bk Block
	var ret []BlockTxDetailedInfoHex
	var needBlock []Block
	var total int64
	fb, err := bk.GetBlocksKey(converter.StringToAddress(key), order)
	if err != nil {
		return nil, total, err
	}
	ioffet := (page - 1) * limit
	var txData int32
	isAddOneBlock := false
	findFirstBlock := false
	var startTxId int
	var endTxiD int
	for i := 0; i < len(fb); i++ {
		txData += fb[i].Tx
		if int(txData) > ioffet {
			if int(txData) < ioffet+limit {
				needBlock = append(needBlock, fb[i])
				if !findFirstBlock {
					findFirstBlock = true
					startTxId = int(txData-(txData-fb[i].Tx)) - (int(txData) - ioffet)
				}
			} else {
				if !isAddOneBlock {
					isAddOneBlock = true
					needBlock = append(needBlock, fb[i])
					endTxiD = (int(txData) - (ioffet + limit))
					if !findFirstBlock {
						findFirstBlock = true
						startTxId = int(txData-(txData-fb[i].Tx)) - (int(txData) - ioffet)
					}
				}
			}
		}
	}
	total = int64(txData)

	for i := 0; i < len(needBlock); i++ {
		_, _, bkts, err := Deal_TransactionBlockTxDetial(&needBlock[i])
		if err != nil {
			return nil, total, err
		}
		for _, obj := range *bkts {
			ret = append(ret, obj)
		}
	}
	if len(ret) >= endTxiD && len(ret) >= startTxId && (len(ret)-endTxiD) >= startTxId {
		ret = ret[startTxId : len(ret)-endTxiD]
	}
	return ret, total, nil
}
