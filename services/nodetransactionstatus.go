/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/
package services

import (
	"context"
	"time"

	"github.com/IBAX-io/go-explorer/storage"

	"github.com/IBAX-io/go-explorer/models"
	log "github.com/sirupsen/logrus"
)

type NodeTransactionStatus struct {
	Nodename     string `yaml:"nodename" json:"nodename"`
	NodePosition int64  `yaml:"nodeposition" json:"nodeposition"`
	Data         *[]models.TransactionStatus
}

var (
	NodeTranStatusDaemonCh = make(chan *NodeTransactionStatus, 100)
)

func DealGetnodetransactionstatus(node *storage.FullNodeDB) (int64, error) {
	bk := &models.TransactionStatus{}
	var count int64
	if node.Enable {
		err := node.DBConn.Table("transactions").Count(&count).Error
		if err != nil {
			log.Info("models.DBconnGetTransactionlist transactions false: " + node.NodeName + err.Error())
		}
		if node.Nodestatusstime.IsZero() {
			ret, err := bk.DBconnGetTransactionlist(node.DBConn)
			if err != nil {
				log.Info("models.DBconnGetTransactionlist false: " + node.NodeName + err.Error())
				return 0, err
			}
			if len(*ret) > 0 {
				dat := NodeTransactionStatus{}
				dat.Nodename = node.NodeName
				dat.NodePosition = node.NodePosition
				dat.Data = ret
				NodeTranStatusDaemonCh <- &dat

				node.Nodestatusstime = time.Now()
				node.Nodestatusstime = node.Nodestatusstime.AddDate(0, 0, -1)
			}
		} else {
			ret, err := bk.DBconnGetTimelimit(node.DBConn, node.Nodestatusstime)
			if err != nil {
				log.Info("models.nodeGetTimelimit false: " + node.NodeName + err.Error())
			} else if len(*ret) > 0 {
				dat := NodeTransactionStatus{}
				dat.Nodename = node.NodeName
				dat.NodePosition = node.NodePosition
				dat.Data = ret
				NodeTranStatusDaemonCh <- &dat
				node.Nodestatusstime = time.Now()
				node.Nodestatusstime = node.Nodestatusstime.AddDate(0, 0, -1)
			}
		}
	}

	return count, nil
}

func DealNodetransactionstatussqlite(ctx context.Context) error {
	bk := &models.TransactionStatus{}
	for {
		select {
		case <-ctx.Done():
			return nil
