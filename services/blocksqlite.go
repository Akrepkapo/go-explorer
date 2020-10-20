/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package services

import (
	"errors"
	"sync/atomic"

	"github.com/IBAX-io/go-explorer/models"
	"github.com/IBAX-io/go-ibax/packages/consts"
	log "github.com/sirupsen/logrus"
)

type NodeBlockData struct {
	Data *[]models.Block
}

var (
	Sqlite_MaxBlockid int64
	bgOnWorkRun       uint32
)

func WorkDealBlock() error {
	if atomic.CompareAndSwapUint32(&bgOnWorkRun, 0, 1) {
		defer atomic.StoreUint32(&bgOnWorkRun, 0)
	} else {
		return nil
	}

	var bm, bc models.BlockID
	fm, errm := bm.GetbyName(consts.MintMax)
	if errm != nil {
		if (errm.Error() == "redis: nil" || errm.Error() == "EOF") && !fm {
			bm.ID = 0
			bm.Time = 1
			bm.Name = consts.MintMax
			err := bm.InsertRedis()
			if err != nil {
				return err
			}
		} else {
			return errm
		}
	}

	fc, errc := bc.GetbyName(consts.ChainMax)
	if errc != nil {
		return errc
	}
	if !fc || !fm {
		return errors.New("mint or chain block id  not found")
	}
	count := bc.ID - bm.ID
	sc := bm.ID + 1
	elen := sc + count
	//fmt.Printf("sc:%d,elen:%d,count:%d\n", sc, elen, count)

	for i := int(sc); i <= int(elen); i++ {
		bid := int64(i)

		var mc models.ScanOut
		f, err := mc.Get(bid)
		if f && err == nil {

			err = mc.Insert_Redis()
			if err != nil {
				log.Info(err.Error())
			}
			if bid > bm.ID {
				bm.ID = bid
				bm.Time = mc.Time
				bm.InsertRedis()
			}
		}

		if err != nil {
			if err.Error() == "redis: nil" || err.Error() == "EOF" {
				log.Info("redis: nil")
				//break
			} else {
				break
			}
		}

		bdid := int64(i)
		mc.Del(bdid - 1)
	}

	return nil
