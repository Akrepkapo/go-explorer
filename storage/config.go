package storage

import (
	"time"

	"gorm.io/gorm"
)

var (
	fullNodedb []*FullNodeDB
)

type FullnodeModel struct {
	Enable       bool      `gorm:"not null" yaml:"enable" json:"enable"`
	Nodename     string    `gorm:"not null" yaml:"nodename" json:"nodename"`
	TCPAddress   string    `gorm:"not null" yaml:"tcp_address" json:"tcp_address"`
	APIAddress   string    `gorm:"not null" yaml:"api_address" json:"api_address"`
	City         string    `yaml:"city" json:"city"`
	Icon         string    `yaml:"icon" json:"icon"`
	IconUrl      string    `yaml:"icon_url" json:"icon_url"`
	NodePosition int64     `gorm:"primary_key;not_null" yaml:"nodeposition" json:"node_position"`
	KeyID        string    `gorm:"not null" yaml:"key_id" json:"key_id"`
	PublicKey    string    `gorm:"not null" yaml:"public_key" json:"public_key"`
	UnbanTime    time.Time `yaml:"unbantime" json:"unbantime,omitempty"`
	Latitude     string    `yaml:"latitude" json:"latitude,omitempty"`
	Longitude    string    `yaml:"longitude" json:"longitude,omitempty"`
	Name         string    `yaml:"name" json:"name,omitempty"`
	Display      bool      `json:"display,omitempty" yaml:"display"`
}

type FullNodeDB struct {
	Enable          bool   `yaml:"enable"`
	NodeName        string `yaml:"node_name" json:"node_name"`
	NodePosition    int64  `yaml:"node_position" json:"node_position"`
	Engine          string `yaml:"engine" json:"engine"`
	Connect         string `yaml:"connect" json:"connect"`
	Nodestatusstime time.Time
	DBConn          *gorm.DB
}

type Crontab struct {
	FullNodeTime   string `yaml:"fullnodeTime"`
	BlockchainTime string `yaml:"blockchainTime"`
	Historyupdate  string `yaml:"historyupdate"`
	Statistics     string `yaml:"statistics"`
	Transaction    string `yaml:"transaction"`
}

type FullNodeModels []*FullnodeModel
	//	dbcom.NodePosition = f[i].NodePosition
	//	db, err := GormDBInit(f[i].Engine, f[i].Connect)
	//	if err != nil {
	//		dbcom.Enable = false
	//	} else {
	//		dbcom.Enable = true
	//	}
	//	dbcom.DBConn = db
	//	fullNodedb = append(fullNodedb, dbcom)
	//}
	return nil
}

func Connes() []*FullNodeDB {
	return fullNodedb
}

func (r FullNodeModels) Infos() FullNodeModels {
	return r
}

func (r FullNodeModels) Close() error {
	for i := 0; i < len(fullNodedb); i++ {
		if fullNodedb[i].DBConn != nil {
			sqlDB, err := fullNodedb[i].DBConn.DB()
			if err != nil {
				return err
			}
			fullNodedb[i].DBConn = nil
			if err = sqlDB.Close(); err != nil {
				return err
			}
		}
	}
	return nil
}
