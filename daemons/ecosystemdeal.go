/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/
package daemons

import (
	//"fmt"
	"context"
	"encoding/json"
	"sync/atomic"
	"time"

	"github.com/IBAX-io/go-explorer/models"
	"github.com/IBAX-io/go-ibax/packages/consts"
	"github.com/IBAX-io/go-ibax/packages/model"

	"github.com/IBAX-io/go-explorer/services"
	//"github.com/IBAX-io/go-explorer/models"
	log "github.com/sirupsen/logrus"
)

//recvice websocket request sendTopData
func EcosystemDealupdate(ctx context.Context) error {
	//bk := &models.TransactionStatus{}
	services.Deal_Redis_Dashboard()
	services.SendWebsocketData = make(chan bool, 1)
	for {
		select {
		case <-ctx.Done():
			log.Error("NodeTranStatusSumupdate done his work")
			return nil
		case <-services.SendWebsocketData:
			if err := services.Deal_Redis_Dashboard(); err != nil {
				log.Info("send topdata err:", err)
			}

		}
	}
}

func Sys_BlockWork(ctx context.Context) {
	models.GetScanOut = make(chan bool, 1)
	for {
		select {
		case <-ctx.Done():
			return
		case <-models.GetScanOut:
			if err := models.GetScanOutDataToRedis(); err != nil {
				log.Info("GetScanOutDataToRedis failed:", err)
			}
			err := services.WorkDealBlock()
			if err != nil {
				log.Info("WorkDealBlock", err)
			}
		}
	}
}

func Sys_Work_ChainValidBlock(ctx context.Context) {
	err := ChainValidBlock()
	if err != nil {
		log.Info("first Sys_Work_ChainValidBlock")
	}

	for {
		select {
		case <-ctx.Done():
			log.Info("Sys_Work_ChainValidBlock done his work")
			return
		case <-time.After(time.Second * 4):
			err := ChainValidBlock()
			if err != nil {
				log.Info("ChainValidBlock")
			}
		}
	}
}

var bgOnChaimRun uint32

func ChainValidBlock() error {
	if atomic.CompareAndSwapUint32(&bgOnChaimRun, 0, 1) {
		defer atomic.StoreUint32(&bgOnChaimRun, 0)
	} else {
		return nil
	}

	var cf model.Confirmation
	var bc models.BlockID
	f, err := cf.GetGoodBlockLast()
	if err != nil {
		return err
	}
	if f {
		fc, err := bc.GetbyName(consts.ChainMax)
		if err != nil {
			if err.Error() == "redis: nil" || err.Error() == "EOF" {
			} else {
				return err
			}
		}
		if !fc {
			bc.Time = cf.Time
			bc.Name = consts.ChainMax
			bc.ID = cf.BlockID
			return bc.InsertRedis()
		}
		if cf.BlockID > bc.ID {
			bc.ID = cf.BlockID
			bc.Time = cf.Time
			return bc.InsertRedis()
		}
	}

	return nil
}

func Sys_CentrifugoWork(ctx context.Context) {
	models.SendScanOut = make(chan bool, 1)
	for {
		select {
		case <-ctx.Done():
			return
	}
}

func SendtoWebsocket(rets *models.ScanOutRet, scanout *models.ScanOut) error {

	data, _ := json.Marshal(rets)
	err := services.WriteChannelByte(services.ChannelDashboard, data)
	if err != nil {
		return err
	}

	//trans, err := scanout.GetBlockTransactions(15)
	//if err != nil {
	//	log.Info("GetBlockTransactions", err.Error())
	//	return err
	//} else {
	//	ds, err := json.Marshal(trans)
	//	if err != nil {
	//		log.Info("json.Marshal", err.Error())
	//		return err
	//	} else {
	//		err := services.WriteChannelByte(services.ChannelBlockAndTxsList, ds)
	//		if err != nil {
	//			log.Info("WriteChannelByte blocktransactionlist", err.Error())
	//			return err
	//		}
	//	}
	//}

	return nil
}
