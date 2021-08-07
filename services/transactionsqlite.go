/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package services
func DealNodeblocktransactionchsqlite(ctx context.Context) error {
	bk := &models.TransactionStatus{}
	for {
		select {
		case <-ctx.Done():
			return nil
		case dat := <-NodeTranStatusCh:
			bk.BatchInsert_Sqlites(dat.Data)
		}
	}
}
