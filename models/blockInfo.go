package models

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/IBAX-io/go-explorer/conf"
	"github.com/IBAX-io/go-ibax/packages/consts"
	"github.com/IBAX-io/go-ibax/packages/converter"
	"github.com/IBAX-io/go-ibax/packages/crypto"
	"github.com/IBAX-io/go-ibax/packages/transaction"
	"github.com/IBAX-io/go-ibax/packages/types"
	"github.com/IBAX-io/go-ibax/packages/utils"
	log "github.com/sirupsen/logrus"
)

const SysEcosytemTitle = "IBXC"

// blocks is storing block data
type blocks struct {
	Header            utils.BlockData
	PrevHeader        *utils.BlockData
	PrevRollbacksHash []byte
	MrklRoot          []byte
	BinData           []byte
	Transactions      []*transaction.Transaction
	SysUpdate         bool
	GenBlock          bool // it equals true when we are generating a new block
	Notifications     []types.Notifications
}

// InfoBlock is model
type InfoBlock struct {
	Hash           []byte `gorm:"not null"`
	EcosystemID    int64  `gorm:"not null default 0"`
	KeyID          int64  `gorm:"not null default 0"`
	NodePosition   string `gorm:"not null default 0"`
	BlockID        int64  `gorm:"not null"`
	Time           int64  `gorm:"not null"`
	CurrentVersion string `gorm:"not null"`
	Sent           int8   `gorm:"not null"`
	RollbacksHash  []byte `gorm:"not null"`
}

// TableName returns name of table
}

func Deal_TransactionBlockTxDetial(mc *Block) (int64, *BlockDetailedInfoHex, *[]BlockTxDetailedInfoHex, error) {
	var (
		ret []BlockTxDetailedInfoHex
		ts  int64
	)

	//bk := &Block{}
	rt, err := GetBlocksDetailedInfoHexByScanOut(mc)
	if err == nil {
		for j := 0; j < int(rt.Tx); j++ {
			bh := BlockTxDetailedInfoHex{}
			bh.BlockID = rt.Header.BlockID
			bh.ContractName = rt.Transactions[j].ContractName
			bh.Hash = rt.Transactions[j].Hash
			bh.KeyID = rt.Transactions[j].KeyID
			//bh.Params = rt.Transactions[j].Params
			bh.Time = rt.Transactions[j].Time
			bh.Type = rt.Transactions[j].Type
			if bh.Time == 0 {
				bh.Time = rt.Time
			}
			if bh.KeyID == "" {
				bh.KeyID = rt.KeyID
			}
			if bh.Ecosystem == 0 {
				bh.Ecosystem = 1
			}
			bh.Ecosystemname = rt.Transactions[j].Ecosystemname
			bh.Ecosystem = rt.Transactions[j].Ecosystem
			if bh.Ecosystem == 1 || bh.Ecosystem == 0 {
				bh.Token_title = SysEcosytemTitle
				if bh.Ecosystemname == "" {
					bh.Ecosystemname = "platform ecosystem"
				}
			} else {
				bh.Token_title = rt.Transactions[j].Token_title
			}
			//dlen := unsafe.Sizeof(rt.Transactions[j])
			bh.Size = rt.Transactions[j].Size
			ts += bh.Size
			lg1, err := json.Marshal(rt.Transactions[j].Params)
			if err == nil {
				bh.Params = string(lg1)
			}

			ret = append(ret, bh)
		}
	}

	return ts, rt, &ret, err
}

func GetBlocksDetailedSizeHexByScanOut(mc *Block) (int64, error) {
	var (
		ts int64
	)

	//bk := &Block{}
	rt, err := GetBlocksDetailedInfoHexByScanOut(mc)
	if err == nil {
		for j := 0; j < int(rt.Tx); j++ {
			bh := BlockTxDetailedInfoHex{}
			bh.Size = rt.Transactions[j].Size
			ts += bh.Size
		}
	}

	return ts, err
}

func GetBlocksDetailedInfoHexByScanOut(mc *Block) (*BlockDetailedInfoHex, error) {
	var (
		transize int64
	)
	result := BlockDetailedInfoHex{}
	blck, err := UnmarshallBlock(bytes.NewBuffer(mc.Data), false)
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
			txDetailedInfo.KeyID = converter.AddressToString(tx.TxKeyID)
			txDetailedInfo.Time = tx.TxTime
			txDetailedInfo.Type = tx.TxType
			txDetailedInfo.Size = int64(len(tx.TxFullData))
			transize += txDetailedInfo.Size
		}

		if txDetailedInfo.Time == 0 {
			txDetailedInfo.Time = mc.Time
		}
		if txDetailedInfo.KeyID == "" {
			txDetailedInfo.KeyID = converter.AddressToString(mc.KeyID)
		}

		if tx.TxHeader != nil {
			es := Ecosystem{}
			f, err := es.Get(tx.TxHeader.EcosystemID)
			if f && err == nil {
				txDetailedInfo.Ecosystem = tx.TxHeader.EcosystemID
				if txDetailedInfo.Ecosystem == 0 {
					txDetailedInfo.Ecosystem = 1
				}
				txDetailedInfo.Ecosystemname = es.Name
				if txDetailedInfo.Ecosystem == 1 || txDetailedInfo.Ecosystem == 0 {
					txDetailedInfo.Token_title = SysEcosytemTitle
					if txDetailedInfo.Ecosystemname == "" {
						txDetailedInfo.Ecosystemname = "platform ecosystem"
					}
				} else {
					txDetailedInfo.Token_title = es.TokenTitle
				}

			}
		} else {
			if txDetailedInfo.Ecosystem == 0 {
				txDetailedInfo.Ecosystem = 1
			}
			if txDetailedInfo.Ecosystem == 1 || txDetailedInfo.Ecosystem == 0 {
				txDetailedInfo.Token_title = SysEcosytemTitle
				if txDetailedInfo.Ecosystemname == "" {
					txDetailedInfo.Ecosystemname = "platform ecosystem"
				}
			}
		}

		txDetailedInfoCollection = append(txDetailedInfoCollection, txDetailedInfo)

		//log.WithFields(log.Fields{"block_id": blockModel.ID, "tx hash": txDetailedInfo.Hash, "contract_name": txDetailedInfo.ContractName, "key_id": txDetailedInfo.KeyID, "time": txDetailedInfo.Time, "type": txDetailedInfo.Type, "params": txDetailedInfoCollection}).Debug("Block Transactions Information")
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
	//prehash
	if mc.ID > 1 {
		var bk Block
		pfb, err := bk.Get(mc.ID - 1)
		if err != nil {
			return &result, err
		}
		if pfb {
			header.PreHash = hex.EncodeToString(bk.Hash)
		}

	}
	if header.EcosystemID == 0 {
		header.EcosystemID = 1
	}

	bdi := BlockDetailedInfoHex{
		Header:        header,
		Hash:          hex.EncodeToString(mc.Hash),
		EcosystemID:   mc.EcosystemID,
		NodePosition:  mc.NodePosition,
		KeyID:         converter.AddressToString(mc.KeyID),
		Time:          mc.Time,
		Tx:            mc.Tx,
		RollbacksHash: hex.EncodeToString(mc.RollbacksHash),
		MrklRoot:      hex.EncodeToString(blck.MrklRoot),
		BinData:       hex.EncodeToString(blck.BinData),
		SysUpdate:     blck.SysUpdate,
		GenBlock:      blck.GenBlock,
		//StopCount:     blck.s,
		BlockSize:     int64(len(mc.Data)),
		TranTotalSize: transize,
		Transactions:  txDetailedInfoCollection,
	}

	if bdi.EcosystemID == 0 {
		bdi.EcosystemID = 1
	}
	return &bdi, nil
}

func UnmarshallBlock(blockBuffer *bytes.Buffer, fillData bool) (*blocks, error) {
	header, prev, err := utils.ParseBlockHeader(blockBuffer)
	if err != nil {
		return nil, err
	}

	logger := log.WithFields(log.Fields{"block_id": header.BlockID, "block_time": header.Time, "block_wallet_id": header.KeyID,
		"block_state_id": header.EcosystemID, "block_hash": header.Hash, "block_version": header.Version})
	transactions := make([]*transaction.Transaction, 0)

	var mrklSlice [][]byte

	// parse transactions
	for blockBuffer.Len() > 0 {
		transactionSize, err := converter.DecodeLengthBuf(blockBuffer)
		if err != nil {
			logger.WithFields(log.Fields{"type": consts.UnmarshallingError, "error": err}).Error("transaction size is 0")
			return nil, fmt.Errorf("bad block format (%s)", err)
		}
		if blockBuffer.Len() < int(transactionSize) {
			logger.WithFields(log.Fields{"size": blockBuffer.Len(), "match_size": int(transactionSize), "type": consts.SizeDoesNotMatch}).Error("transaction size does not matches encoded length")
			return nil, fmt.Errorf("bad block format (transaction len is too big: %d)", transactionSize)
		}

		if transactionSize == 0 {
			logger.WithFields(log.Fields{"type": consts.EmptyObject}).Error("transaction size is 0")
			return nil, fmt.Errorf("transaction size is 0")
		}

		bufTransaction := bytes.NewBuffer(blockBuffer.Next(int(transactionSize)))
		t, err := transaction.UnmarshallTransaction(bufTransaction, fillData)
		if err != nil {
			if t != nil && t.TxHash != nil {
				transaction.MarkTransactionBad(t.DbTransaction, t.TxHash, err.Error())
			}
			return nil, fmt.Errorf("parse transaction error(%s)", err)
		}
		t.BlockData = &header

		transactions = append(transactions, t)

		// build merkle tree
		if len(t.TxFullData) > 0 {
			dSha256Hash := crypto.DoubleHash(t.TxFullData)
			//if err != nil {
			//	logger.WithFields(log.Fields{"type": consts.CryptoError, "error": err}).Error("double hashing tx full data")
			//	return nil, err
			//}
			dSha256Hash = converter.BinToHex(dSha256Hash)
			mrklSlice = append(mrklSlice, dSha256Hash)
		}
	}

	if len(mrklSlice) == 0 {
		mrklSlice = append(mrklSlice, []byte("0"))
	}
	mrkl, err := utils.MerkleTreeRoot(mrklSlice)
	if err != nil {
		return nil, err
	}
	return &blocks{
		Header:            header,
		PrevRollbacksHash: prev.RollbacksHash,
		Transactions:      transactions,
		MrklRoot:          mrkl,
	}, nil
}
