/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX. All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/
package cmd

import (
	"context"
	"os"
	"path/filepath"

	"github.com/IBAX-io/go-explorer/models"

	"github.com/IBAX-io/go-explorer/conf"
	"github.com/IBAX-io/go-explorer/daemons"
	"github.com/IBAX-io/go-explorer/route"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-explorer",
	Short: "scan application",
}

func init() {
	rootCmd.AddCommand(
		initDatabaseCmd,
		startCmd,
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.WithError(err).Fatal("Executing root command")
	}
}

func loadStartRun() error {
	conf.Initer()

	defer models.GormClose()

	go models.CreateCrontab()
	//wg := &sync.WaitGroup{}
	//wg.Add(1)

	daemonsChan := daemons.StartDaemons(context.Background())

	route.Run(conf.GetEnvConf().ServerInfo.Str())
	<-daemonsChan
	//
	//sigChan := make(chan os.Signal)
	//signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	//go func() {
	//	for {
	//		select {
	//		case <-daemonsChan:
	//			wg.Done()
	//		case <-sigChan:
	//			models.GormClose()
	//		}
	//	}
	//}()
	//wg.Wait()
	return nil
}

// Load the configuration from file
func loadInitDatabase() error {
	models.InitDatabase()
	return nil
}

//
func loadConfigWKey(cmd *cobra.Command, args []string) {
	conf.LoadConfig(conf.GetEnvConf().ConfigPath)
}
