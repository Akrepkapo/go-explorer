package models

import (
	"encoding/hex"
	"fmt"
	"github.com/IBAX-io/go-ibax/packages/converter"
	"strconv"
	"time"

	"github.com/IBAX-io/go-explorer/consts"
	"github.com/vmihailenco/msgpack/v5"
)

type ScanOut struct {
	Blockid           int64  `gorm:"not null" ` //
	BlockSizes        int64  `gorm:"not null" ` //
	BlockTranscations int64  `gorm:"not null" `
	Hash              string `gorm:"not null"`
	RollbacksHash     string `gorm:"not null"`
	//Data           string `gorm:"not null"`
	//Tx             int64  `gorm:"not null"`
	EcosystemID    int64  `gorm:"not null default 0"`
	KeyID          string `gorm:"not null default 0"`
	NodePosition   string `gorm:"not null default 0"`
	Time           int64  `gorm:"not null"`
	CurrentVersion string `gorm:"not null"`

	TotalCounts          int64 `gorm:"not null" ` //total count
	TotalCapacitys       int64 `gorm:"not null" ` //total capacity
	BlockTranscationSize int64 `gorm:"not null" `
	QueueTranscations    int64 `gorm:"not null" `
	GuardianNodes        int64 `gorm:"not null" `
	StorageCapacitys     int64 `gorm:"not null" `
	CastNodes            int64 `gorm:"not null" `
	Ecosystems           int64 `gorm:"not null" `
	SubNodes             int64 `gorm:"not null" `
	CLBNodes             int64 `gorm:"not null" `

	Circulations   string `gorm:"not null" `
	TotalAmounts   string `gorm:"not null" `
	MintAmounts    string `gorm:"not null" `
	StakeAmounts   string `gorm:"not null" `
	BlockDetail    BlockDetailedInfoHex
	BlockTxDetails []BlockTxDetailedInfoHex `gorm:"not null" `
}

type ScanOutRet struct {
	Blockid           int64  `gorm:"not null" json:"block_id"`    //
	BlockSizes        int64  `gorm:"not null" json:"block_sizes"` //
	BlockTranscations int64  `gorm:"not null" json:"block_transcations"`
	Hash              string `gorm:"not null" json:"hash"`
	RollbacksHash     string `gorm:"not null" json:"rollbacks_hash"`
	//Data           string `gorm:"not null"`
	//BlockSizeTotals      int64 `gorm:"not null" json:"block_size_totals"` //
	BlockTranscationSize int64 `gorm:"not null" json:"block_transcation_size" `
	QueueTranscations    int64 `gorm:"not null" json:"queue_transcations"`
	GuardianNodes        int64 `gorm:"not null" json:"guardian_nodes"`
	StorageCapacitys     int64 `gorm:"not null" json:"storage_capacitys"`
	CastNodes            int64 `gorm:"not null" json:"cast_nodes"`
	Ecosystems           int64 `gorm:"not null" json:"ecosystems"`
	SubNodes             int64 `gorm:"not null" json:"sub_nodes"`
	CLBNodes             int64 `gorm:"not null" json:"clb_nodes"`

	Circulations string `gorm:"not null" json:"circulations_amount"`
	TotalAmounts string `gorm:"not null" json:"total_amount"`
	MintAmounts  string `gorm:"not null" json:"mint_amount"`
	StakeAmounts string `gorm:"not null" json:"stake_amount"`
	NodeBlocks   string `gorm:"not null" json:"node_blocks"`
}

type ScanOutBlockTransactionRet struct {
	BlockId           int64 `json:"block_id"`            //
	BlockSizes        int64 `json:"block_size" `         //
	BlockTranscations int64 `json:"block_transcations" ` //
}

var ScanPrefix = "scan-"
var ScanOutStPrefix = "scan-out-"
var ScanOutLastest = "lastest"
var ScanOutBlock = "blockTransaction"
var GetScanOut chan bool
var SendScanOut chan bool

func (s *ScanOutBlockTransactionRet) Marshal(q []ScanOutBlockTransactionRet) (string, error) {
	if res, err := msgpack.Marshal(q); err != nil {
		return "", err
	} else {
		return string(res), err
	}
}

func (s *ScanOutBlockTransactionRet) Unmarshal(bt string) (q []ScanOutBlockTransactionRet, err error) {
	if err := msgpack.Unmarshal([]byte(bt), &q); err != nil {
		return nil, err
	}
	return q, nil
}

func (s *ScanOut) Marshal() ([]byte, error) {
	if res, err := msgpack.Marshal(s); err != nil {
		return nil, err
	} else {
		return res, err
	}
}

func (s *ScanOut) Unmarshal(bt []byte) error {
	if err := msgpack.Unmarshal(bt, &s); err != nil {
		return err
	}
	return nil
}

func (s *ScanOut) Get(id int64) (bool, error) {
	rp := &RedisParams{
		Key: ScanPrefix + strconv.FormatInt(id, 10),
	}
	for i := 0; i < 10; i++ {
		err := rp.Get()
		if err == nil {
			err = s.Unmarshal([]byte(rp.Value))
			return true, err
		}
		if err.Error() == "redis: nil" || err.Error() == "EOF" {
			break
		} else {
			time.Sleep(200 * time.Millisecond)
		}

	}

	return false, nil
}

func (m *ScanOut) Del(id int64) error {
	rp := &RedisParams{
		Key: ScanPrefix + strconv.FormatInt(id, 10),
	}

	for i := 0; i < 5; i++ {
		err := rp.Del()
		if err == nil {
			break
		}
	}

	return nil
}

func (m *ScanOut) DelRange(id, count int64) error {

	for i := int64(0); i < count; i++ {
		rp := &RedisParams{
			Key: ScanPrefix + strconv.FormatInt(id+i, 10),
		}

		err := rp.Del()
		if err != nil {
			return err
		}
	}

	return nil
}
func (m *ScanOut) Del_Redis(id int64) error {
	rd := RedisParams{
		Key:   ScanOutStPrefix + strconv.FormatInt(id, 10),
		Value: "",
	}
	if err := rd.Del(); err != nil {
		//log.WithFields(log.Fields{"err": err}).Warn("Del_Redis failed")
		return err
	}
	return nil
}

func (m *ScanOut) Insert_Redis() error {
	val, err := m.Marshal()
	if err != nil {
		return err
	}

	var so ScanOut
	f, err := so.Get_RedisId(m.Blockid)
	if err != nil {
		if err.Error() == "redis: nil" || err.Error() == "EOF" {
			var nb NodeBlocks
			nb.AddOneFromRedis(m.NodePosition)
			var bs BlockSizeTotal
			bs.ReSetRedis(m.Blockid, m.BlockSizes)
		} else {
			return err
		}

	}
	if f {
		if so.NodePosition != m.NodePosition {
			var nb NodeBlocks
			nb.SubOneRedis(so.NodePosition)
			nb.AddOneFromRedis(m.NodePosition)
			var bs BlockSizeTotal
			bs.ReSetRedis(m.Blockid, m.BlockSizes-so.BlockSizes)
		}
	}

	rd := RedisParams{
		Key:   ScanOutStPrefix + ScanOutLastest,
		Value: string(val),
	}
	err = rd.Set()
	if err != nil {
		return err
	}

	//rd = RedisParams{
	//	Key:   ScanOutStPrefix + strconv.FormatInt(m.Blockid, 10),
	//	Value: string(val),
	//}
	//err = rd.Set()

	return err
}
func (m *ScanOut) Get_Db(id int64) (bool, error) {
	var bk Block
	fb, err := bk.Get(id)
	if err != nil {
		return false, err
	}
	if fb {
		m.Blockid = bk.ID
		m.BlockTranscations = int64(bk.Tx)
		m.BlockSizes = int64(len(bk.Data))
	} else {
		return false, nil
	}
	return true, nil
}
func (m *ScanOut) Get_Redis(id int64) (bool, error) {
	rd := RedisParams{
		Key:   ScanOutStPrefix + strconv.FormatInt(id, 10),
		Value: "",
	}
	err := rd.Get()
	if err != nil {
		if err.Error() == "redis: nil" || err.Error() == "EOF" {
			return false, nil
		}
		return false, err
	}
	err = m.Unmarshal([]byte(rd.Value))
	if err != nil {
		return false, err
	}
	return true, err
}

func (m *ScanOut) Get_RedisId(id int64) (bool, error) {
	rd := RedisParams{
		Key:   ScanOutStPrefix + strconv.FormatInt(id, 10),
		Value: "",
	}
	err := rd.Get()
	if err != nil {
		return false, err
	}
	err = m.Unmarshal([]byte(rd.Value))
	if err != nil {
		return false, err
	}
	return true, err
}

func GetScanOutDataToRedis() error {
	err := processScanOutBlocks()
	return err
}

func (ret *ScanOut) Changes() error {

	ret.TotalAmounts = consts.TotalSupplyToken

	var mh MineIncomehistory
	f, err := mh.GetID(ret.Blockid)
	if err != nil {
		return err
	}
	if f {
		ret.TotalCounts = mh.Nonce
	}

	var key Key
	tm, err := key.GetTotalAmount()
	if err != nil {
		return err
	}
	stakeamount, err := key.GetStakeAmount()
	if err != nil {
		return err
	}
	ret.Circulations = tm.String()
	ret.StakeAmounts = stakeamount

	ess, _, err := GetAllSystemStatesIDs()
	if err != nil {
		return err
	}
	ret.Ecosystems = int64(len(ess))

	var mst MinePledgeStatus
	gnode, casts, err := mst.GetCastNodeandGuardianNode()
	if err != nil {
		return err
	}
	ret.GuardianNodes = gnode
	ret.CastNodes = casts

	var mi MineInfo
	capacity, err := mi.GetGuardianNodeCapacity()
	if err != nil {
		return err
	}
	ret.StorageCapacitys = capacity

	var sp StateParameter
	sp.ecosystem = 1
	mb, err := sp.GetMintAmount()
	if err != nil {
		return err
	}
	ret.MintAmounts = mb
	return nil
}

func processScanOutBlocks() error {
	var bk Block
	var ret ScanOut
	cbk := InfoBlock{}
	if _, err := cbk.Get(); err != nil {
		return err
	}
	if cbk.BlockID == 2 {
		err := processScanOutFirstBlocks()
		if err != nil {
			return err
		}
	}
	ret.Blockid = cbk.BlockID
	ret.KeyID = converter.AddressToString(cbk.KeyID)
	ret.EcosystemID = cbk.EcosystemID
	ret.Time = cbk.Time
	ret.CurrentVersion = cbk.CurrentVersion
	ret.Hash = hex.EncodeToString(cbk.Hash)
	ret.NodePosition = cbk.NodePosition
	if ret.EcosystemID == 0 {
		ret.EcosystemID = 1
	}
	f, err := bk.Get(cbk.BlockID)
	if err != nil {
		return err
	}
	if f {
		ret.BlockTranscations = int64(bk.Tx)
		ret.RollbacksHash = hex.EncodeToString(bk.RollbacksHash)
		ret.BlockSizes = int64(len(bk.Data))
		ts, err := GetBlocksDetailedSizeHexByScanOut(&bk)
		if err != nil {
			return err
		}
		ret.BlockTranscationSize = ts
	}

	err = ret.Insert_redisdb()
	if err != nil {
		return fmt.Errorf("Insert_redisdb scanout:%s\n", err.Error())
	}

	return nil
}

func processScanOutFirstBlocks() error {
	var so ScanOut
	var bk Block

	fb, err := bk.Get(1)
	if err != nil {
		return err
	}
	if fb {
		so.Blockid = 1
		so.KeyID = converter.AddressToString(bk.KeyID)
		so.EcosystemID = bk.EcosystemID
		so.Time = bk.Time
		//so.CurrentVersion =
		so.Hash = hex.EncodeToString(bk.Hash)
		so.NodePosition = converter.Int64ToStr(bk.NodePosition)
		if so.EcosystemID == 0 {
			so.EcosystemID = 1
		}

		so.BlockTranscations = int64(bk.Tx)
		so.RollbacksHash = hex.EncodeToString(bk.RollbacksHash)
		so.BlockSizes = int64(len(bk.Data))
		ts, bdt, _, err := Deal_TransactionBlockTxDetial(&bk)
		if err != nil {
			return err
		}
		so.CurrentVersion = converter.IntToStr(bdt.Header.Version)
		so.BlockTranscationSize = ts
	}

	err = so.Insert_redisdb()
	if err != nil {
		return err
	}
	return nil
}

func (s *ScanOut) Insert_redisdb() error {
	errCs := s.Changes()
	if errCs != nil {
		return fmt.Errorf("changes err:%s\n", errCs.Error())
	}
	val, err := s.Marshal()
	if err != nil {
		return err
	}
	rp := RedisParams{
		Key:   ScanPrefix + strconv.FormatInt(s.Blockid, 10),
		Value: string(val),
	}

	for i := 0; i < 10; i++ {
		err = rp.Set()
		if err == nil {
			break
		} else {
			time.Sleep(10 * time.Millisecond)
		}

	}

	return err
}

func (m *ScanOut) GetRedisdashboard() (*ScanOutRet, error) {
	var rets ScanOutRet
	rd := RedisParams{
		Key:   ScanOutStPrefix + ScanOutLastest,
		Value: "",
	}
	err := rd.Get()
	if err != nil {
		if err.Error() == "redis: nil" || err.Error() == "EOF" {
			return &rets, nil
		}
		return &rets, err
	}

	err = m.Unmarshal([]byte(rd.Value))
	if err != nil {
		return &rets, err
	}

	var bc BlockID
	fc, _ := bc.GetbyName(consts.TransactionsMax)
	if fc {
		rets.QueueTranscations = bc.ID
	}

	np, err := strconv.Atoi(m.NodePosition)
	if err != nil {
		return &rets, err
	}
	nps := strconv.Itoa(np + 1)
	rets.NodePosition = nps

	var nb NodeBlocks
	fb, _ := nb.Get_redis(m.NodePosition)
	if fb {
		rets.NodeBlocks = strconv.FormatInt(nb.Count, 10)
	} else {
		rets.NodeBlocks = "0"
	}

	var bs BlockSizeTotal
	fs, _ := bs.Get_Redis()
	if fs {
		rets.TotalCapacitys = m.TocapacityString(bs.Count)
	}

	rets.Blockid = m.Blockid

	rets.Hash = m.Hash
	rets.RollbacksHash = m.RollbacksHash
	//rets.Tx    = m.Blockid
	rets.EcosystemID = m.EcosystemID
	rets.KeyID = m.KeyID
	//rets.NodePosition = m.NodePosition
	rets.Time = m.Time
	rets.CurrentVersion = m.CurrentVersion

	rets.TotalCounts = m.TotalCounts
	//rets.TotalCapacitys = m.b
	//rets.BlockSizes = m.BlockSizes
	rets.BlockTranscations = m.BlockTranscations
	rets.BlockTranscationSize = m.BlockTranscationSize
	//rets.QueueTranscations = m.QueueTranscations
	rets.GuardianNodes = m.GuardianNodes
	rets.StorageCapacitys = m.StorageCapacitys
	rets.CastNodes = m.CastNodes
	rets.Ecosystems = m.Ecosystems
	rets.SubNodes = m.SubNodes
	rets.CLBNodes = m.CLBNodes

	rets.Circulations = m.Circulations
	rets.TotalAmounts = consts.TotalSupplyToken
	rets.MintAmounts = m.MintAmounts
	rets.StakeAmounts = m.StakeAmounts

	return &rets, err
}

func (m *ScanOut) GetRedisLastest() (bool, error) {
	rd := RedisParams{
		Key:   ScanOutStPrefix + ScanOutLastest,
		Value: "",
	}
	err := rd.Get()
	if err != nil {
		if err.Error() == "redis: nil" || err.Error() == "EOF" {
			return false, nil
		}
		return false, err
	}
	err = m.Unmarshal([]byte(rd.Value))
	if err != nil {
		return false, err
	}

	return true, err
}

func (m *ScanOut) GetBlockTransactions(count int) (*[]ScanOutBlockTransactionRet, error) {
	var (
		ret []ScanOutBlockTransactionRet
	)

	f, err := m.GetRedisLastest()
	if err != nil {
		return &ret, err
	}
	if f {
		var st ScanOutBlockTransactionRet
		st.BlockId = m.Blockid
		st.BlockTranscations = m.BlockTranscations
		st.BlockSizes = m.BlockSizes
		ret = append(ret, st)

		for i := 1; i <= count; i++ {
			var so ScanOut
			bid := m.Blockid - int64(i)
			if bid > 0 {
				fs, err := so.Get_Db(bid)
				if err != nil {
					return &ret, err
				}
				if fs {
					var sti ScanOutBlockTransactionRet
					sti.BlockId = so.Blockid
					sti.BlockTranscations = so.BlockTranscations
					sti.BlockSizes = so.BlockSizes
					ret = append(ret, sti)
				}
			} else {
				break
			}

		}

	}

	return &ret, err
}

func (m *ScanOut) GetBlockDetialRespones(page int, limit int) *BlockDetailedInfoHexRespone {
	var (
		ret     BlockDetailedInfoHexRespone
		st, end int
	)
	ret.Limit = limit
	ret.Page = page
	ret.Total = int64(len(m.BlockDetail.Transactions))

	ret.Header = m.BlockDetail.Header

	ret.NodePosition = m.BlockDetail.NodePosition
	ret.Hash = m.BlockDetail.Hash
	ret.RollbacksHash = m.BlockDetail.RollbacksHash
	ret.KeyID = m.BlockDetail.KeyID
	ret.GenBlock = m.BlockDetail.GenBlock
	ret.Time = m.BlockDetail.Time
	ret.EcosystemID = m.BlockDetail.EcosystemID
	ret.BlockSize = m.BlockDetail.BlockSize
	ret.TranTotalSize = m.BlockDetail.TranTotalSize
	ret.MrklRoot = m.BlockDetail.MrklRoot
	ret.Tx = m.BlockDetail.Tx
	ret.BinData = m.BlockDetail.BinData
	ret.StopCount = m.BlockDetail.StopCount
	ret.SysUpdate = m.BlockDetail.SysUpdate

	if page > 0 {
		st = (page - 1) * limit
		end = page * limit
	} else {
		st = 0
		end = limit
	}
	if end > int(ret.Total) {
		end = int(ret.Total)
	}
	if st < int(ret.Total) {
		ret.Transactions = m.BlockDetail.Transactions[st:end]
	}
	return &ret
}

func (m *ScanOut) TocapacityString(count int64) string {
	rs := float64(count) / float64(1048576)
	if rs >= 1024 {
		rs = rs / float64(1024)
		return strconv.FormatFloat(rs, 'f', 2, 64) + "G"
	}
	return strconv.FormatFloat(rs, 'f', 2, 64) + "M"
}

func DealRedisDashboardHistoryMap() error {
	return GetDBDealTraninfo(30)
}

func SendStatisticsSignal() {
	for len(GetScanOut) > 0 {
		<-GetScanOut
	}
	select {
	case GetScanOut <- true:
	default:
	}
	for len(SendScanOut) > 0 {
		<-SendScanOut
	}
	select {
	case SendScanOut <- true:
	default:
	}
}
