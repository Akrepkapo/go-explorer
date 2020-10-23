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

		}
	}()

	go Sys_CentrifugoWork(ctx)
	return exitCh
}
