package models

import (
	"github.com/IBAX-io/go-explorer/conf"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

func CreateCrontab() {
	CrontabInfo := conf.GetEnvConf().Crontab
	if CrontabInfo != nil {
		go CreateCronTimeFromFullNode(CrontabInfo.FullNodeTime)

}

func CreateCronTimeFromFullNode(timeSet string) {
	c := NewWithSecond()
	_, err := c.AddFunc(timeSet, func() {
		getFullNodeInfoFromDb()
	})
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("CreateCronTimeFromFullNode addfunc failed")
	}
	c.Start()
}

func NewWithSecond() *cron.Cron {
	secondParser := cron.NewParser(cron.Second | cron.Minute |
		cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor)
	return cron.New(cron.WithParser(secondParser), cron.WithChain())
}
func CreateCronTimeFromBlockchain(timeSet string) {
	c := NewWithSecond()
	_, err := c.AddFunc(timeSet, func() {
		SyncBlockinfoToRedis()
	})
	if err != nil {
		log.WithFields(log.Fields{"error": err, "timeset": timeSet}).Error("CreateCronTimeFromBlockchain addfunc failed")
	}
	c.Start()
}
func CreateCronTimeFromStatistics(timeSet string) {
	c := NewWithSecond()
	_, err := c.AddFunc(timeSet, func() {
		//if err := getStatisticsToRedis(); err != nil {
		//	log.WithFields(log.Fields{"error": err}).Error("getStatisticsToRedis failed")
		//}
		SendStatisticsSignal()
	})
	if err != nil {
		log.WithFields(log.Fields{"error": err, "timeset": timeSet}).Error("CreateCronTimeFromStatistics addfunc failed")
	}
	c.Start()
}
func EcosystemDashboard_historyupdate(timeSet string) {
	c := NewWithSecond()
	_, err := c.AddFunc(timeSet, func() {
		if err := DealRedisDashboardHistoryMap(); err != nil {
			log.WithFields(log.Fields{"error": err}).Error("DealRedisDashboardHistoryMap failed")
		}
	})
	if err != nil {
		log.WithFields(log.Fields{"error": err, "timeset": timeSet}).Error("EcosystemDashboard_historyupdate addfunc failed")
	}
	c.Start()
}
func CreateCrontabFromTransaction(timeSet string) {
	c := NewWithSecond()
	_, err := c.AddFunc(timeSet, func() {
		if err := getTransactionBlockToRedis(); err != nil {
			log.WithFields(log.Fields{"error": err}).Error("getTransactionBlockToRedis failed")
		}
	})
	if err != nil {
		log.WithFields(log.Fields{"error": err, "timeset": timeSet}).Error("CreateCrontabFromTransaction addfunc failed")
	}
	c.Start()
}
