/*---------------------------------------------------------------------------------------------
 *  Copyright (c) IBAX All rights reserved.
 *  See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/
package conf

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/IBAX-io/go-explorer/storage"
	DatabaseInfo *storage.DatabaseModel `yaml:"database"`
	RedisInfo    *storage.RedisModel    `yaml:"redis"`
	Url          *UrlModel
	Centrifugo   *storage.CentrifugoConfig
	Crontab      *storage.Crontab `yaml:"crontab"`
}

func GetEnvConf() *EnvConf {
	return &configInfo
}

func GetDbConn() *storage.DatabaseModel {
	return GetEnvConf().DatabaseInfo
}

func GetFullNodesDbConn() []*storage.FullNodeDB {
	return storage.Connes()
}

func GetRedisDbConn() *storage.RedisModel {
	return GetEnvConf().RedisInfo
}
func GetCentrifugoConn() *storage.CentrifugoConfig {
	return GetEnvConf().Centrifugo
}

func LoadConfig(configPath string) {
	filePath := path.Join(configPath, "config.yml")
	configData, err := os.ReadFile(filePath)
	if err != nil {
		logrus.WithError(err).Fatal("config file read failed")
	}
	// expand environment variables
	configData = []byte(os.ExpandEnv(string(configData)))
	err = yaml.Unmarshal(configData, &configInfo)
	data,_ :=json.Marshal(&configInfo)
	fmt.Printf("config: %v\n",string(data))
	if err != nil {
		logrus.WithError(err).Fatal("config parse failed")
	}
}

func Initer() {
	DatabaseInfo := GetEnvConf().DatabaseInfo
	RedisInfo := GetEnvConf().RedisInfo
	Centrifugo := GetEnvConf().Centrifugo

	if err := DatabaseInfo.Initer(); err != nil {
		logrus.WithError(err).Fatal("postgres database connect failed: %v", DatabaseInfo.Connect)
	}
	if err := RedisInfo.Initer(); err != nil {
		logrus.WithError(err).Fatal("redis database config information: %v", RedisInfo)
	}
	if err := Centrifugo.Initer(); err != nil {
		logrus.WithError(err).Fatal("centrifugo config information: %v", Centrifugo)
	}
	if err := initLogs(); err != nil {
		logrus.WithError(err).Fatal("init log file")
	}
}

func initLogs() error {
	fileName := path.Join(GetEnvConf().ConfigPath, "logrus.log")
	openMode := os.O_APPEND
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		openMode = os.O_CREATE
	}
	f, err := os.OpenFile(fileName, os.O_WRONLY|openMode, 0755)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Can't open log file: ", fileName)
		return err
	}
	logrus.SetOutput(f)
	return nil
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
