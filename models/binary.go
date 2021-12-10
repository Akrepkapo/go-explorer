package models

import (
	"path"
	"strconv"
	"strings"

	"github.com/IBAX-io/go-explorer/conf"

	"github.com/IBAX-io/go-ibax/packages/converter"
)

const BinaryTableSuffix = "_binaries"

// Binary represents record of {prefix}_binaries table
type Binary struct {
	Ecosystem int64
	ID        int64
	AppId     int64
	Name      string
	Data      []byte
	Hash      string
	MimeType  string
	Account   string
}

// SetTablePrefix is setting table prefix
func (b *Binary) SetTablePrefix(prefix string) {
	b.Ecosystem = converter.StrToInt64(prefix)
}

// SetTableName sets name of table
func (b *Binary) SetTableName(tableName string) {
	ecosystem, _ := converter.ParseName(tableName)
	b.Ecosystem = ecosystem
}

// TableName returns name of table
func (b *Binary) TableName() string {
	if b.Ecosystem == 0 {
		b.Ecosystem = 1
	}
	return `1_binaries`
}

// GetByID is retrieving model from db by id
func (b *Binary) GetByID(id int64) (bool, error) {
	return isFound(conf.GetDbConn().Conn().Where("id=?", id).First(b))
}

func (b *Binary) GetByPng(d *Binary) (bool, error) {
	return isFound(conf.GetDbConn().Conn().Where("ecosystem = ? and app_id=? and hash = ? and mime_type !=?", d.Ecosystem, d.AppId, d.Hash, d.MimeType).First(b))
}

func (b *Binary) GetByJpeg() string {
	file := ""
	fileSuffix := path.Ext(b.Name) //
	if fileSuffix == "" {
		if b.MimeType == "image/jpeg" || b.MimeType == "image/jpg" || b.MimeType == "image/png" {
			ns := strings.Split(b.MimeType, "/")
			if len(ns) == 2 {

				file = strconv.FormatInt(b.ID, 10) + "-" + b.Name + "." + ns[1]
				return file
			}
		}
	}

	file = strconv.FormatInt(b.ID, 10) + "-" + b.Name
	return file
}
