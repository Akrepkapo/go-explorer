/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import (
	"encoding/hex"
	"encoding/json"
	"github.com/IBAX-io/go-explorer/conf"
	"github.com/IBAX-io/go-explorer/consts"
	log "github.com/sirupsen/logrus"
	"strconv"
	"unsafe"
)

// LogTransaction is model
type LogTransaction struct {
	Hash  []byte `gorm:"primary_key;not null"`
	Block int64  `gorm:"not null"`
}

var (
	GLogTranHash map[string]int64
)

func (lt *LogTransaction) GetBlockTransactionsDashboard(page int, size int) (*[]BlockTxDetailedInfoDashboard, error) {
	var (
		tss    []LogTransaction
		ret    []BlockTxDetailedInfoDashboard
		num    int
		ioffet int
		i      int
		j      int32
		err    error
	)
	ioffet = (page - 1) * size
	if page < 1 || size < 1 {
		return &ret, err
	}
	err = GetDB(nil).Offset(ioffet).Limit(size).Order("block desc").Find(&tss).Error
	//err = DBConn.Order(order).Find(&tss).Error
	if err != nil {
		return &ret, err
	}
	num = len(tss)
	if num < page*size {
		size = num % size
	}

	if num < ioffet || num < 1 {
		return &ret, err
	}
		found, err := bk.Get(Blocks[i-1])
		if err == nil && found {
			rt, err := GetBlocksDetailedInfoHex(bk)
			if err != nil {
				return nil, err
			}
			for j = 0; j < rt.Tx; j++ {
				bh := BlockTxDetailedInfoDashboard{}
				bh.BlockID = rt.Header.BlockID
				bh.ContractName = rt.Transactions[j].ContractName
				bh.Hash = rt.Transactions[j].Hash
				bh.KeyID = rt.Transactions[j].KeyID
				bh.Time = rt.Transactions[j].Time
				bh.Ecosystemname = rt.Transactions[j].Ecosystemname
				if Thash[rt.Transactions[j].Hash] {
					ret = append(ret, bh)
				}
			}
		} else {
			if err != nil {
				return nil, err
			}
		}
	}
	return &ret, err
}

func (lt *LogTransaction) Get_BlockTransactions(page int, size int, order string) (*[]BlockTxDetailedInfoHex, int, error) {
	var (
		tss    []LogTransaction
		ret    []BlockTxDetailedInfoHex
		num    int
		ioffet int
		i      int
		j      int32
		err    error
	)
	if page < 1 || size < 1 {
		return &ret, num, err
	}
	err = GetDB(nil).Order("block desc").Find(&tss).Error
	//err = DBConn.Order(order).Find(&tss).Error
	if err != nil {
		return &ret, num, err
	}
	num = len(tss)
	ioffet = (page - 1) * size
	if num < page*size {
		size = num % size
	}

	if num < ioffet || num < 1 {
		return &ret, num, err
	}
	TBlock := make(map[string]int64)
	Thash := make(map[string]bool)
	for i = 0; i < size; i++ {
		hash := hex.EncodeToString(tss[ioffet+i].Hash)
		Thash[hash] = true

		key := strconv.FormatInt(tss[ioffet+i].Block, 10)
		TBlock[key] = tss[ioffet+i].Block
	}

	var Blocks []int64
	for _, k := range TBlock {
		Blocks = append(Blocks, k)
	}

	quickSort(Blocks, 0, int64(len(Blocks)-1))

	for i = len(Blocks); i > 0; i-- {
		bk := &Block{}
		found, err := bk.Get(Blocks[i-1])
		if err == nil && found {
			rt, err := GetBlocksDetailedInfoHex(bk)
			if err != nil {
				return nil, 0, err
			}
			for j = 0; j < rt.Tx; j++ {
				bh := BlockTxDetailedInfoHex{}
				bh.BlockID = rt.Header.BlockID
				bh.ContractName = rt.Transactions[j].ContractName
				bh.Hash = rt.Transactions[j].Hash
				bh.KeyID = rt.Transactions[j].KeyID
				bh.Params = rt.Transactions[j].Params
				bh.Time = rt.Transactions[j].Time
				bh.Type = rt.Transactions[j].Type
				bh.Ecosystem = rt.Transactions[j].Ecosystem
				bh.Ecosystemname = rt.Transactions[j].Ecosystemname
				if bh.Ecosystem == 1 {
					bh.Token_title = consts.SysEcosytemTitle
					if bh.Ecosystemname == "" {
						bh.Ecosystemname = "platform ecosystem"
					}
				} else {
					bh.Token_title = rt.Transactions[j].Token_title
				}
				Ten := unsafe.Sizeof(rt.Transactions[j])
				bh.Size = int64(Ten)
				if Thash[rt.Transactions[j].Hash] {
					ret = append(ret, bh)
				}
			}
		} else {
			if err != nil {
				return nil, 0, err
			}
		}
	}
	return &ret, num, err
}

func SetTransactionBlockLastToRedis(blockid int64) {
	ts := &LogTransaction{}
	ret, _, err := ts.Get_BlockTransactionsLast(blockid, 1, 5, "block desc")
	blockInfo, err := json.Marshal(ret)
	if err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("SyncBlockinfoToRedis json err")
		return
	}
	value, err := GzipEncode(blockInfo)
	if err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("SyncBlockinfoToRedis GzipEncode err")
		return
	}
	rd := RedisParams{
		Key:   "transaction-block-last",
		Value: string(value),
	}
	if err := rd.Set(); err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("SyncBlockinfoToRedis setdb err")
	}
	if err := sendTransactionBlockLastToWebsocket(ret); err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("sendBlockListToWebsocket err")
	}
}

func (lt *LogTransaction) GetTransactionBlockLastFromRedis() ([]BlockTxDetailedInfoHex, error) {
	var err error
	var blockTx []BlockTxDetailedInfoHex
	rd := RedisParams{
		Key:   "transaction-block-last",
		Value: "",
	}
	if err := rd.Get(); err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("GetBlocksListFromRedis getdb err")
		return nil, err
	}
	value, err1 := GzipDecode([]byte(rd.Value))
	if err1 != nil {
		log.WithFields(log.Fields{"warn": err1}).Warn("GetBlocksListFromRedis GzipDecode err")
		return nil, err1
	}
	if err = json.Unmarshal(value, &blockTx); err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("GetBlocksListFromRedis json err")
		return nil, err
	}
	return blockTx, err
}
func sendTransactionBlockLastToWebsocket(toptransactions interface{}) error {
	dat := ResponseTopTitle{}
	dat.Cmd = CmdTopTransactions
	dat.Data = toptransactions

	ds, err := json.Marshal(dat)
	if err != nil {
		return err
	}
	err = WriteChannelByte(ChannelTopData, ds)
	if err != nil {
		return err
	}
	return nil
}
func (lt *LogTransaction) Get_BlockTransactionsLast(id int64, page int64, size int64, order string) (*[]BlockTxDetailedInfoHex, int64, error) {
	var (
		tss []LogTransaction
		ret []BlockTxDetailedInfoHex
		num int64
		//ioffet int64
		i int64
		j int32
		//err    error
	)

	num = 0
	if page < 1 || size < 1 {
		return &ret, num, nil
	}
	err := GetDB(nil).Limit(int(size)).Offset(int((page-1)*size)).Order(order).Where("block <=?", id).Find(&tss).Error
	//err := DBConn.Order("block desc").Where("block <=?", id).Find(&tss).Error
	if err != nil {
		return &ret, num, err
	}
	//fmt.Println("page:%d, size:%d", page, size)
	//if page < 1 || size < 1 {
	//	return &ret, num, err
	//}
	//num = int64(len(tss))
	////fmt.Println("tr_blocks Error: %d", num)
	//ioffet = (page - 1) * size
	//if num < page*size {
	//	size = num % size
	//	//ioffet = 0
	//}
	//
	//if num < ioffet || num < 1 {
	//	return &ret, num, err
	//}
	TBlock := make(map[string]int64)
	Thash := make(map[string]bool)
	if len(tss) >= int(size) {
		for i = 0; i < size; i++ {
			hash := hex.EncodeToString(tss[i].Hash)
			Thash[hash] = true

			key := strconv.FormatInt(tss[i].Block, 10)
			TBlock[key] = tss[i].Block
		}
	} else {
		for k := 0; k < len(tss); k++ {
			hash := hex.EncodeToString(tss[k].Hash)
			Thash[hash] = true

			key := strconv.FormatInt(tss[k].Block, 10)
			TBlock[key] = tss[k].Block
		}
	}

	var Blocks []int64
	for _, k := range TBlock {
		Blocks = append(Blocks, k)
	}

	quickSort(Blocks, 0, int64(len(Blocks)-1))

	for i = int64(len(Blocks)); i > 0; i-- {
		//seelog.Info("logtran GetBlocksDetailedInfoHex i:", Blocks[i-1])
		bk := &Block{}
		found, err := bk.Get(Blocks[i-1])
		if err == nil && found {
			rt, err := GetBlocksDetailedInfoHex(bk)
			if err == nil {
				for j = 0; j < rt.Tx; j++ {
					bh := BlockTxDetailedInfoHex{}
					bh.BlockID = rt.Header.BlockID
					bh.ContractName = rt.Transactions[j].ContractName
					bh.Hash = rt.Transactions[j].Hash
					bh.KeyID = rt.Transactions[j].KeyID
					bh.Params = rt.Transactions[j].Params
					bh.Time = rt.Transactions[j].Time
					bh.Type = rt.Transactions[j].Type

					bh.Ecosystem = rt.Transactions[j].Ecosystem
					bh.Ecosystemname = rt.Transactions[j].Ecosystemname
					if bh.Ecosystem == 1 {
						bh.Token_title = consts.SysEcosytemTitle
					} else {
						bh.Token_title = rt.Transactions[j].Token_title
					}
					//bh.Token_title = rt.Transactions[j].Token_title

					Ten := unsafe.Sizeof(rt.Transactions[j])
					bh.Size = int64(Ten)

					if Thash[rt.Transactions[j].Hash] {
						ret = append(ret, bh)
					} else {
						//seelog.Info("logtran ihash: %s ", hash)
						log.WithFields(log.Fields{"hash": rt.Transactions[j].Hash}).Warn("logtran jhash:")
					}

					//seelog.Info("logtran GetBlocksDetailedInfoHex i: %d  j:%d ", i, j)
					//seelog.Flush()
				}
			} else {
				//
				log.WithFields(log.Fields{"warn": err.Error()}).Warn("logtran GetBlocksDetailedInfoHex")
			}

		} else {
			if err != nil {
				log.WithFields(log.Fields{"warn": err.Error()}).Warn("logtran GetBlocks  DetailedInfoHex")
			}
		}
	}

	return &ret, num, err
}

func quickSort(arr []int64, start, end int64) {
	if start < end {
		i, j := start, end
		key := arr[(start+end)/2]
		for i <= j {
			for arr[i] < key {
				i++
			}
			for arr[j] > key {
				j--
			}
			if i <= j {
				arr[i], arr[j] = arr[j], arr[i]
				i++
				j--
			}
		}

		if start < j {
			quickSort(arr, start, j)
		}
		if end > i {
			quickSort(arr, i, end)
		}
	}
}

func getTransactionBlockToRedis() error {
	ts := &LogTransaction{}
	ret, err := ts.GetBlockTransactionsDashboard(1, 10)
	if err != nil {
		return err
	}
	value, err := json.Marshal(ret)
	if err != nil {
		return err
	}
	rd := RedisParams{
		Key:   "transaction-block",
		Value: string(value),
	}
	if err := rd.Set(); err != nil {
		return err
	}
	return nil
}

func GetTransactionBlockFromRedis() (*[]BlockTxDetailedInfoHex, int, error) {
	var ret []BlockTxDetailedInfoHex
	rd := RedisParams{
		Key:   "transaction-block",
		Value: "",
	}
	if err := rd.Get(); err != nil {
		return nil, 0, err
	}
	if err := json.Unmarshal([]byte(rd.Value), &ret); err != nil {
		return nil, 0, err
	}
	return &ret, len(ret), nil
}

func (lt *LogTransaction) getTransactionIdFromHash(hash []byte) (bool, error) {
	f, err := isFound(conf.GetDbConn().Conn().Where("hash = ?", hash).First(&lt))
	if f && err == nil {
		return f, err
	}
	return f, err
}
