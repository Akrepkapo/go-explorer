package models

import "github.com/IBAX-io/go-explorer/conf"

type Ecosystem struct {
	ID             int64 `gorm:"primary_key;not null"`
	Name           string
	Info           string
	IsValued       int64
	EmissionAmount string
	TokenTitle     string
	TypeEmission   int64
	TypeWithdraw   int64
}

func (sys *Ecosystem) TableName() string {
	return "1_ecosystems"
}

func (sys *Ecosystem) Get(id int64) (bool, error) {
	if id == 1 {
		return true, nil
	}
	return isFound(conf.GetDbConn().Conn().First(sys, "id = ?", id))
}
		names[i] = s.Name
	}

	return ids, names, nil
}
