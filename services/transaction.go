/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package services

import (
	"github.com/IBAX-io/go-explorer/models"
	"github.com/IBAX-io/go-ibax/packages/converter"
	"github.com/shopspring/decimal"
)

func Get_Group_TransactionStatus(ids int, icount int, order string) (*[]models.TransactionStatusHex, int64, error) {
	ts := &models.TransactionStatus{}
	ret, num, err := ts.GetTransactions(ids, icount, order)
	return ret, num, err
	//	ret, num, err := bt.Get_BlockTransactionsLast_Sqlite(id, ids, icount, order)
	//	if err == nil && *ret != nil && num > 0 {
	//		//fmt.Println("Get_BlockTransactions_Sqlite  ok ids:%d icount:%d", ids, icount)
	//		return ret, num, err
	//	}
	//}
	ret, err := ts.GetTransactionBlockLastFromRedis()
	//ret, num, err := ts.Get_BlockTransactionsLast(id, int64(ids), int64(icount), order)
	//fmt.Println("Get_BlockTransactions pg  ok ids:%d icount:%d", ids, icount)
	return &ret, err
	//return nil, int(0), nil
}
func Get_Group_TransactionWallet(ids int, icount int, wallet string, searchType string) (*[]models.HistoryHex, int64, decimal.Decimal, error) {
	ts := &models.History{}

	ret, num, total, err := ts.GetWallets(ids, icount, wallet, searchType)
	return ret, num, total, err
}

func Get_Group_TransactionEcosytemWallet(id int64, ids int, icount int, wallet string, searchType string) (*[]models.HistoryHex, int64, decimal.Decimal, error) {
	ts := &models.History{}
	return ts.GetEcosytemWallets(id, ids, icount, wallet, searchType)
}

func Get_Group_WalletHistory(id int64, wallet string) (*models.WalletHistoryHex, error) {
	ts := &models.History{}
	key := &models.Key{}
	ret, err := ts.GetWalletTotals(wallet)
	if err != nil {
		return ret, err
	}

	dat, err := key.Get(id, wallet)
	if err != nil {
		return ret, err
	}
	ret.Amount = dat.Amount

	return ret, err
}

func Get_Group_Wallet_Total(ids int, icount int, order string, wallet string) (int64, int, *[]models.EcosyKeyTotalHex, error) {
	key := &models.Key{}
	return key.GetTotal(ids, icount, order, wallet)
}

func Get_transaction_HashHistory(logourl string, hash []byte) (*models.HistoryExplorer, error) {
	ts := &models.History{}
	//ret, err := ts.Get(hash)
	ret, err := ts.GetExplorer(logourl, hash)
	return ret, err
}

func Get_transaction_Hash(logourl string, hash string) (*map[string]interface{}, error) {
	ret := make(map[string]interface{})
	hashdat := []byte(converter.HexToBin(hash))

	ret2, err2 := Get_transaction_HashHistory(logourl, hashdat)
	//ret1, err1 := Get_transaction_Hashstatus(hashdat)
	if err2 != nil {
		return &ret, err2
	}

	var ret3 models.BlockTxDetailedInfoHex
	f, err3 := ret3.GetDb_txdetailedHash(hash)
	//ret3, err3 := GetBlockTransactionHash(hashdat)
	if err3 != nil {
		return &ret, err3
	}
	if f {
		ret["TxDetailedInfoHex"] = ret3
	}

	if ret2.Blockid == 0 {
		ret2.Senderid = ret3.KeyID
		ret2.Ecosystemname = ret3.Ecosystemname
		ret2.Ecosystem = ret3.Ecosystem
		ret2.Token_title = ret3.Token_title
		ret2.Txhash = ret3.Hash
		ret2.Createdat = ret3.Time
		ret2.CreateSetup = ret3.Time
	}
	ret["TransactionHistory"] = ret2

	return &ret, nil
}
