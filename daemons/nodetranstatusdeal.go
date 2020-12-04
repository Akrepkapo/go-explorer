/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package daemons

import (
	"context"
	"time"

	"github.com/IBAX-io/go-explorer/conf"

	"github.com/IBAX-io/go-explorer/consts"
	"github.com/IBAX-io/go-explorer/models"
	"github.com/IBAX-io/go-explorer/services"
	//"encoding/hex"
)

func NodeTranStatusSumupdate(ctx context.Context) error {
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(30 * time.Second):
			dlen := len(conf.GetFullNodesDbConn())
			maxlen = 0
			for i := 0; i < dlen; i++ {
				mlen, _ := services.DealGetnodetransactionstatus(conf.GetFullNodesDbConn()[i])
				if mlen > maxlen {
					maxlen = mlen
				}
			}
			//set
			var bc models.BlockID
			bc.Time = time.Now().Unix()
			bc.Name = consts.TransactionsMax
			bc.ID = maxlen
			bc.InsertRedis()

		}
	}
}
