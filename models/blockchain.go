/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package models

import (
	"bytes"
	"compress/gzip"
	"encoding/hex"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"

	"github.com/IBAX-io/go-explorer/consts"

	"github.com/IBAX-io/go-explorer/conf"

	"github.com/IBAX-io/go-ibax/packages/converter"

	//"strconv"
	"strings"
	"time"

	"github.com/IBAX-io/go-ibax/packages/block"
)

// Block is model
type Block struct {
	ID            int64  `gorm:"primary_key;not_null"`
	Hash          []byte `gorm:"not null"`
	RollbacksHash []byte `gorm:"not null"`
	Data          []byte `gorm:"not null"`
	EcosystemID   int64  `gorm:"not null"`
	KeyID         int64  `gorm:"not null"`
	NodePosition  int64  `gorm:"not null"`
	Time          int64  `gorm:"not null"`
	Tx            int32  `gorm:"not null"`
}

// TableName returns name of table
func (Block) TableName() string {
	return "block_chain"
}

func (b *Block) Get(blockID int64) (bool, error) {
	//f, err := b.GetRedisByid(blockID)
	//if f && err == nil {
	//	return f, err
	//}

	f, err := isFound(conf.GetDbConn().Conn().Where("id = ?", blockID).First(b))
	if f && err == nil {
		return f, err
	}
	return f, err
}

// GetNodeBlocksAtTime returns records of blocks for time interval and position of node
func (b *Block) GetBlocksHash(hash []byte) (bool, error) {
	//f, err := b.GetRedisByhash(hash)
	//if f && err == nil {
	//	fmt.Println("return redis !!!\n")
	//	return f, err
	//}

	f, err := isFound(conf.GetDbConn().Conn().Where("hash = ?", hash).First(&b))
	if f && err == nil {
		//b.InsertRedis()
		return f, err
	}
	return f, err
}

//
func (b *Block) GetBlocksKey(key int64, order string) ([]Block, error) {
	var err error
	var blockchain []Block

	err = conf.GetDbConn().Conn().Order(order).Where("key_id = ?", key).Find(&blockchain).Error
	return blockchain, err
}

// GetMaxBlock returns last block existence
func (b *Block) GetMaxBlock() (bool, error) {
	return isFound(conf.GetDbConn().Conn().Last(b))
}

// GetBlockchain is retrieving chain of blocks from database
func GetBlockchain(startBlockID int64, endblockID int64, order string) (*[]Block, error) {
	var err error
	blockchain := new([]Block)

	orderStr := "id " + string(order)
	query := conf.GetDbConn().Conn().Model(&Block{}).Order(orderStr)
	if endblockID > 0 {
		query = query.Where("id > ? AND id <= ?", startBlockID, endblockID).Find(&blockchain)
	} else {
		query = query.Where("id > ?", startBlockID).Find(&blockchain)
	}

	if query.Error != nil {
		return nil, err
	}
	return blockchain, nil
}

// GetBlocks is retrieving limited chain of blocks from database
func (b *Block) GetBlocks(startFromID int, limit int, order string) (*[]Block, error) {
	var err error
	blockchain := new([]Block)

	if startFromID > 0 {
		err = conf.GetDbConn().Conn().Limit(limit + 1).Offset((startFromID - 1) * limit).Order(order).Find(blockchain).Error
		//err = conf.GetDbConn().Conn().Order("id desc").Limit(limit).Where("id > ?", startFromID).Find(&blockchain).Error
	} else {
		err = conf.GetDbConn().Conn().Order("id desc").Limit(limit).Find(blockchain).Error
	}
	return blockchain, err
}

// GetBlocksFrom is retrieving ordered chain of blocks from database
func (b *Block) GetBlocksFrom(startFromID int64, ordering string, limit int) ([]Block, error) {
	var err error
	blockchain := new([]Block)

	if limit == 0 {
		err = conf.GetDbConn().Conn().Order("id "+ordering).Where("id > ?", startFromID).Find(&blockchain).Error
	} else {
		err = conf.GetDbConn().Conn().Order("id "+ordering).Where("id > ?", startFromID).Limit(limit).Find(&blockchain).Error
	}
	return *blockchain, err
}

// GetReverseBlockchain returns records of blocks in reverse ordering
func (b *Block) GetReverseBlockchain(endBlockID int64, limit int) ([]Block, error) {
	var err error
	blockchain := new([]Block)

	err = conf.GetDbConn().Conn().Model(&Block{}).Order("id DESC").Where("id <= ?", endBlockID).Limit(limit).Find(&blockchain).Error
	return *blockchain, err
}

// GetNodeBlocksAtTime returns records of blocks for time interval and position of node
func (b *Block) GetNodeBlocksAtTime(from, to time.Time, node int64) ([]Block, error) {
	var err error
	blockchain := new([]Block)

	err = conf.GetDbConn().Conn().Model(&Block{}).Where("node_position = ? AND time BETWEEN ? AND ?", node, from.Unix(), to.Unix()).Find(&blockchain).Error
	return *blockchain, err
}

// Get is retrieving model from database
func (b *Block) Get_Sqlite(blockID int64) (bool, error) {
	return isFound(conf.GetDbConn().Conn().Where("id = ?", blockID).First(b))
}

// GetMaxBlock returns last block existence
func (b *Block) GetMaxBlock_Sqlite() (bool, error) {
	return isFound(conf.GetDbConn().Conn().Last(b))
}

// GetBlockchain is retrieving chain of blocks from database
func GetBlockchain_Sqlite(startBlockID int64, endblockID int64, order string) ([]Block, error) {
	var err error
	blockchain := new([]Block)

	orderStr := "id " + string(order)
	query := conf.GetDbConn().Conn().Model(&Block{}).Order(orderStr)
	if endblockID > 0 {
		query = query.Where("id > ? AND id <= ?", startBlockID, endblockID).Find(&blockchain)
	} else {
		query = query.Where("id > ?", startBlockID).Find(&blockchain)
	}

	if query.Error != nil {
		return nil, err
	}
	return *blockchain, nil
}

// GetBlocks is retrieving limited chain of blocks from database
func (b *Block) GetBlocks_Sqlite(startFromID int, limit int, order string) (*[]Block, error) {
	var (
		err        error
		blockchain []Block
	)
	//blockchain := new([]Block)
	if startFromID >= 0 {
		//err = conf.GetDbConn().Conn().Order("id desc").Limit(limit).Where("id > ?", startFromID).Find(&blockchain).Error
		//err = conf.GetDbConn().Conn().Order(order).Where("id > ? and id<= ?", startFromID, startFromID+limit).Find(&blockchain).Error
		err = conf.GetDbConn().Conn().Limit(limit + 1).Offset((startFromID - 1) * limit).Order(order).Find(&blockchain).Error
	} else {
		err = conf.GetDbConn().Conn().Order(order).Limit(limit).Find(&blockchain).Error
	}
	return &blockchain, err
}

//func (b *Block) GetBlocks_Sqlite(startFromID int64, limit int64, order string) (*[]Block, error) {
//	var (
//		err        error
//		blockchain []Block
//	)
//	//blockchain := new([]Block)
//	if startFromID >= 0 {
//		//err = conf.GetDbConn().Conn().Order("id desc").Limit(limit).Where("id > ?", startFromID).Find(&blockchain).Error
//		//err = conf.GetDbConn().Conn().Order(order).Where("id > ? and id<= ?", startFromID, startFromID+limit).Find(&blockchain).Error
//		err = conf.GetDbConn().Conn().Order(order).Offset(startFromID).Limit(limit).Find(&blockchain).Error
//	} else {
//		err = conf.GetDbConn().Conn().Order(order).Limit(limit).Find(&blockchain).Error
//	}
//	return &blockchain, err
//}

// GetBlocksFrom is retrieving ordered chain of blocks from database
func (b *Block) GetBlocksFrom_Sqlite(startFromID int64, ordering string, limit int) ([]Block, error) {
	var err error
	blockchain := new([]Block)
	if limit == 0 {
		err = conf.GetDbConn().Conn().Order("id "+ordering).Where("id > ?", startFromID).Find(&blockchain).Error
	} else {
		err = conf.GetDbConn().Conn().Order("id "+ordering).Where("id > ?", startFromID).Limit(limit).Find(&blockchain).Error
	}
	return *blockchain, err
}

// GetReverseBlockchain returns records of blocks in reverse ordering
func (b *Block) GetReverseBlockchain_Sqlite(endBlockID int64, limit int) ([]Block, error) {
	var err error
	blockchain := new([]Block)
	err = conf.GetDbConn().Conn().Model(&Block{}).Order("id DESC").Where("id <= ?", endBlockID).Limit(limit).Find(&blockchain).Error
	return *blockchain, err
}

// GetNodeBlocksAtTime returns records of blocks for time interval and position of node
func (b *Block) GetNodeBlocksAtTime_Sqlite(from, to time.Time, node int64) ([]Block, error) {
	var err error
	blockchain := new([]Block)
	err = conf.GetDbConn().Conn().Model(&Block{}).Where("node_position = ? AND time BETWEEN ? AND ?", node, from.Unix(), to.Unix()).Find(&blockchain).Error
	return *blockchain, err
}

// GetNodeBlocksAtTime returns records of blocks for time interval and position of node
func (b *Block) GetBlocksHash_Sqlite(hash []byte) (bool, error) {
	return isFound(conf.GetDbConn().Conn().Where("hash = ?", hash).First(b))
}

func GetBlocklist(ids int, icount int, order string) (*[]BlocksResult, int64, error) {
	var (
		ret []BlocksResult
		dat []Block
		i   int64
		err error
	)
	sum := icount
	//err := Errsqlite
	bk := &Block{}
	ret2, err := bk.GetMaxBlock()
	if !ret2 {
		return &ret, bk.ID, err
	}
	if ids < 1 || icount < 1 {
		return &ret, bk.ID, err
	}

	bks := &Block{}
	ret1, err := bks.GetBlocks(ids, icount, order)
	if err == nil {
		dat = *ret1
		for i = 0; i < int64(len(dat)); i++ {
			//ioffet
			//fmt.Println("SendGet detailed_blocks Error: "+strconv.FormatInt(ioffet, 10)+"  "+strconv.FormatInt(icount, 10), err)
			bkx, err := GetBlocksDetailedInfoHex(&dat[i])
			if err == nil {
				var blocks = BlocksResult{
					BlockID:      bkx.Header.BlockID,
					Time:         bkx.Time,
					EcosystemID:  bkx.EcosystemID,
					KeyID:        bkx.KeyID,
					NodePosition: bkx.NodePosition,
					Hash:         bkx.Hash,
					Tx:           bkx.Tx,
				}
				if i > 0 {
					//ret[i-1].PreHash = block.Hash
				}
				ret = append(ret, blocks)

			} else {
				if len(ret) > int(sum) {
					ret = ret[:len(ret)-1]
				}

				return &ret, bk.ID, err
			}
		}
		if len(ret) > int(sum) {
			ret = ret[:len(ret)-1]
		}
	} else {
		//ret = ret[:len(ret)-1]
		return &ret, bk.ID, err
	}

	return &ret, bk.ID, err
}
func GetBlockListFromRedis() (*[]BlocksResult, int64, error) {
	var (
		ret []BlocksResult
		dat []Block
		i   int64
		err error
	)
	bks := &Block{}
	ret1, err, maxid := bks.GetBlocksListFromRedis()
	if err == nil {
		dat = *ret1
		for i = 0; i < int64(len(dat)); i++ {
			bkx, err := GetBlocksDetailedInfoHex(&dat[i])
			if err == nil {
				var blocks = BlocksResult{
					BlockID:      bkx.Header.BlockID,
					Time:         bkx.Time,
					EcosystemID:  bkx.EcosystemID,
					KeyID:        bkx.KeyID,
					NodePosition: bkx.NodePosition,
					Hash:         bkx.Hash,
					Tx:           bkx.Tx,
				}
				ret = append(ret, blocks)
			} else {
				return &ret, maxid, err
			}
		}
	} else {
		return &ret, maxid, err
	}

	return &ret, maxid, err
}

func GetBlockTpslistsFromRedis() (*[]ScanOutBlockTransactionRet, error) {

	ret1, err := GetTraninfoFromRedis(30)
	if err == nil {
		return ret1, err
	} else {
		return nil, err
	}
}

func GetBlocksDetailedInfoHex(bk *Block) (*BlockDetailedInfoHex, error) {
	var (
		transize int64
	)
	result := BlockDetailedInfoHex{}

	blck, err := block.UnmarshallBlock(bytes.NewBuffer(bk.Data), false)
	if err != nil {
		return &result, err
	}

	txDetailedInfoCollection := make([]TxDetailedInfoHex, 0, len(blck.Transactions))
	for _, tx := range blck.Transactions {
		txDetailedInfo := TxDetailedInfoHex{
			Hash: hex.EncodeToString(tx.TxHash),
		}

		if tx.TxContract != nil {
			txDetailedInfo.ContractName, txDetailedInfo.Params = GetMineParam(tx.TxHeader.EcosystemID, tx.TxContract.Name, tx.TxData, tx.TxHash)
			//txDetailedInfo.ContractName = tx.TxContract.Name
			//txDetailedInfo.Params = tx.TxData
			txDetailedInfo.KeyID = converter.AddressToString(tx.TxKeyID)
			txDetailedInfo.Time = tx.TxTime
			txDetailedInfo.Type = tx.TxType
			txDetailedInfo.Size = int64(len(tx.TxFullData))
			transize += txDetailedInfo.Size
		}

		if tx.TxHeader != nil {
			es := Ecosystem{}
			f, err := es.Get(tx.TxHeader.EcosystemID)
			if f && err == nil {
				txDetailedInfo.Ecosystem = tx.TxHeader.EcosystemID
				txDetailedInfo.Ecosystemname = es.Name
				if txDetailedInfo.Ecosystem == 0 {
					txDetailedInfo.Ecosystem = 1
				}
				if txDetailedInfo.Ecosystem == 1 {
					txDetailedInfo.Token_title = consts.SysEcosytemTitle
					if txDetailedInfo.Ecosystemname == "" {
						txDetailedInfo.Ecosystemname = "platform ecosystem"
					}
				} else {
					txDetailedInfo.Token_title = es.TokenTitle
				}

			}
		}

		txDetailedInfoCollection = append(txDetailedInfoCollection, txDetailedInfo)

		//log.WithFields(log.Fields{"block_id": blockModel.ID, "tx hash": txDetailedInfo.Hash, "contract_name": txDetailedInfo.ContractName, "key_id": txDetailedInfo.KeyID, "time": txDetailedInfo.Time, "type": txDetailedInfo.Type, "params": txDetailedInfoCollection}).Debug("Block Transactions Information")
	}
	if blck.Header.EcosystemID == 0 {
		blck.Header.EcosystemID = 1
	}
	if bk.EcosystemID == 0 {
		bk.EcosystemID = 1
	}
	header := BlockHeaderInfoHex{
		BlockID:      blck.Header.BlockID,
		Time:         blck.Header.Time,
		EcosystemID:  blck.Header.EcosystemID,
		KeyID:        converter.AddressToString(blck.Header.KeyID),
		NodePosition: blck.Header.NodePosition,
		Sign:         hex.EncodeToString(blck.Header.Sign),
		Hash:         hex.EncodeToString(blck.Header.Hash),
		Version:      blck.Header.Version,
	}

	bdi := BlockDetailedInfoHex{
		Header:        header,
		Hash:          hex.EncodeToString(bk.Hash),
		EcosystemID:   bk.EcosystemID,
		NodePosition:  bk.NodePosition,
		KeyID:         converter.AddressToString(bk.KeyID),
		Time:          bk.Time,
		Tx:            bk.Tx,
		RollbacksHash: hex.EncodeToString(bk.RollbacksHash),
		MrklRoot:      hex.EncodeToString(blck.MrklRoot),
		BinData:       hex.EncodeToString(blck.BinData),
		SysUpdate:     blck.SysUpdate,
		GenBlock:      blck.GenBlock,
		//StopCount:     blck.s,
		BlockSize:     int64(len(bk.Data)),
		TranTotalSize: transize,
		Transactions:  txDetailedInfoCollection,
	}

	return &bdi, nil
}
func GetMineParam(ecosystem int64, name string, param map[string]interface{}, TxHash []byte) (string, string) {
	escape := func(value interface{}) string {
		return strings.Replace(fmt.Sprint(value), `'`, `''`, -1)
	}
	if name == "@1CallDelayedContract" && ecosystem == 1 {
		v, ok := param["Id"]
		if ok {
			idstr := escape(v)
			if idstr == "4" {
				return "@1Mint", GetMineIncomeParam(TxHash)
			}
		}

	}
	dataBytes, _ := json.Marshal(param)
	return name, string(dataBytes)
}

func GetMineIncomeParam(hash []byte) string {
	ret := make(map[string]interface{})
	ts := &MineIncomehistory{}
	f, err := ts.Get(hash)
	if err == nil && f {
		ret["miner"] = ts.Devid
		ret["minerowner"] = ts.Keyid
		ret["profiter"] = ts.Mineid
		ret["type"] = ts.Type
		ret["staked"] = ts.Nonce
		ret["earnings"] = ts.Amount
		//return ret
	}
	dataBytes, _ := json.Marshal(ret)
	return string(dataBytes)
}
func SyncBlockinfoToRedis() {
	var trans []Block
	if err := GetDB(nil).Limit(10).Order("id desc").Find(&trans).Error; err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("SyncBlockinfoToRedis find err")
		return
	}
	go SetTransactionBlockLastToRedis(GetMaxBlockFromRedis(trans))
	blockInfo, err := json.Marshal(trans)
	if err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("SyncBlockinfoToRedis json err")
		return
	}
	//blocks, err := msgpack.Marshal(blockInfo)
	//if err != nil {
	//	log.WithFields(log.Fields{"warn": err}).Warn("SyncBlockinfoToRedis msgpack err")
	//	return
	//}
	value, err := GzipEncode(blockInfo)
	if err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("SyncBlockinfoToRedis GzipEncode err")
		return
	}
	rd := RedisParams{
		Key:   "blockChain-list",
		Value: string(value),
	}
	if err := rd.Set(); err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("SyncBlockinfoToRedis setdb err")
	}
	if err := sendBlockListToWebsocket(&trans); err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("sendBlockListToWebsocket err")
	}
}
func GzipEncode(in []byte) ([]byte, error) {
	var (
		buffer bytes.Buffer
		out    []byte
		err    error
	)
	writer := gzip.NewWriter(&buffer)
	_, err = writer.Write(in)
	if err != nil {
		err = writer.Close()
		return out, err
	}
	if err = writer.Close(); err != nil {
		return out, err
	}
	return buffer.Bytes(), nil
}

func GzipDecode(in []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(in))
	if err != nil {
		var out []byte
		return out, err
	}
	defer func() {
		if err = reader.Close(); err != nil {
			println("Gzip unzip error", err.Error())
		}
	}()
	return ioutil.ReadAll(reader)
}

func (b *Block) GetBlocksListFromRedis() (*[]Block, error, int64) {
	var err error
	//blockchain := new([]Block)
	var blockchain []Block
	rd := RedisParams{
		Key:   "blockChain-list",
		Value: "",
	}
	if err := rd.Get(); err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("GetBlocksListFromRedis getdb err")
		return nil, err, 0
	}
	value, err1 := GzipDecode([]byte(rd.Value))
	if err1 != nil {
		log.WithFields(log.Fields{"warn": err1}).Warn("GetBlocksListFromRedis GzipDecode err")
		return nil, err1, 0
	}
	//var blockInfo []byte
	//err = msgpack.Unmarshal(value, &blockInfo)
	//if err != nil {
	//	log.WithFields(log.Fields{"warn": err}).Warn("GetBlocksListFromRedis msgpack err")
	//	return nil, err, 0
	//}

	if err = json.Unmarshal(value, &blockchain); err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("GetBlocksListFromRedis json err")
		return nil, err, 0
	}
	maxId := GetMaxBlockFromRedis(blockchain)
	return &blockchain, nil, maxId
}
func sendBlockListToWebsocket(ret1 *[]Block) error {
	var (
		ret []BlocksResult
		dat []Block
		i   int64
		err error
	)
	dat = *ret1
	for i = 0; i < int64(len(dat)); i++ {
		bkx, err := GetBlocksDetailedInfoHex(&dat[i])
		if err == nil {
			var blocks = BlocksResult{
				BlockID:      bkx.Header.BlockID,
				Time:         bkx.Time,
				EcosystemID:  bkx.EcosystemID,
				KeyID:        bkx.KeyID,
				NodePosition: bkx.NodePosition,
				Hash:         bkx.Hash,
				Tx:           bkx.Tx,
			}
			ret = append(ret, blocks)
		} else {
			return err
		}
	}
	if err = SendBlockList(&ret); err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("Send Block Transaction err")
		return err
	}
	return nil
}
func SendBlockList(topblocks *[]BlocksResult) error {
	dat := ResponseTopTitle{}
	dat.Cmd = CmdTopBlocks
	dat.Data = topblocks

	ds, err := json.Marshal(dat)
	if err != nil {
		return err
	}
	err = WriteChannelByte(ChannelTopData, ds)
	}
	return maxid
}
