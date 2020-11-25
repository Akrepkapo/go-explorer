/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package daemons

import (
	"context"
	"github.com/IBAX-io/go-explorer/services"
)

func StartDaemons(ctx context.Context) <-chan error {
	exitCh := make(chan error)
	go func() {
		err := services.DealNodetransactionstatussqlite(ctx)
		if err != nil {
			exitCh <- err
		}
	}()

	go func() {
		err := services.DealNodeblocktransactionchsqlite(ctx)
		if err != nil {
			exitCh <- err

		}
	}()
	go Sys_BlockWork(ctx)

	go Sys_Work_ChainValidBlock(ctx)

	go func() {
		err := EcosystemDealupdate(ctx)
		if err != nil {
			exitCh <- err
		}
	}()
	go func() {
		err := NodeTranStatusSumupdate(ctx)
		if err != nil {
			exitCh <- err

