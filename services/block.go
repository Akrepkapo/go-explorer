/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/
package services

import (

	//"strings"
	//	"fmt"

	"github.com/IBAX-io/go-explorer/models"
	"github.com/IBAX-io/go-ibax/packages/converter"
	//"github.com/gin-gonic/gin"
	//"github.com/gin-gonic/gin/binding"
)

func Get_Group_Block_Lists(ids int, icount int, order string) (*[]models.BlocksResult, error) {

	//ret, _, err := models.GetBlocklist(ids, icount, order)
	ret, _, err := models.GetBlockListFromRedis()
	return ret, err

}

func Get_Group_Block_TpsLists() (*[]models.ScanOutBlockTransactionRet, error) {
	ret, err := models.GetBlockTpslistsFromRedis()
	return ret, err
}

func Get_Group_Block_Details(id int64) (*models.BlockDetailedInfoHex, error) {
	ret, _, err := GetBlockDetailed(id)
	if err == nil && ret.Header.BlockID > 0 {

		if ret.Header.BlockID-1 > 0 {
			ret1, _, err1 := GetBlockDetailed(ret.Header.BlockID - 1)
			if err1 == nil || ret.Header.BlockID > 0 {
				ret.Header.PreHash = ret1.Hash
				return ret, err1
			}

	bks := &models.Block{}
	ret, err := bks.Get(id)
	if err == nil && ret {
		bkx, err := models.GetBlocksDetailedInfoHex(bks)
		if err == nil {
			//rb.Data = bk
			var block = models.BlocksResult{
				BlockID:      bkx.Header.BlockID,
				Time:         bkx.Time,
				EcosystemID:  bkx.EcosystemID,
				KeyID:        bkx.KeyID,
				NodePosition: bkx.NodePosition,
				Hash:         bkx.Hash,
				Tx:           bkx.Tx,
			}

			return bkx, &block, err
		} else {
			return nil, nil, err
		}

	}
	return nil, nil, err
}

func GetBlockHash(hash []byte) (*models.BlockDetailedInfoHex, bool, error) {
	retd := &models.BlockDetailedInfoHex{}
	bk := &models.Block{}
	found, err := bk.GetBlocksHash(hash)
	if err != nil || !found {
		return retd, found, err
	}
	ret, err := models.GetBlocksDetailedInfoHex(bk)
	if err != nil {
		return retd, found, err
	}
	if ret.Header.BlockID-1 > 0 {
		found1, err1 := bk.Get(ret.Header.BlockID - 1)
		if err1 != nil || !found1 {
			return ret, found, err
		}
		bkx1, err1 := models.GetBlocksDetailedInfoHex(bk)
		if err1 != nil {
			return retd, found, err1
		}
		ret.Header.PreHash = bkx1.Hash
		err = err1
	}

	return ret, found, err
}
